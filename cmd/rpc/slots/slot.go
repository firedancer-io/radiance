package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"sync"
	"time"

	"go.firedancer.io/radiance/pkg/envfile"
	"go.firedancer.io/radiance/pkg/kafka"
	"go.firedancer.io/radiance/pkg/leaderschedule"
	envv1 "go.firedancer.io/radiance/proto/env/v1"
	networkv1 "go.firedancer.io/radiance/proto/network/v1"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/golang/protobuf/proto"
	"github.com/twmb/franz-go/pkg/kgo"
	"k8s.io/klog/v2"
)

var (
	flagEnv  = flag.String("env", ".env.prototxt", "Env file (.prototxt)")
	flagOnly = flag.String("only", "", "Only watch specified nodes (comma-separated)")

	flagType = flag.String("type", "", "Only print specific types to log")

	flagKafka      = flag.Bool("kafka", false, "Enable Kafka publishing")
	flagKafkaTopic = flag.String("kafkaTopic", "slot_status", "Kafka topic suffix to publish to")

	flagDebugAddr = flag.String("debugAddr", "localhost:6060", "pprof/metrics listen address")
)

func init() {
	klog.InitFlags(nil)
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

	highest := &sync.Map{}

	sched := &leaderschedule.Tracker{}

	go sched.Run(ctx, env.Nodes)

	var kcl *kgo.Client
	var topic string
	if *flagKafka {
		if *flagKafkaTopic == "" {
			klog.Exitf("Kafka enabled but no topic specified")
		}
		topic = strings.Join([]string{env.Kafka.TopicPrefix, *flagKafkaTopic}, ".")
		klog.Infof("Publishing to topic %s", topic)
		kcl, err = kafka.NewClientFromEnv(env.Kafka)
		if err != nil {
			klog.Exitf("Failed to create kafka client: %v", err)
		}
	}

	for _, node := range nodes {
		node := node
		go func() {
			for {
				if err := watchSlotUpdates(ctx, node, highest, sched, kcl, topic); err != nil {
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

func watchSlotUpdates(ctx context.Context, node *envv1.RPCNode, highest *sync.Map, sched *leaderschedule.Tracker, kcl *kgo.Client, topic string) error {
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
		delay := time.Since(ts)

		sched.Update(m.Slot)

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
			prop = ts.Sub(first).Milliseconds()
		} else {
			prop = -1
		}

		if *flagType != "" && string(m.Type) != *flagType {
			continue
		}

		leader := sched.Get(m.Slot)

		if kcl != nil {
			var stats *networkv1.TxStats
			if m.Type == ws.SlotsUpdatesFrozen {
				stats = &networkv1.TxStats{
					NumTransactionEntries:     m.Stats.NumTransactionEntries,
					NumSuccessfulTransactions: m.Stats.NumSuccessfulTransactions,
					NumFailedTransactions:     m.Stats.NumFailedTransactions,
					MaxTransactionsPerEntry:   m.Stats.MaxTransactionsPerEntry,
				}
			}

			st := &networkv1.SlotStatus{
				Slot:      m.Slot,
				Timestamp: uint64(ts.UnixMilli()),
				Delay:     uint64(delay.Milliseconds()),
				Type:      convertUpdateType(m.Type),
				Parent:    m.Parent,
				Stats:     stats,
				Err:       "", // TODO
				Leader:    leader.String(),
				Source:    node.Name,
			}

			// Fixed-length proto encoding
			buf := proto.NewBuffer([]byte{})
			if err := buf.EncodeMessage(st); err != nil {
				panic(err)
			}

			r := &kgo.Record{Topic: topic, Value: buf.Bytes()}
			kcl.Produce(ctx, r, func(_ *kgo.Record, err error) {
				if err != nil {
					klog.Warningf("failed to publish message to %s: %v", topic, err)
				}
			})
		}

		klog.V(1).Infof("%s: slot=%d type=%s delay=%dms prop=%dms parent=%d stats=%v leader=%s",
			node.Name, m.Slot, m.Type, delay.Milliseconds(), prop, m.Parent, m.Stats, leader)
	}
}

func convertUpdateType(t ws.SlotsUpdatesType) networkv1.SlotStatus_UpdateType {
	switch t {
	case ws.SlotsUpdatesFirstShredReceived:
		return networkv1.SlotStatus_UPDATE_TYPE_FIRST_SHRED_RECEIVED
	case ws.SlotsUpdatesCompleted:
		return networkv1.SlotStatus_UPDATE_TYPE_COMPLETED
	case ws.SlotsUpdatesCreatedBank:
		return networkv1.SlotStatus_UPDATE_TYPE_CREATED_BANK
	case ws.SlotsUpdatesFrozen:
		return networkv1.SlotStatus_UPDATE_TYPE_FROZEN
	case ws.SlotsUpdatesDead:
		return networkv1.SlotStatus_UPDATE_TYPE_DEAD
	case ws.SlotsUpdatesOptimisticConfirmation:
		return networkv1.SlotStatus_UPDATE_TYPE_OPTIMISTIC_CONFIRMATION
	case ws.SlotsUpdatesRoot:
		return networkv1.SlotStatus_UPDATE_TYPE_ROOT
	default:
		panic("unknown slot update type " + t)
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

		klog.V(1).Infof("%s: slot=%d root=%d parent=%d",
			node.Name, m.Slot, m.Root, m.Parent)
	}
}
