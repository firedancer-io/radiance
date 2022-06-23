package main

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/gagliardetto/solana-go/text"
	"k8s.io/klog/v2"
)

type pingData struct {
	Slot uint64    `json:"slot"`
	Ts   time.Time `json:"ts"`
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

func sendPing(m *ws.SlotsUpdatesResult, b solana.Hash, signer solana.PrivateKey, g *rpc.GetClusterNodesResult) {
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

	sendUDP(*g.TPU, txb)
}
