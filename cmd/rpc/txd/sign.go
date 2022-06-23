package main

import (
	"os"
	"path"

	"github.com/gagliardetto/solana-go"
)

func loadLocalSigner() (solana.PrivateKey, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	key := path.Join(home, ".config/solana/id.json")
	return solana.PrivateKeyFromSolanaKeygenFile(key)
}
