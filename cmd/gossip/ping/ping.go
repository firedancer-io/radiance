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
	"sync/atomic"
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

func init() {
	klog.InitFlags(nil)
	flag.Parse()
}

func main() {
	if *flagAddr == "" {
		klog.Exit("No address to ping specified")
	}

	ctx := context.Background()

	_, privkey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	conn, err := net.Dial("udp", *flagAddr)
	if err != nil {
		klog.Exit(err)
	}
	udpConn := conn.(*net.UDPConn)

	klog.Infof("GOSSIP PING %s (%s)", *flagAddr, conn.RemoteAddr())

	s := Session{
		privkey: privkey,
		udpConn: udpConn,
		reqs:    make(map[[32]byte]pending),
	}

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		return s.send(ctx)
	})
	group.Go(func() error {
		return s.receive(ctx)
	})
	_ = group.Wait()
	_ = conn.Close()

	klog.Infof("--- %s gossip ping statistics ---", udpConn.RemoteAddr())

	numSuccess := atomic.LoadUint64(&s.numSuccess)
	numTimeout := atomic.LoadUint64(&s.numTimeout)
	klog.Infof("%d packets transmitted, %d packets received, %.1f%% packet loss",
		numSuccess+numTimeout, numSuccess, (1-(float64(numSuccess)/float64(numSuccess+numTimeout)))*100)
}

type pending struct {
	c      int
	t      time.Time
	ping   gossip.Ping
	cancel context.CancelFunc
}

type Session struct {
	privkey ed25519.PrivateKey
	udpConn *net.UDPConn
	lock    sync.Mutex
	reqs    map[[32]byte]pending

	numSuccess  uint64
	numTimeout  uint64
	numSendFail uint64
}

func (s *Session) send(ctx context.Context) error {
	defer s.udpConn.Close()
	ticker := time.NewTicker(*flagDelay)
	for c := 0; c < *flagCount || *flagCount == -1; c++ {
		s.sendPing(ctx, c)
		select {
		case <-ticker.C:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return errors.New("done")
}

func (s *Session) sendPing(ctx context.Context, c int) {
	t := time.Now()

	var token [32]byte
	if _, err := rand.Read(token[:]); err != nil {
		panic(err)
	}

	ping := gossip.NewPing(token, s.privkey)
	pingMsg := gossip.Message__Ping{Value: ping}
	frame, err := pingMsg.BincodeSerialize()
	if err != nil {
		panic(err)
	}

	_, _, err = s.udpConn.WriteMsgUDP(frame, nil, nil)
	if err != nil {
		klog.Warningf("Send ping: %s", err)
		atomic.AddUint64(&s.numSendFail, 1)
		return
	}

	pongToken := gossip.HashPingToken(ping.Token)
	pingCtx, pingCancel := context.WithTimeout(ctx, *flagTimeout)
	go func() {
		<-pingCtx.Done()
		if err := pingCtx.Err(); err == context.DeadlineExceeded {
			s.lock.Lock()
			defer s.lock.Unlock()
			delete(s.reqs, pongToken)
			klog.V(3).Infof("Request timeout for seq %d", c)
			atomic.AddUint64(&s.numTimeout, 1)
		}
	}()

	s.lock.Lock()
	s.reqs[pongToken] = pending{
		c:      c,
		t:      t,
		ping:   ping,
		cancel: pingCancel,
	}
	s.lock.Unlock()
}

func (s *Session) receive(ctx context.Context) error {
	for ctx.Err() == nil {
		var packet [132]byte
		n, remote, err := s.udpConn.ReadFromUDPAddrPort(packet[:])
		klog.V(7).Infof("Packet from %s", remote)
		if n >= len(packet) {
			s.handlePong(packet[:], remote)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Session) handlePong(packet []byte, remote netip.AddrPort) {
	msg, err := gossip.BincodeDeserializeMessage(packet)
	if err != nil {
		return
	}
	pongMsg, ok := msg.(*gossip.Message__Pong)
	if !ok {
		return
	}

	s.lock.Lock()
	defer s.lock.Unlock()
	req, ok := s.reqs[pongMsg.Value.Token]
	if !ok {
		return
	}
	delete(s.reqs, pongMsg.Value.Token)

	req.cancel()

	klog.V(3).Infof("Pong from %s seq=%d time=%v", remote, req.c, time.Since(req.t))
	atomic.AddUint64(&s.numSuccess, 1)
}
