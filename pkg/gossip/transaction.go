package gossip

import (
	"github.com/gagliardetto/solana-go"
	"github.com/novifinancial/serde-reflection/serde-generate/runtime/golang/serde"
)

// TODO write codegen for this

type Transaction solana.Transaction

func DeserializeTransaction(deserializer serde.Deserializer) (Transaction, error) {
	var obj Transaction
	numSigs, err := deserializer.DeserializeU8()
	if err != nil {
		return obj, err
	}
	for i := uint8(0); i < numSigs; i++ {
		sig, err := DeserializeSignature(deserializer)
		if err != nil {
			return obj, err
		}
		obj.Signatures = append(obj.Signatures, solana.Signature(sig))
	}
	obj.Message, err = DeserializeTxMessage(deserializer)
	if err != nil {
		return obj, err
	}
	return obj, nil
}

func (obj *Transaction) Serialize(serializer serde.Serializer) error {
	panic("not implemented")
}

func DeserializeTxMessage(deserializer serde.Deserializer) (solana.Message, error) {
	var obj solana.Message
	numRequiredSigs, err := deserializer.DeserializeU8()
	if err != nil {
		return obj, err
	}
	obj.Header.NumRequiredSignatures = numRequiredSigs
	numReadonlySignedAccs, err := deserializer.DeserializeU8()
	if err != nil {
		return obj, err
	}
	obj.Header.NumReadonlySignedAccounts = numReadonlySignedAccs
	numReadonlyUnsignedAccs, err := deserializer.DeserializeU8()
	if err != nil {
		return obj, err
	}
	obj.Header.NumReadonlyUnsignedAccounts = numReadonlyUnsignedAccs
	numAccountKeys, err := deserializer.DeserializeU8()
	if err != nil {
		return obj, err
	}
	for i := uint8(0); i < numAccountKeys; i++ {
		address, err := DeserializePubkey(deserializer)
		if err != nil {
			return obj, err
		}
		obj.AccountKeys = append(obj.AccountKeys, solana.PublicKey(address))
	}
	recentBlockHash, err := DeserializeHash(deserializer)
	if err != nil {
		return obj, err
	}
	obj.RecentBlockhash = solana.Hash(recentBlockHash)
	numInsns, err := deserializer.DeserializeU8()
	if err != nil {
		return obj, err
	}
	for i := uint8(0); i < numInsns; i++ {
		insn, err := DeserializeInstruction(deserializer)
		if err != nil {
			return obj, err
		}
		obj.Instructions = append(obj.Instructions, insn)
	}
	return obj, nil
}

func DeserializeInstruction(deserializer serde.Deserializer) (solana.CompiledInstruction, error) {
	var obj solana.CompiledInstruction
	programIdIdx, err := deserializer.DeserializeU8()
	if err != nil {
		return obj, err
	}
	obj.ProgramIDIndex = uint16(programIdIdx)
	numAccs, err := deserializer.DeserializeU8()
	if err != nil {
		return obj, err
	}
	for i := uint8(0); i < numAccs; i++ {
		idx, err := deserializer.DeserializeU8()
		if err != nil {
			return obj, err
		}
		obj.Accounts = append(obj.Accounts, uint16(idx))
	}
	// This is brain-dead
	dataLen, err := deserializer.DeserializeU8()
	if err != nil {
		return obj, err
	}
	obj.Data = make([]byte, dataLen)
	for i := uint8(0); i < dataLen; i++ {
		_byte, err := deserializer.DeserializeU8()
		if err != nil {
			return obj, err
		}
		obj.Data[i] = _byte
	}
	return obj, nil
}
