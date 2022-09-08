package gossip

import (
	"net/netip"

	"github.com/novifinancial/serde-reflection/serde-generate/runtime/golang/serde"
)

type SocketAddr struct {
	netip.AddrPort
}

func DeserializeSocketAddr(deserializer serde.Deserializer) (sa SocketAddr, err error) {
	var raw RawSocketAddr
	raw, err = DeserializeRawSocketAddr(deserializer)
	if err != nil {
		return
	}
	sa.AddrPort = netip.AddrPortFrom(raw.Addr.Addr, raw.Port)
	return
}

func (sa SocketAddr) Serialize(serializer serde.Serializer) error {
	raw := RawSocketAddr{
		Addr: Addr{sa.Addr()},
		Port: sa.Port(),
	}
	return raw.Serialize(serializer)
}

type Addr struct {
	netip.Addr
}

func DeserializeAddr(deserializer serde.Deserializer) (sa Addr, err error) {
	var raw RawAddr
	raw, err = DeserializeRawAddr(deserializer)
	if err != nil {
		return
	}
	switch x := raw.(type) {
	case *RawAddr__V4:
		sa.Addr = netip.AddrFrom4(*x)
		if sa.As4() == [4]byte{0, 0, 0, 0} {
			// All zero IP serves as a placeholder
			sa.Addr = netip.Addr{}
		}
	case *RawAddr__V6:
		sa.Addr = netip.AddrFrom16(*x)
	default:
		panic("unexpected RawSocketAddr")
	}
	return
}

func (a Addr) Serialize(serializer serde.Serializer) error {
	var raw RawAddr
	if a.Is4() {
		v4 := RawAddr__V4(a.As4())
		raw = &v4
	} else {
		v6 := RawAddr__V6(a.As16())
		raw = &v6
	}
	return raw.Serialize(serializer)
}
