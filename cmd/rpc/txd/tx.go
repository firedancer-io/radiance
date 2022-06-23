package main

import (
	"encoding/json"
	"time"

	"github.com/gagliardetto/solana-go"
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
