package nftables

import (
	"fmt"
	"github.com/google/nftables"
	"github.com/google/nftables/binaryutil"
	"github.com/google/nftables/expr"
	"golang.org/x/sys/unix"
	"os"
	"os/exec"
)

func getNFTConn() nftables.Conn {
	return nftables.Conn{}
}

// EnsureKernelModules ensures that nft_tproxy kernel module is loaded.
func EnsureKernelModules() error {
	// Check if directory /sys/module/nft_tproxy exists
	_, err := os.Stat("/sys/module/nft_tproxy")
	if err == nil {
		return nil
	}

	// If not, run "modprobe nft_tproxy"
	_, err = exec.Command("modprobe", "nft_tproxy").Output()
	if err != nil {
		return err
	}

	return nil
}

func ifname(n string) []byte {
	b := make([]byte, 16)
	copy(b, []byte(n+"\x00"))
	return b
}

// InsertProxyChain create a new nftables table with a single chain with a tproxy rule in it.
//
// We create a new chain with a custom priority to ensure that we can cleanly add and remove it
// without affecting other rules or knowing the setup.
//
// Generates a chain that looks like this:
//
// table ip filter {
//   chain tpuproxy {
//     type filter hook prerouting priority -2147483648; policy accept;
//     iifname "enp6s0" udp dport { 8003, 8004, 8005 } tproxy to :51211
//   }
// }
//
func InsertProxyChain(destPorts []uint16, redirectPort uint16, iface string) error {
	c := getNFTConn()

	// Create a new table just in case the host isn't using nft.
	table := c.AddTable(&nftables.Table{
		Name:   "filter",
		Family: nftables.TableFamilyIPv4,
	})

	// Create our own chain so we can cleanly remove it later.
	chain := c.AddChain(&nftables.Chain{
		Name:     "tpuproxy",
		Table:    table,
		Type:     nftables.ChainTypeFilter,
		Hooknum:  nftables.ChainHookPrerouting,
		Priority: nftables.ChainPriorityFirst,
	})

	c.FlushChain(chain)

	set := &nftables.Set{
		Anonymous: true,
		Constant:  true,
		Table:     table,
		KeyType:   nftables.TypeInetService,
	}

	elements := make([]nftables.SetElement, len(destPorts))
	for i, port := range destPorts {
		elements[i] = nftables.SetElement{
			Key: binaryutil.BigEndian.PutUint16(port),
		}
	}
	if err := c.AddSet(set, elements); err != nil {
		return fmt.Errorf("failed to add set: %v", err)
	}

	// Add a rule to the chain that matches on the destination port and redirects to the redirect port.
	c.AddRule(&nftables.Rule{
		Table: table,
		Chain: chain,
		Exprs: []expr.Any{
			// [ meta load iifname => reg 1 ]
			&expr.Meta{Key: expr.MetaKeyIIF, Register: 1},
			// [ cmp eq reg 1 0x696c7075 0x00306b6e 0x00000000 0x00000000 ]
			&expr.Cmp{Op: expr.CmpOpNeq, Register: 1, Data: binaryutil.NativeEndian.PutUint32(uint32(1)) /* lo */}, // TODO
			// [ payload load 1b @ network header + 9 => reg 1 ]
			&expr.Payload{DestRegister: 1, Base: expr.PayloadBaseNetworkHeader, Offset: 9, Len: 1},
			// [ cmp eq reg 1 <proto> ]
			&expr.Cmp{Op: expr.CmpOpEq, Register: 1, Data: []byte{unix.IPPROTO_UDP}},
			// [ payload load 2b @ transport header + 2 => reg 1 ]
			&expr.Payload{DestRegister: 1, Base: expr.PayloadBaseTransportHeader, Offset: 2, Len: 2},
			// [ lookup reg 1 set __set%d ]
			&expr.Lookup{SourceRegister: 1, SetName: set.Name, SetID: set.ID},
			//	[ immediate reg 1 <redirectPort> ]
			&expr.Immediate{Register: 1, Data: binaryutil.BigEndian.PutUint16(redirectPort)},
			//	[ tproxy ip port reg 1 ]
			&expr.TProxy{
				Family:      byte(nftables.TableFamilyIPv4),
				TableFamily: byte(nftables.TableFamilyIPv4),
				RegPort:     1,
			},
		},
	})

	if err := c.Flush(); err != nil {
		return err
	}

	return nil
}

func DeleteProxyChain() error {
	c := getNFTConn()

	// Delete the chain.
	chain := c.AddChain(&nftables.Chain{
		Name: "tpuproxy",
		Table: &nftables.Table{
			Name:   "filter",
			Family: nftables.TableFamilyIPv4,
		},
	})
	c.DelChain(chain)

	if err := c.Flush(); err != nil {
		return err
	}

	return nil
}
