package blockchain

import (
	"blockchain-app/handlers"
	"blockchain-app/merkle"
	"bytes"
	"encoding/gob"
	"github.com/fatih/color"
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
	PrintBlockDivider()
	color.Green("Block; Height: %d, Nonce: %d, Timestamp: %d, PoW: %s", b.Height, b.Nonce, b.Timestamp, strconv.FormatBool(NewProof(b).Validate()))
	color.Yellow(getSpace(2)+"- Block Hash: %x\n", b.Hash)
	color.Yellow(getSpace(2)+"- Prev. Hash: %x\n", b.PrevHash)
	for _, tx := range b.Transactions {
		tx.Print()
	}
}

func PrintBlockDivider() {
	color.Red("---------------------------------------------------------------------------------------")
}

func getSpace(count int) string {
	space := ""
	for i := 0; i < count; i++ {
		space += " "
	}
	return space
}

func (b *Block) IsGenesis() bool {
	return bytes.Compare(b.PrevHash, make([]byte, 32)) == 0
}
