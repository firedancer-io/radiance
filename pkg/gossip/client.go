package gossip

import (
	"context"
	"net"
	"net/netip"
	"sync/atomic"
)

// Driver implements the network main loop.
//
// Note: This uses Go 1.19 standard library networking, which processes packets one-by-one. (slow!)
// Looks like Go 1.20 will add batch packet receive: https://github.com/golang/go/issues/45886
type Driver struct {
	handler *Handler
	so      *net.UDPConn
}

func NewDriver(handler *Handler, so *net.UDPConn) *Driver {
	return &Driver{
		handler: handler,
		so:      so,
	}
}

// Run processes packets until the context is cancelled.
//
// Destroys all handlers and closes the socket after returning.
// Returns any network error or nil if the context closed.
func (c *Driver) Run(ctx context.Context) error {
	defer c.handler.Close()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var graceful atomic.Bool
	// Monitor close signal
	go func() {
		defer c.so.Close()
		defer graceful.Store(true)
		<-ctx.Done()
	}()

	var buf [PacketSize]byte
	for {
		n, _, _, addr, err := c.so.ReadMsgUDPAddrPort(buf[:], nil)
		if n > 0 {
			c.handler.HandlePacket(buf[:n], addr)
		}
		if err != nil {
			if graceful.Load() {
				return nil
			}
			return err
		}
	}
}

// Handler is a network-agnostic multiplexer for incoming gossip messages.
type Handler struct {
	*PullClient
	*PingClient
	*PingServer

	numInvalidMsgs uint64
	numIgnoredMsgs uint64
}

// HandlePacket is the entrypoint of the RX side.
func (h *Handler) HandlePacket(packet []byte, from netip.AddrPort) {
	msg, err := BincodeDeserializeMessage(packet)
	if err != nil {
		atomic.AddUint64(&h.numInvalidMsgs, 1)
		return
	}
	switch x := msg.(type) {
	case *Message__PullResponse:
		if h.PullClient != nil {
			h.PullClient.HandlePullResponse(x, from)
			return
		}
	case *Message__Ping:
		if h.PingServer != nil {
			h.PingServer.HandlePing(x, from)
			return
		}
	case *Message__Pong:
		if h.PingClient != nil {
			h.PingClient.HandlePong(x, from)
			return
		}
	}
	atomic.AddUint64(&h.numIgnoredMsgs, 1)
}

// Close destroys all handlers.
func (h *Handler) Close() {
	h.PingClient.Close()
}

type udpSender interface {
	WriteToUDPAddrPort(b []byte, addr netip.AddrPort) (int, error)
}
