// Pull downloads CRDS from a peer.
package main

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"flag"
	"net"
	"net/netip"
	"os"
	"os/signal"
	"time"

	"github.com/certusone/radiance/pkg/gossip"
	"golang.org/x/sync/errgroup"
	"k8s.io/klog/v2"
)

var (
	flagAddr = flag.String("addr", "", "Address to ping (<host>:<port>)")
)

var target netip.AddrPort

func init() {
	klog.InitFlags(nil)
	flag.Parse()
}

func main() {
	if *flagAddr == "" {
		klog.Exit("No address specified")
	}

	udpAddr, err := net.ResolveUDPAddr("udp", *flagAddr)
	if err != nil {
		klog.Exitf("invalid target address: %s", err)
	}
	target = udpAddr.AddrPort()

	ctx := context.Background()

	_, identity, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		klog.Exit(err)
	}

	pingClient := gossip.NewPingClient(identity, conn)
	pingServer := gossip.NewPingServer(identity, conn)
	pullClient := gossip.NewPullClient(identity, conn)
	handler := &gossip.Handler{
		PingClient: pingClient,
		PingServer: pingServer,
		PullClient: pullClient,
	}
	client := gossip.NewDriver(handler, conn)

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		return client.Run(ctx)
	})
	group.Go(func() error {
		_, _, err := pingClient.Ping(ctx, target)
		if err != nil {
			klog.Exitf("ping ignored: %s", err)
		}
		klog.Info("Pinged")
		err = pullClient.Pull(target)
		if err != nil {
			klog.Exitf("failed to send pull request: %s", err)
		}
		err = pullClient.Pull(target)
		if err != nil {
			klog.Exitf("failed to send pull request: %s", err)
		}

		time.Sleep(2 * time.Second)

		err = pullClient.Pull(target)
		if err != nil {
			klog.Exitf("failed to send pull request: %s", err)
		}

		time.Sleep(3 * time.Second)
		return context.Canceled // done
	})
	_ = group.Wait()
	_ = conn.Close()
}
