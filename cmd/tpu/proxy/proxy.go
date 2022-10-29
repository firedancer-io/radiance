package main

import (
	"bytes"
	"context"
	"flag"
	"github.com/LiamHaworth/go-tproxy"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.firedancer.io/radiance/pkg/endpoints"
	"go.firedancer.io/radiance/pkg/netlink"
	"go.firedancer.io/radiance/pkg/nftables"
	"golang.org/x/sys/unix"
	"k8s.io/klog/v2"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var (
	flagDebugAddr = flag.String("debug-addr", ":6060", "Metrics and pprof listen address")
	flagIface     = flag.String("iface", "", "External interface to receive packets from")
	flagPorts     = flag.String("ports", "", "Destination ports to proxy (comma-separated), asks local RPC if empty")

	metricPacketsCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "tproxy_packets_count",
		Help: "Number of packets received by the proxy",
	}, []string{"port"})
	metricBytesCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "tproxy_bytes_count",
		Help: "Number of bytes received by the proxy",
	}, []string{"port"})
)

func main() {
	flag.Parse()

	if *flagIface == "lo" {
		klog.Exitf("proxying lo would lead to a loopback packet loop")
	}

	if *flagIface == "" {
		klog.Exitf("no interface specified, use -iface to specify one")
	}

	dst, err := netlink.GetInterfaceIP(*flagIface)
	if err != nil {
		klog.Exit("failed to get IP: ", err)
	}

	klog.Infof("interface %s has primary IP %s", *flagIface, dst)

	ports := make([]uint16, 0)

	if *flagPorts == "" {
		klog.Infof("no ports specified, asking local RPC for ports")
		ports, err = endpoints.GetNodeTPUPorts(context.Background(), endpoints.RPCLocalhost, dst)
		if err != nil {
			klog.Exit("failed to get ports: ", err)
		}
		klog.Infof("found ports: %v", ports)
	} else {
		for _, port := range strings.Split(*flagPorts, ",") {
			p, err := strconv.ParseUint(port, 10, 16)
			if err != nil {
				klog.Exit("failed to parse port: ", err)
			}
			ports = append(ports, uint16(p))
		}
	}

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		klog.Infof("Starting pprof and Prometheus server on %s", *flagDebugAddr)
		klog.Fatal(http.ListenAndServe(*flagDebugAddr, nil))
	}()

	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		klog.Fatalf("Failed to get hostname: %v", err)
	}

	klog.Infof("Running on  %s", hostname)

	addr, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		klog.Exitf("Failed to resolve address: %v", err)
	}

	conn, err := tproxy.ListenUDP("udp", addr)
	if err != nil {
		klog.Exitf("Failed to listen on %s: %v", addr, err)
	}

	localPort := uint16(conn.LocalAddr().(*net.UDPAddr).Port)

	defer conn.Close()

	klog.Infof("Listening on %s", conn.LocalAddr())

	go listen(conn)

	if err := nftables.EnsureKernelModules(); err != nil {
		klog.Exitf("Failed to ensure kernel modules: %v", err)
	}

	if err := nftables.InsertProxyChain(ports, localPort, *flagIface); err != nil {
		klog.Exitf("Failed to insert nft tproxy chain: %v", err)
	}
	defer func() {
		err := nftables.DeleteProxyChain()
		if err != nil {
			klog.Warningf("Failed to delete nft tproxy chain: %v", err)
		}
		klog.Infof("Deleted nft tproxy chain")
	}()

	klog.Infof("Inserted nft tproxy chain")

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, unix.SIGINT)
	signal.Notify(sigint, unix.SIGTERM)
	<-sigint

	klog.Infof("Shutting down")
}

func listen(conn *net.UDPConn) {
	var inBytes *uint64
	var inPackets *uint64
	inBytes = new(uint64)
	inPackets = new(uint64)

	// Periodically log stats
	go func() {
		for {
			klog.Infof("InBytes: %d, InPackets: %d", atomic.LoadUint64(inBytes), atomic.LoadUint64(inPackets))
			atomic.StoreUint64(inBytes, 0)
			atomic.StoreUint64(inPackets, 0)

			time.Sleep(time.Second)
		}
	}()

	for {
		buf := make([]byte, 1024)
		n, src, dst, err := tproxy.ReadFromUDP(conn, buf)
		if err != nil {
			klog.Errorf("Failed to read from UDP: %v", err)
			continue
		}

		if bytes.Equal(src.IP, dst.IP) && src.Port == dst.Port {
			klog.V(2).Infof("src and dst are identical, dropping packet")
			continue
		}

		atomic.AddUint64(inBytes, uint64(n))
		atomic.AddUint64(inPackets, 1)
		metricBytesCount.WithLabelValues(strconv.Itoa(int(dst.Port))).Add(float64(n))
		metricPacketsCount.WithLabelValues(strconv.Itoa(int(dst.Port))).Add(1)

		go handlePacket(conn, buf, src, dst)
	}
}

func handlePacket(conn *net.UDPConn, buf []byte, src, dst *net.UDPAddr) {
	klog.V(2).Infof("Received %d bytes from %s", len(buf), src)

	_, err := conn.WriteToUDP(buf, dst)
	if err != nil {
		klog.Errorf("Failed to write to UDP: %v", err)
		return
	}

	klog.V(2).Infof("Sent %d bytes to %s", len(buf), dst)
}
