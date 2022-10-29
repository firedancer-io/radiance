package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"go.firedancer.io/radiance/pkg/endpoints"
	"go.firedancer.io/radiance/pkg/netlink"
	"go.firedancer.io/radiance/pkg/tpu"
	"k8s.io/klog/v2"
	"net"
	"strconv"
	"strings"
)

type packet struct {
	data []byte
	port uint16
	src  net.IP
}

func readPacketsFromInterface(iface string, ports []uint16, dst net.IP) (chan packet, error) {
	// bpf filter
	var filter string
	for _, port := range ports {
		filter += fmt.Sprintf(" or dst port %d", port)
	}
	filter = fmt.Sprintf("udp and dst host %s and (%s)", dst.String(), filter[4:])

	klog.Info("filter: ", filter)

	handle, err := pcap.OpenLive(iface, 1600, false, pcap.BlockForever)
	if err != nil {
		return nil, err
	}

	// set filter
	err = handle.SetBPFFilter(filter)
	if err != nil {
		return nil, err
	}

	packets := make(chan packet)
	go func() {
		defer close(packets)
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for p := range packetSource.Packets() {
			udpLayer := p.Layer(layers.LayerTypeUDP)
			if udpLayer == nil {
				continue
			}
			packets <- packet{
				data: udpLayer.(*layers.UDP).Payload,
				port: uint16(udpLayer.(*layers.UDP).DstPort),
			}
		}
	}()
	return packets, nil
}

var (
	flagIface = flag.String("iface", "", "interface to read packets from")
	flagPorts = flag.String("ports", "", "destination ports to sniff (comma-separated), asks local RPC if empty")
)

func main() {
	flag.Parse()
	if *flagIface == "" {
		klog.Exit("-iface is required")
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

	packets, err := readPacketsFromInterface(*flagIface, ports, dst)
	if err != nil {
		klog.Exit("error reading packets: ", err)
	}

	for p := range packets {
		tx, err := tpu.ParseTx(p.data)
		if err != nil {
			klog.Warning("port %d error parsing tx: ", p.port, err)
			continue
		}

		signers := tpu.ExtractSigners(tx)
		klog.Infof("port %d sig %s signers %v", p.port, tx.Signatures[0], signers)
	}
}
