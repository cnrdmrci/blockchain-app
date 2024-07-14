package blockchain

import (
	"blockchain-app/handlers"
	"blockchain-app/merkle"
	"blockchain-app/pow"
	"bytes"
	"encoding/gob"
	"time"
)

type Block struct {
	Timestamp    int64
	Hash         []byte
	Transactions []*Transaction
	PrevHash     []byte
	Nonce        int
	Height       int
}

func CreateGenesisBlock(tx *Transaction) *Block {
	return CreateBlock([]*Transaction{tx}, make([]byte, 32), 0)
}

func CreateBlock(txs []*Transaction, prevHash []byte, height int) *Block {
	block := &Block{time.Now().Unix(), []byte{}, txs, prevHash, 0, height}
	pow := pow.NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.Serialize())
	}
	tree := merkle.CreateMerkleTree(txHashes)

	return tree.RootNode.Data
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(b)
	handlers.HandleErrors(err)
	return res.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	handlers.HandleErrors(err)
	return &block
}
