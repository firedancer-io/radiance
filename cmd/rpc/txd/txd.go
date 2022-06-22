package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"sync/atomic"
	"time"

	"github.com/certusone/radiance/pkg/envfile"
	"github.com/certusone/radiance/pkg/leaderschedule"
	envv1 "github.com/certusone/radiance/proto/env/v1"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"k8s.io/klog/v2"
)

var (
	flagEnv  = flag.String("env", ".env.prototxt", "Env file (.prototxt)")
	flagOnly = flag.String("only", "", "Only watch specified nodes (comma-separated)")

	flagDebugAddr = flag.String("debugAddr", "localhost:6060", "pprof/metrics listen address")
)

func init() {
	klog.InitFlags(nil)
	flag.Parse()
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

	go func() {
		klog.Error(http.ListenAndServe(*flagDebugAddr, nil))
	}()

	nodes = envfile.FilterNodes(nodes, envfile.ParseOnlyFlag(*flagOnly))

	if len(nodes) == 0 {
		klog.Exitf("No nodes in environment or all nodes filtered")
	}
	klog.Infof("Watching %d nodes", len(nodes))

	ctx := context.Background()

	// Leader schedule helper
	sched := &leaderschedule.Tracker{}
	go sched.Run(ctx, env.Nodes)

	var highest uint64

	for _, node := range nodes {
		node := node
		go func() {
			for {
				if err := watchSlotUpdates(ctx, node, sched, &highest); err != nil {
					klog.Errorf("watchSlotUpdates on node %s, reconnecting: %v", node.Name, err)
				}
				time.Sleep(time.Second * 5)
			}
		}()
	}

	select {}
}

func watchSlotUpdates(ctx context.Context, node *envv1.RPCNode, sched *leaderschedule.Tracker, highest *uint64) error {
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

		sched.Update(m.Slot)

		if m.Type == ws.SlotsUpdatesFirstShredReceived {
			klog.V(1).Infof("%s: first shred received for slot %d", node.Name, m.Slot)

			if m.Slot > atomic.LoadUint64(highest) {
				atomic.StoreUint64(highest, m.Slot)
				klog.Infof("%s: highest slot is now %d", node.Name, m.Slot)
			}
		}
	}
}
