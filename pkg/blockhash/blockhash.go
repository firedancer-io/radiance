package blockhash

import (
	"context"
	"sync"
	"time"

	envv1 "github.com/certusone/radiance/proto/env/v1"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"k8s.io/klog/v2"
)

type Tracker struct {
	mu     sync.Mutex
	byNode map[string]struct {
		Blockhash        solana.Hash
		HighestValidSlot uint64
	}
	nodes []*envv1.RPCNode

	c map[string]*rpc.Client
}

func New(nodes []*envv1.RPCNode) *Tracker {
	t := &Tracker{
		byNode: make(map[string]struct {
			Blockhash        solana.Hash
			HighestValidSlot uint64
		}),
		c:     make(map[string]*rpc.Client),
		nodes: nodes,
	}

	for _, node := range nodes {
		t.c[node.Name] = rpc.New(node.Http)
	}

	return t
}

func (t *Tracker) MostPopular() solana.Hash {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Return the most frequently occurring value of t.byNode
	var mostPopular solana.Hash
	var mostPopularCount int
	for _, h := range t.byNode {
		count := 0
		for _, h2 := range t.byNode {
			if h.Blockhash.Equals(h2.Blockhash) {
				count++
			}
		}
		if count > mostPopularCount {
			mostPopular = h.Blockhash
			mostPopularCount = count
		}

		klog.V(2).Infof("%s: %d", h, count)
	}
	return mostPopular
}

func (t *Tracker) MostRecent() solana.Hash {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Return the blockhash which has the highest valid slot
	var mostRecent solana.Hash
	var highestValidSlot uint64
	for _, h := range t.byNode {
		if h.HighestValidSlot > highestValidSlot {
			highestValidSlot = h.HighestValidSlot
			mostRecent = h.Blockhash
		}
	}
	return mostRecent
}

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

	for _, node := range t.nodes {
		node := node
		go func() {
			ctx, cancel := context.WithTimeout(ctx, time.Second*5)
			defer cancel()
			klog.V(1).Infof("Fetching blockhash from %s", node.Http)
			h, err := t.c[node.Name].GetLatestBlockhash(ctx, rpc.CommitmentConfirmed)
			if err != nil {
				klog.Errorf("%s: failed to request blockhash: %v", node.Name, err)
				return
			}

			klog.V(1).Infof("%s: fetched blockhash %d -> %s in %v",
				node.Name, h.Value.LastValidBlockHeight, h.Value.Blockhash, time.Since(now))
			t.mu.Lock()
			t.byNode[node.Name] = struct {
				Blockhash        solana.Hash
				HighestValidSlot uint64
			}{
				Blockhash:        h.Value.Blockhash,
				HighestValidSlot: h.Value.LastValidBlockHeight,
			}
			t.mu.Unlock()
		}()
	}
}
