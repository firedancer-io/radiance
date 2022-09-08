package gossip

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"net"
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestPingServer(t *testing.T) {
	conn, err := net.ListenUDP("udp", net.UDPAddrFromAddrPort(netip.MustParseAddrPort("[::1]:0")))
	require.NoError(t, err)
	defer conn.Close()

	lo := conn.LocalAddr().(*net.UDPAddr).AddrPort()

	_, identity, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)

	handler := &Handler{
		PingClient: NewPingClient(identity, conn),
		PingServer: NewPingServer(identity, conn),
	}
	client := NewDriver(handler, conn)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		return client.Run(ctx)
	})
	group.Go(func() error {
		defer cancel()
		pinger := handler.PingClient

		for i := 0; i < 100; i++ {
			pong, responder, err := pinger.Ping(ctx, lo)
			require.NoError(t, err)
			assert.Equal(t, lo, responder)
			assert.True(t, pong.Verify())
		}
		return nil
	})
	err = group.Wait()
	assert.NoError(t, err)
}

func BenchmarkPing_SignHashVerify(b *testing.B) {
	_, identity, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		ping := NewPingRandom(identity)
		assert.True(b, ping.Verify())
		pong := NewPing(HashPingToken(ping.Token), identity)
		assert.True(b, pong.Verify())
	}
}
