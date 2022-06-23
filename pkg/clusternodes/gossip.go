package clusternodes

import (
	"context"
	"math/rand"
	"sync"
	"time"

	envv1 "github.com/certusone/radiance/proto/env/v1"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"k8s.io/klog/v2"
)

type Tracker struct {
	mu       sync.Mutex
	current  []*rpc.GetClusterNodesResult
	byPubkey map[solana.PublicKey]*rpc.GetClusterNodesResult
}

func New() *Tracker {
	return &Tracker{
		byPubkey: make(map[solana.PublicKey]*rpc.GetClusterNodesResult),
	}
}

// Run periodically fetches the gossip
func (t *Tracker) Run(ctx context.Context, nodes []*envv1.RPCNode, interval time.Duration) {
	t.update(ctx, nodes)
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(interval):
			t.update(ctx, nodes)
		}
	}
}

func (t *Tracker) update(ctx context.Context, nodes []*envv1.RPCNode) {
	now := time.Now()

	// Fetch gossip
	node := nodes[rand.Intn(len(nodes))]
	c := rpc.New(node.Http)
	klog.Infof("Fetching cluster nodes from %s", node.Http)
	out, err := c.GetClusterNodes(ctx)
	if err != nil {
		klog.Errorf("Failed to update cluster nodes: %v", err)
		return
	}

	klog.Infof("Fetched %d nodes in %v", len(out), time.Since(now))

	t.mu.Lock()
	t.current = out
	for _, n := range out {
		t.byPubkey[n.Pubkey] = n
	}
	t.mu.Unlock()
}

func (t *Tracker) GetByPubkey(pubkey solana.PublicKey) *rpc.GetClusterNodesResult {
	t.mu.Lock()
	entry := t.byPubkey[pubkey]
	t.mu.Unlock()
	return entry
}
