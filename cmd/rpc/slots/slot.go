package main

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/certusone/radiance/pkg/envfile"
	"github.com/certusone/radiance/proto/envv1"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"k8s.io/klog/v2"
)

var (
	flagEnv  = flag.String("env", ".env.prototxt", "Env file (.prototxt)")
	flagOnly = flag.String("only", "", "Only watch specified nodes (comma-separated)")
	flagType = flag.String("type", "", "Only print specific types")
)

func init() {
	flag.Parse()
}

/*
	I0612 20:37:10.076826  916547 slot.go:111] val1.ffm1: slot=137326466 type=firstShredReceived delta=7ms parent=0
	I0612 20:37:10.428919  916547 slot.go:111] val1.ffm1: slot=137326466 type=completed delta=7ms parent=0
	I0612 20:37:10.687256  916547 slot.go:111] val1.ffm1: slot=137326466 type=createdBank delta=4ms parent=137326465
	I0612 20:37:10.691104  916547 slot.go:136] val1.ffm1: slot=137326466 root=137326431 parent=137326465
	I0612 20:37:11.232413  916547 slot.go:111] val1.ffm1: slot=137326466 type=frozen delta=8ms parent=0
	I0612 20:37:12.480333  916547 slot.go:111] val1.ffm1: slot=137326466 type=optimisticConfirmation delta=8ms parent=0
	I0612 20:37:43.279139  916547 slot.go:111] val1.ffm1: slot=137326466 type=root delta=9ms parent=0
	I0612 20:37:43.805364  916547 slot.go:111] val1.ffm1: slot=137326466 type=root delta=8ms parent=0
*/

func parseOnlyFlag(only string) []string {
	if only == "" {
		return nil
	}
	return strings.Split(only, ",")
}

func filterNodes(nodes []*envv1.RPCNode, only []string) []*envv1.RPCNode {
	if len(only) == 0 {
		return nodes
	}
	var filtered []*envv1.RPCNode
	for _, node := range nodes {
		for _, o := range only {
			if node.Name == o {
				filtered = append(filtered, node)
			}
		}
	}
	return filtered
}

func main() {
	env, err := envfile.Load(*flagEnv)
	if err != nil {
		klog.Fatalf("Failed to load env file: %v", err)
	}

	nodes := env.GetNodes()
	if len(nodes) == 0 {
		klog.Fatalf("No nodes found in env file")
	}

	nodes = filterNodes(nodes, parseOnlyFlag(*flagOnly))

	klog.Infof("Watching %d nodes", len(nodes))

	ctx := context.Background()

	highest := &sync.Map{}

	for _, node := range nodes {
		node := node
		go func() {
			for {
				if err := watchSlotUpdates(ctx, node, highest); err != nil {
					klog.Errorf("watchSlotUpdates on node %s, reconnecting: %v", node.Name, err)
				}
				time.Sleep(time.Second * 5)
			}
		}()

		if *flagType == "" {
			go func() {
				for {
					if err := watchSlots(ctx, node); err != nil {
						klog.Errorf("watchSlots on node %s, reconnecting: %v", node.Name, err)
					}
					time.Sleep(time.Second * 5)
				}
			}()
		}
	}

	select {}
}

func watchSlotUpdates(ctx context.Context, node *envv1.RPCNode, highest *sync.Map) error {
	timeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	c, err := ws.Connect(timeout, node.Ws)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	sub, err := c.SlotsUpdatesSubscribe()
	if err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}

	for {
		m, err := sub.Recv()
		if err != nil {
			return fmt.Errorf("recv: %w", err)
		}

		ts := m.Timestamp.Time()
		delta := time.Since(ts)

		if *flagType != "" && string(m.Type) != *flagType {
			continue
		}

		var first time.Time
		if m.Type == ws.SlotsUpdatesFirstShredReceived {
			value, _ := highest.LoadOrStore(m.Slot, ts)
			first = value.(time.Time)
		} else {
			value, ok := highest.Load(m.Slot)
			if ok {
				first = value.(time.Time)
			}
		}

		if m.Type == ws.SlotsUpdatesRoot {
			highest.Delete(m.Slot)
		}

		var prop int64
		if !first.IsZero() {
			prop = time.Since(first).Milliseconds()
		} else {
			prop = -1
		}

		klog.Infof("%s: slot=%d type=%s delta=%dms prop=%dms parent=%d stats=%v",
			node.Name, m.Slot, m.Type, delta.Milliseconds(), prop, m.Parent, m.Stats)
	}
}

func watchSlots(ctx context.Context, node *envv1.RPCNode) error {
	timeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	c, err := ws.Connect(timeout, node.Ws)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	sub, err := c.SlotSubscribe()
	if err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}

	for {
		m, err := sub.Recv()
		if err != nil {
			return fmt.Errorf("recv: %w", err)
		}

		klog.Infof("%s: slot=%d root=%d parent=%d",
			node.Name, m.Slot, m.Root, m.Parent)
	}
}
