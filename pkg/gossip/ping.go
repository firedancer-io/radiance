package gossip

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"math/rand"
	"net/netip"
	"sync"
	"sync/atomic"

	"k8s.io/klog/v2"
)

// PingSize is the size of a serialized ping message.
const PingSize = 128

// NewPing creates and signs a new ping message.
//
// Panics if the provided private key is invalid.
func NewPing(token [32]byte, key ed25519.PrivateKey) (p Ping) {
	sig := ed25519.Sign(key, token[:])
	copy(p.From[:], key.Public().(ed25519.PublicKey))
	copy(p.Token[:], token[:])
	copy(p.Signature[:], sig[:])
	return p
}

func NewPingRandom(key ed25519.PrivateKey) Ping {
	var token [32]byte
	rand.Read(token[:])
	return NewPing(token, key)
}

// Verify checks the Ping's signature.
func (p *Ping) Verify() bool {
	return ed25519.Verify(p.From[:], p.Token[:], p.Signature[:])
}

// HashPingToken returns the pong token given a ping token.
func HashPingToken(token [32]byte) [32]byte {
	return sha256.Sum256(token[:])
}

// PingClient implements the stateful client (initiator) side of the gossip ping protocol.
//
// It tracks every pending request to match it with solicited pong frames.
type PingClient struct {
	identity ed25519.PrivateKey
	so       udpSender

	lock sync.RWMutex
	reqs map[Hash]*pingSession // pong token => session

	NumSent     atomic.Uint64 // ping messages sent
	NumOK       atomic.Uint64 // successful ping transaction
	NumInvalid  atomic.Uint64 // invalid sig in pong
	NumTimeout  atomic.Uint64 // context errored before pong arrived
	NumSendFail atomic.Uint64 // socket refused to send (tx buffer full)
	NumMartian  atomic.Uint64 // unsolicited pong
}

func NewPingClient(identity ed25519.PrivateKey, so udpSender) *PingClient {
	return &PingClient{
		identity: identity,
		so:       so,
		reqs:     make(map[Hash]*pingSession),
	}
}

type pingSession struct {
	out  chan<- pongResponse
	done atomic.Bool
}

type pongResponse struct {
	from netip.AddrPort
	pong Ping
}

// Ping sends a gossip ping packet.
// Blocks until a valid matching pong packet arrives or the context is cancelled.
//
// Note that this mechanism is unrelated to ICMP pings.
func (p *PingClient) Ping(ctx context.Context, target netip.AddrPort) (pong Ping, responder netip.AddrPort, err error) {
	// Allocate session for lifetime of scope
	ping, pongToken, resp := p.createSession()
	defer p.destroySession(pongToken)

	// Send ping to "server"
	pingMsg := &Message__Ping{ping}
	packet, err := pingMsg.BincodeSerialize()
	if err != nil {
		klog.Errorf("Failed to serialize ping: %s", err)
		return
	}
	if _, err = p.so.WriteToUDPAddrPort(packet, target); err != nil {
		p.NumSendFail.Add(1)
		return
	}
	p.NumSent.Add(1)

	// Block until something happens
	select {
	case <-ctx.Done():
		err = ctx.Err()
		p.NumTimeout.Add(1)
		return
	case resp, ok := <-resp:
		if !ok {
			// sanity check: cancellation can only happen before or after select
			panic("race condition")
		}
		pong = resp.pong
		responder = resp.from
		p.NumOK.Add(1)
		return
	}
}

// HandlePong processes incoming gossip pong messages.
func (p *PingClient) HandlePong(msg *Message__Pong, from netip.AddrPort) {
	pong := msg.Value

	// map lookup is cheaper than Ed25519 verify, so do that first
	sess := p.getSession(pong.Token)
	if sess == nil {
		p.NumMartian.Add(1)
		return
	}

	// We might receive two valid pongs before the initiating goroutine cleans up.
	// Bail here because we are only allowed to send one pong back to the channel.
	if sess.done.Swap(true) {
		p.NumMartian.Add(1)
		return
	}

	if !pong.Verify() {
		p.NumInvalid.Add(1)
		return
	}

	// Upgrade to write lock to prevent initiating goroutine
	// from closing the channel we're about to send on.
	// TODO: this is probably very slow
	p.lock.Lock()
	defer p.lock.Unlock()
	sess = p.reqs[pong.Token]
	if sess == nil {
		// session was cancelled while we were verifying the pong
		p.NumMartian.Add(1)
		return
	}
	sess.out <- pongResponse{
		from: from,
		pong: pong,
	}
}

func (p *PingClient) createSession() (ping Ping, pongToken Hash, resp <-chan pongResponse) {
	ping = NewPingRandom(p.identity)
	pongToken = HashPingToken(ping.Token)

	respBi := make(chan pongResponse, 1)
	resp = respBi // recv-only

	p.lock.Lock()
	defer p.lock.Unlock()
	p.reqs[pongToken] = &pingSession{
		out: respBi, // send-only
	}
	return
}

func (p *PingClient) getSession(pongToken Hash) *pingSession {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.reqs[pongToken]
}

func (p *PingClient) destroySession(pongToken Hash) {
	p.lock.Lock()
	defer p.lock.Unlock()

	session, ok := p.reqs[pongToken]
	if !ok {
		return
	}
	close(session.out)
	delete(p.reqs, pongToken)
}

func (p *PingClient) Close() {
	p.lock.Lock()
	defer p.lock.Unlock()
}

// PingServer implements the stateless server (reactor) side of the gossip ping protocol.
//
// It implements no rate-limits and is thus vulnerable to packet floods.
type PingServer struct {
	identity ed25519.PrivateKey
	so       udpSender

	NumOK       atomic.Uint64 // handled pings, though uncertain whether pong arrived
	NumInvalid  atomic.Uint64 // invalid sig in ping
	NumSendFail atomic.Uint64 // socket refused to send (tx buffer full)
}

func NewPingServer(identity ed25519.PrivateKey, so udpSender) *PingServer {
	return &PingServer{
		identity: identity,
		so:       so,
	}
}

// HandlePing processes incoming gossip ping messages.
func (p *PingServer) HandlePing(ping *Message__Ping, from netip.AddrPort) {
	// Verify signature of ping.
	if !ping.Value.Verify() {
		p.NumInvalid.Add(1)
		return
	}

	// SHA-256 hash token and sign with identity.
	// Note: Possible signature forging attack vector.
	pong := NewPing(HashPingToken(ping.Value.Token), p.identity)
	pongMsg := &Message__Pong{pong}

	// Respond to sender.
	packet, err := pongMsg.BincodeSerialize()
	if err != nil {
		klog.Errorf("Failed to serialize pong: %s", err)
		return
	}
	if _, err = p.so.WriteToUDPAddrPort(packet, from); err != nil {
		p.NumSendFail.Add(1)
		return
	}
	p.NumOK.Add(1)
}
