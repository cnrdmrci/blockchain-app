package blockchain

import (
	"blockchain-app/encoders"
	"blockchain-app/handlers"
	"bytes"
	"encoding/gob"
)

type TxOutput struct {
	Value      int
	PubKeyHash []byte
}

type TxOutputs struct {
	Outputs []TxOutput
}

type TxInput struct {
	TxID      []byte
	Out       int
	Signature []byte
	PubKey    []byte
}

func (out *TxOutput) Lock(address string) {
	pubKeyHash := encoders.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}

func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

func NewTXOutput(amount int, address string) *TxOutput {
	txo := &TxOutput{amount, nil}
	txo.Lock(address)

	return txo
}

func (outs TxOutputs) Serialize() []byte {
	var buffer bytes.Buffer
	encode := gob.NewEncoder(&buffer)
	err := encode.Encode(outs)
	handlers.HandleErrors(err)
	return buffer.Bytes()
}

func DeserializeOutputs(data []byte) TxOutputs {
	var outputs TxOutputs
	decode := gob.NewDecoder(bytes.NewReader(data))
	err := decode.Decode(&outputs)
	handlers.HandleErrors(err)
	return outputs
}
