package tpu

import (
	"crypto/ed25519"
	"errors"
	"fmt"
	"github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
)

func ParseTx(p []byte) (tx *solana.Transaction, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("ParseTx panic")
		}
	}()

	tx, err = solana.TransactionFromDecoder(bin.NewBinDecoder(p))
	if err != nil {
		return nil, err
	}
	return
}

func VerifyTxSig(tx *solana.Transaction) (ok bool) {
	msg, err := tx.Message.MarshalBinary()
	if err != nil {
		panic(err)
	}

	signers := ExtractSigners(tx)

	if len(signers) != len(tx.Signatures) {
		return false
	}

	for i, sig := range tx.Signatures {
		if !ed25519.Verify(signers[i][:], msg, sig[:]) {
			fmt.Printf("invalid signature by %s\n", signers[i].String())
			return false
		}
	}

	return true
}

func ExtractSigners(tx *solana.Transaction) []solana.PublicKey {
	signers := make([]solana.PublicKey, 0, len(tx.Signatures))
	for _, acc := range tx.Message.AccountKeys {
		if tx.IsSigner(acc) {
			signers = append(signers, acc)
		}
	}
	return signers
}
