// Ping sends gossip pings to a Solana node.
package main

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"flag"
	"net"
	"net/netip"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/certusone/radiance/pkg/gossip"
	"golang.org/x/sync/errgroup"
	"k8s.io/klog/v2"
)

var (
	flagCount   = flag.Int("c", -1, "Number of pings to send, -1 for infinite")
	flagDelay   = flag.Duration("i", 1*time.Second, "Delay between pings")
	flagTimeout = flag.Duration("timeout", 3*time.Second, "Ping timeout")
	flagAddr    = flag.String("addr", "", "Address to ping (<host>:<port>)")
)

var target netip.AddrPort

func init() {
	klog.InitFlags(nil)
	flag.Parse()
}

func main() {
	if *flagAddr == "" {
		klog.Exit("No address to ping specified")
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
	handler := &gossip.Handler{
		PingClient: pingClient,
	}
	client := gossip.NewDriver(handler, conn)

	klog.Infof("GOSSIP PING %s (%s)", *flagAddr, target.String())

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		return client.Run(ctx)
	})
	group.Go(func() error {
		pingLoop(ctx, pingClient)
		return context.Canceled // done
	})
	_ = group.Wait()
	_ = conn.Close()

	klog.Infof("--- %s gossip ping statistics ---", target.String())

	numSuccess := pingClient.NumOK.Load()
	numTimeout := pingClient.NumTimeout.Load()
	klog.Infof("%d packets transmitted, %d packets received, %.1f%% packet loss",
		numSuccess+numTimeout, numSuccess, (1-(float64(numSuccess)/float64(numSuccess+numTimeout)))*100)
}

func pingLoop(ctx context.Context, pinger *gossip.PingClient) {
	var wg sync.WaitGroup
	defer wg.Wait()

	ticker := time.NewTicker(*flagDelay)
	count := *flagCount
	for seq := 0; seq < count || count == -1; seq++ {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			wg.Add(1)
			go sendPing(ctx, &wg, pinger, seq)
		}
	}
}

func sendPing(ctx context.Context, wg *sync.WaitGroup, pinger *gossip.PingClient, seq int) {
	defer wg.Done()

	ctx, cancel := context.WithTimeout(ctx, *flagTimeout)
	defer cancel()

	start := time.Now()
	_, responder, err := pinger.Ping(ctx, target)
	if err == nil {
		klog.Infof("Pong from %s seq=%d time=%v", responder, seq, time.Since(start))
	} else if errors.Is(err, context.Canceled) {
		return
	} else if errors.Is(err, context.DeadlineExceeded) {
		klog.Infof("Request timeout for seq %d", seq)
	} else {
		klog.Warning(err)
	}
}
