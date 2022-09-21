package clusternodes

import (
	"context"
	"math/rand"
	"sync"
	"time"

	envv1 "go.firedancer.io/radiance/proto/env/v1"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"k8s.io/klog/v2"
)

type Tracker struct {
	mu       sync.Mutex
	current  []*rpc.GetClusterNodesResult
	byPubkey map[solana.PublicKey]*rpc.GetClusterNodesResult
	c        map[string]*rpc.Client
	nodes    []*envv1.RPCNode
}

func New(nodes []*envv1.RPCNode) *Tracker {
	c := make(map[string]*rpc.Client)
	for _, node := range nodes {
		c[node.Name] = rpc.New(node.Http)
	}

	return &Tracker{
		byPubkey: make(map[solana.PublicKey]*rpc.GetClusterNodesResult),
		c:        c,
		nodes:    nodes,
	}
}

// Run periodically fetches the gossip
func (t *Tracker) Run(ctx context.Context, interval time.Duration) {
	t.update(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(interval):
			t.update(ctx)
		}
	}
}

func (t *Tracker) update(ctx context.Context) {
	now := time.Now()

	// Fetch gossip
	node := t.nodes[rand.Intn(len(t.nodes))]
	c := t.c[node.Name]
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
