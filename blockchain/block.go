package blockchain

import (
	"blockchain-app/handlers"
	"blockchain-app/merkle"
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"
	"time"
)

type Block struct {
	Timestamp    int64
	Hash         []byte
	Transactions []*Transaction
	PrevHash     []byte
	Nonce        int64
	Height       int64
}

func CreateGenesisBlock(tx *Transaction) *Block {
	return CreateBlock([]*Transaction{tx}, make([]byte, 32), 1)
}

func CreateBlock(txs []*Transaction, prevHash []byte, height int64) *Block {
	block := &Block{time.Now().Unix(), []byte{}, txs, prevHash, 0, height}
	pow := NewProof(block)
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

func (b *Block) PrintBlockDetails() {
	fmt.Printf("Hash: %x\n", b.Hash)
	fmt.Printf("Prev. hash: %x\n", b.PrevHash)
	pow := NewProof(b)
	fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
	for _, tx := range b.Transactions {
		fmt.Println(tx)
	}
	fmt.Println()
}

func (b *Block) IsGenesis() bool {
	return bytes.Compare(b.PrevHash, make([]byte, 32)) == 0
}
