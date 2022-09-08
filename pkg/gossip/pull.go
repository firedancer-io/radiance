package gossip

import (
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"net/netip"
	"time"
)

const PacketSize = 1232

// PullClient implements the stateful client (initiator) side of the gossip pull protocol.
type PullClient struct {
	identity ed25519.PrivateKey
	so       udpSender
}

func NewPullClient(identity ed25519.PrivateKey, so udpSender) *PullClient {
	return &PullClient{
		identity: identity,
		so:       so,
	}
}

func (p *PullClient) Pull(target netip.AddrPort) error {
	filters := NewCrdsFilterSet(65536, MaxBloomSize)
	for _, filter := range filters {
		if err := p.sendPullRequest(target, filter); err != nil {
			return err
		}
	}
	return nil
}

func (p *PullClient) sendPullRequest(target netip.AddrPort, filter CrdsFilter) error {
	msg := &Message__PullRequest{
		Filter: filter,
		Value: CrdsValue{
			Data: &CrdsData__ContactInfo{
				Value: ContactInfo{
					Wallclock: uint64(time.Now().UnixMilli()),
				},
			},
		},
	}
	err := msg.Value.Sign(p.identity)
	if err != nil {
		panic("failed to sign pull request: " + err.Error())
	}

	packet, err := msg.BincodeSerialize()
	if err != nil {
		panic("failed to serialize packet: " + err.Error())
	}

	_, err = p.so.WriteToUDPAddrPort(packet, target)
	return err
}

func (p *PullClient) HandlePullResponse(msg *Message__PullResponse, _ netip.AddrPort) {
	jsonBuf, _ := json.MarshalIndent(msg, "", "\t")
	fmt.Println(string(jsonBuf))
}
