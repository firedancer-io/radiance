package endpoints

import (
	"context"
	"fmt"
	"github.com/gagliardetto/solana-go/rpc"
	"net"
	"strconv"
)

const RPCLocalhost = "http://127.0.0.1:8899"

func GetNodeTPUPorts(ctx context.Context, rpcHost string, nodeIP net.IP) ([]uint16, error) {
	c := rpc.New(rpcHost)

	out, err := c.GetClusterNodes(ctx)
	if err != nil {
		return nil, err
	}

	for _, node := range out {
		if node.TPU == nil {
			continue
		}
		tpuAddr := *node.TPU
		host, port, err := net.SplitHostPort(tpuAddr)
		if err != nil {
			return nil, fmt.Errorf("error parsing node TPU %s: %v", tpuAddr, err)
		}
		if host == nodeIP.String() {
			port, err := strconv.Atoi(port)
			if err != nil {
				return nil, fmt.Errorf("error parsing node TPU %s: %v", tpuAddr, err)
			}
			return []uint16{
				uint16(port),     // TPU
				uint16(port + 1), // TPUfwd
				uint16(port + 2), // TPUvote
			}, nil
		}
	}

	return nil, fmt.Errorf("node %s not found in cluster", nodeIP)
}
