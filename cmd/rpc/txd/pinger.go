package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"syscall"
	"time"

	envv1 "github.com/certusone/radiance/proto/env/v1"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/gagliardetto/solana-go/text"
	"k8s.io/klog/v2"
)

var (
	flagPingerLog = flag.String("pinger_log", "", "JSON log file for pinger")
)

const (
	confirmPollInterval = time.Second * 5
	confirmRetries      = 10
)

type pingData struct {
	Slot uint64    `json:"slot"`
	Ts   time.Time `json:"ts"`
}

type logEntry struct {
	Slot      uint64           `json:"slot"`
	SendDelay time.Duration    `json:"send_delay"`
	Ts        time.Time        `json:"ts"`
	Signature solana.Signature `json:"signature"`
	Leader    solana.PublicKey `json:"leader"`

	Confirmed     bool          `json:"confirmed"`
	ConfirmedSlot uint64        `json:"confirmed_slot,omitempty"`
	SlotDelay     *int          `json:"slot_delay,omitempty"`
	TimeDelay     time.Duration `json:"time_delay,omitempty"`

	Timeout time.Duration `json:"timeout,omitempty"`
}

func buildTransaction(slot uint64, now time.Time, blockhash solana.Hash, feePayer solana.PublicKey) *solana.Transaction {
	payload := &pingData{Slot: slot, Ts: now}
	b, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	ins := solana.NewInstruction(solana.MemoProgramID, solana.AccountMetaSlice{}, b)

	tx, err := solana.NewTransaction(
		[]solana.Instruction{ins}, blockhash, solana.TransactionPayer(feePayer))
	if err != nil {
		panic(err)
	}

	return tx
}

func sendPing(ctx context.Context, m *ws.SlotsUpdatesResult, b solana.Hash, signer solana.PrivateKey, g *rpc.GetClusterNodesResult, nodes []*envv1.RPCNode, c map[string]*rpc.Client) {
	tx := buildTransaction(m.Slot, time.Now(), b, signer.PublicKey())
	sigs, err := tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key != signer.PublicKey() {
			panic("no private key for unknown signer " + key.String())
		}
		return &signer
	})
	if err != nil {
		panic(err)
	}

	if klog.V(2).Enabled() {
		tx.EncodeTree(text.NewTreeEncoder(os.Stdout, "Ping memo"))
	}

	txb, err := tx.MarshalBinary()
	if err != nil {
		panic(err)
	}

	klog.Infof("Sending tx %s", sigs[0].String())
	klog.V(2).Infof("tx: %s", hex.EncodeToString(txb))

	sendUDP(*g.TPU, txb, 20)

	go waitForConfirmation(ctx, m, nodes, c, sigs, g)
}

func waitForConfirmation(ctx context.Context, m *ws.SlotsUpdatesResult, nodes []*envv1.RPCNode, c map[string]*rpc.Client, sigs []solana.Signature, g *rpc.GetClusterNodesResult) {
	for i := 0; i < confirmRetries; i++ {
		select {
		case <-ctx.Done():
			return
		case <-time.After(confirmPollInterval):
			// pick random node to query status
			node := nodes[rand.Intn(len(nodes))]
			st, err := c[node.Name].GetSignatureStatuses(ctx, false, sigs[0])
			if err != nil {
				klog.Errorf("Failed to fetch signature status: %v", err)
				continue
			}

			for _, s := range st.Value {
				if s == nil || s.ConfirmationStatus != rpc.ConfirmationStatusConfirmed {
					continue
				}
				delay := int(s.Slot) - int(m.Slot)
				klog.Infof("%s confirmed in slot %d (offset %d)", sigs[0], s.Slot, delay)
				if *flagPingerLog != "" {
					log(logEntry{
						Slot:          m.Slot,
						SendDelay:     time.Duration(0),
						Ts:            time.Now(),
						Signature:     sigs[0],
						Leader:        g.Pubkey,
						Confirmed:     true,
						ConfirmedSlot: s.Slot,
						SlotDelay:     &delay,
						TimeDelay:     0, // TODO - we need precise slot timings for this first
					})
					return
				}
			}
		}
	}
	klog.Infof("%s failed to confirm after %v",
		sigs[0].String(), confirmPollInterval*confirmRetries)

	if *flagPingerLog != "" {
		log(logEntry{
			Slot:      m.Slot,
			SendDelay: time.Duration(0),
			Ts:        time.Now(),
			Signature: sigs[0],
			Leader:    g.Pubkey,
			Confirmed: false,
			Timeout:   confirmPollInterval * confirmRetries,
		})
	}
}

func log(entry logEntry) {
	f, err := os.OpenFile(*flagPingerLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		klog.Errorf("failed to open %s: %v", *flagPingerLog, err)
		return
	}
	defer f.Close()

	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		panic(fmt.Sprintf("failed to lock %s: %v", *flagPingerLog, err))
	}

	if err := json.NewEncoder(f).Encode(entry); err != nil {
		klog.Errorf("failed to write to %s: %v", *flagPingerLog, err)
	}

	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_UN); err != nil {
		panic(fmt.Sprintf("failed to unlock %s: %v", *flagPingerLog, err))
	}
}
