package leaderschedule

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	envv1 "github.com/certusone/radiance/proto/env/v1"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"k8s.io/klog/v2"
)

type Tracker struct {
	// mu guards max and cur
	mu sync.Mutex
	// max is the current leader schedule's highest slot
	max uint64
	// cur is the current highest slot observed on the network
	cur uint64

	// bySlot maps slots to their leader. Grows indefinitely as epochs progress.
	bySlot   map[uint64]solana.PublicKey
	bySlotMu sync.Mutex

	// slotsPerEpoch is the number of slots per epoch on the network.
	// Fetched once via RPC. Used to calculate epoch boundaries.
	slotsPerEpoch uint64

	// initCh is used to signal that the leader schedule is available.
	initCh chan struct{}
}

const (
	// prefetchSlots is the number of slots to prefetch.
	prefetchSlots = 1000
)

// FirstSlot returns the epoch number and first slot of the epoch.
func (t *Tracker) FirstSlot(slotInEpoch uint64) (uint64, uint64) {
	epoch := slotInEpoch / t.slotsPerEpoch
	firstSlot := epoch * t.slotsPerEpoch
	return epoch, firstSlot
}

// Update updates the current highest slot. Non-blocking.
func (t *Tracker) Update(slot uint64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if slot > t.cur {
		t.cur = slot
	}
}

// Run periodically updates the leader schedule. When the current slot + prefetchSlots
// is greater than the highest slot in the schedule, the schedule is updated.
//
// A random node is picked from nodes to do the request against.
func (t *Tracker) Run(ctx context.Context, nodes []*envv1.RPCNode) {
	t.initCh = make(chan struct{})

	for {
		// Fetch slots per epoch
		node := nodes[rand.Intn(len(nodes))]
		c := rpc.New(node.Http)
		klog.Infof("Fetching epoch schedule from %s", node.Http)
		out, err := c.GetEpochSchedule(ctx)
		if err != nil {
			klog.Errorf("get epoch schedule: %w", err)
			time.Sleep(time.Second)
			continue
		}
		if out.FirstNormalEpoch != 0 {
			panic("first normal epoch should be 0")
		}
		if out.FirstNormalSlot != 0 {
			panic("first normal slot should be 0")
		}
		if out.LeaderScheduleSlotOffset <= prefetchSlots {
			panic("leader schedule slot offset should be greater than prefetch slots")
		}
		t.slotsPerEpoch = out.SlotsPerEpoch
		klog.Infof("Got epoch schedule: slotsPerEpoch=%d", t.slotsPerEpoch)
		break
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			t.mu.Lock()
			fetch := false

			// If we're less than prefetchSlots slots away from the epoch boundary,
			// fetch the new leader schedule.
			threshold := t.cur + prefetchSlots

			if threshold > t.max && t.cur != 0 {
				fetch = true
			}

			// If we have no current leader schedule, fetch the current one.
			var slot uint64
			var prefetch bool
			if t.max == 0 {
				slot = t.cur
			} else {
				// If we have a leader schedule, prefetch the next one
				slot = t.max + 1
				prefetch = true
			}
			t.mu.Unlock()
			if fetch {
				if prefetch {
					klog.Infof("Prefetching leader schedule for cur=%d, threshold=%d, max=%d, slot=%d",
						t.cur, threshold, t.max, slot)
					if err := t.fetch(ctx, nodes, slot); err != nil {
						klog.Errorf("Failed to fetch leader schedule: %v", err)
					}
				} else {
					klog.Infof("Fetching initial leader schedule for cur=%d, threshold=%d, max=%d, slot=%d",
						t.cur, threshold, t.max, slot)
					if err := t.fetch(ctx, nodes, slot); err != nil {
						klog.Errorf("Failed to fetch leader schedule: %v", err)
					}

					// Signal that the leader schedule is available
					close(t.initCh)
				}
			}
			time.Sleep(time.Second)
		}
	}
}

func (t *Tracker) fetch(ctx context.Context, nodes []*envv1.RPCNode, slot uint64) error {
	t.bySlotMu.Lock()
	defer t.bySlotMu.Unlock()

	if t.bySlot == nil {
		t.bySlot = make(map[uint64]solana.PublicKey)
	}

	// Pick random node from nodes
	node := nodes[rand.Intn(len(nodes))]
	klog.Infof("Using node %s", node.Http)

	// Fetch the leader schedule
	c := rpc.New(node.Http)
	out, err := c.GetLeaderScheduleWithOpts(ctx, &rpc.GetLeaderScheduleOpts{
		Epoch: &slot,
	})
	if err != nil {
		return fmt.Errorf("get leader schedule: %w", err)
	}

	epoch, firstSlot := t.FirstSlot(slot)

	// Update the leader schedule
	m := uint64(0)
	for pk, slots := range out {
		for _, s := range slots {
			t.bySlot[firstSlot+s] = pk
			if firstSlot+s > m {
				m = firstSlot + s
			}
		}
	}

	t.mu.Lock()
	t.max = m
	t.mu.Unlock()

	klog.Infof("Updated leader schedule for epoch=%d, slot=%d, first=%d, max=%d",
		epoch, slot, firstSlot, t.max)
	return nil
}

// Get returns the scheduled leader for the given slot.
// It blocks until the leader schedule is available.
func (t *Tracker) Get(slot uint64) solana.PublicKey {
	// Block until the leader schedule is available
	if t.initCh != nil {
		<-t.initCh
	}

	t.bySlotMu.Lock()
	defer t.bySlotMu.Unlock()

	return t.bySlot[slot]
}
