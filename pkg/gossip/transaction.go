package gossip

import (
	"github.com/gagliardetto/solana-go"
	"github.com/novifinancial/serde-reflection/serde-generate/runtime/golang/serde"
)

type Transaction solana.Transaction

func DeserializeTransaction(deserializer serde.Deserializer) (Transaction, error) {
	panic("not implemented")
}

func (obj *Transaction) Serialize(serializer serde.Serializer) error {
	panic("not implemented")
}
