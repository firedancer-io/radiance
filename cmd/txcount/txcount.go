package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/jsonrpc"
	"k8s.io/klog/v2"
	"os"
	"time"
)

const localRPC = "http://localhost:8899"

var (
	flagAddr  = flag.String("addr", localRPC, "RPC address")
	flagCount = flag.Bool("count", true, "Print the number of transactions")
	flagAfter = flag.Uint64("after", 0, "Only print transactions after this slot")
	voteAcc   = flag.String("voteAcc", "Certusm1sa411sMpV9FPqU5dXAYhmmhygvxJ23S6hJ24", "Vote account address")
)

type logLine struct {
	TS    time.Time
	Slot  uint64
	NumTx int
}

func main() {
	flag.Parse()

	r := rpc.New(*flagAddr)

	our, err := solana.PublicKeyFromBase58(*voteAcc)
	if err != nil {
		klog.Fatalf("Failed to parse vote account: %v", err)
	}

	klog.Infof("Our vote account: %v", our)

	epr, err := r.GetEpochInfo(context.Background(), rpc.CommitmentConfirmed)
	if err != nil {
		klog.Exitf("GetEpochSchedule: %v", err)
	}

	if epr == nil {
		klog.Exitf("GetEpochInfo: empty response")
	}

	epoch := *epr
	offset := epoch.AbsoluteSlot - epoch.SlotIndex
	klog.Infof("Epoch: %v, SlotIndex: %v, AbsoluteSlot: %v, Offset: %v", epoch.Epoch, epoch.SlotIndex, epoch.AbsoluteSlot, offset)

	resp, err := r.GetLeaderScheduleWithOpts(context.Background(), &rpc.GetLeaderScheduleOpts{
		Epoch:    &epoch.AbsoluteSlot,
		Identity: &our,
	})
	if err != nil {
		klog.Exitf("GetLeaderSchedule: %v", err)
	}

	if resp == nil {
		klog.Exitf("GetLeaderSchedule: empty response")
	}

	current, err := r.GetSlot(context.Background(), rpc.CommitmentConfirmed)
	if err != nil {
		klog.Exitf("GetSlot: %v", err)
	}

	klog.Infof("Current slot: %d", current)

	sched := *resp
	slots := sched[our]

	klog.Infof("%d slots for %s", len(slots), our)

	for _, slot := range slots {
		slot := slot + offset
		if *flagAfter > 0 && slot < *flagAfter {
			continue
		}
		if slot > current {
			break
		}

		if *flagCount {
			block, err := r.GetBlock(context.Background(), slot)
			if err != nil {
				var rpcErr *jsonrpc.RPCError
				if errors.As(err, &rpcErr) && (rpcErr.Code == -32007 /* SLOT_SKIPPED */ || rpcErr.Code == -32004 /* BLOCK_NOT_AVAILABLE */) {
					fmt.Printf("slot=%d skipped=true\n", slot)
					continue
				}
				klog.Warningf("GetBlock: %v", err)
				continue
			}

			bt := time.Unix(*block.BlockTime, 0)

			json.NewEncoder(os.Stdout).Encode(logLine{TS: bt, Slot: slot, NumTx: len(block.Transactions)})
		} else {
			fmt.Printf("https://explorer.solana.com/block/%d\n", slot)
		}

	}
}
