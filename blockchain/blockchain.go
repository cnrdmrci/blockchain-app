package blockchain

import (
	"blockchain-app/database"
	"blockchain-app/handlers"
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"runtime"
)

type BlockChain struct {
	LastHash []byte
}

func InitBlockChain(rewardAddress, nodeID string) *BlockChain {
	lastBlockHash := database.Get([]byte(lastBlockHashKey), nodeID)
	if lastBlockHash != nil {
		fmt.Println("Blockchain already exists")
		runtime.Goexit()
	}

	cbtx := CoinbaseTx(rewardAddress, "First Transaction from Genesis")
	genesis := CreateGenesisBlock(cbtx)
	database.Set(genesis.Hash, genesis.Serialize(), nodeID)
	database.Set([]byte(lastBlockHashKey), genesis.Hash, nodeID)

	blockchain := BlockChain{genesis.Hash}
	return &blockchain
}

func GetLastBlock(nodeID string) *Block {
	lastHash := database.Get([]byte(lastBlockHashKey), nodeID)
	blockByte := database.Get(lastHash, nodeID)
	return Deserialize(blockByte)
}

func (b *Block) GetPreviousBlock(nodeID string) *Block {
	prevBlock := database.Get(b.PrevHash, nodeID)
	return Deserialize(prevBlock)
}

func (c *BlockChain) AddBlock(block *Block, nodeID string) {
	blockData := block.Serialize()
	database.Set(block.Hash, blockData, nodeID)
	lastBlockHash := database.Get([]byte(lastBlockHashKey), nodeID)
	lastBlockData := database.Get(lastBlockHash, nodeID)
	lastBlock := Deserialize(lastBlockData)
	if block.Height > lastBlock.Height {
		database.Set([]byte(lastBlockHashKey), block.Hash, nodeID)
		c.LastHash = block.Hash
	}
}

func GetBestHeight(nodeID string) int {
	lastBlockHash := database.Get([]byte(lastBlockHashKey), nodeID)
	lastBlockData := database.Get(lastBlockHash, nodeID)
	lastBlock := Deserialize(lastBlockData)
	return lastBlock.Height
}

func GetBlock(blockHash []byte, nodeID string) Block {
	blockData := database.Get(blockHash, nodeID)
	if blockData == nil {
		handlers.HandleErrors(errors.New("block is not found"))
	}
	block := Deserialize(blockData)

	return *block
}

func GetBlockHashes(nodeID string) [][]byte {
	var blocks [][]byte

	block := GetLastBlock(nodeID)
	blocks = append(blocks, block.Hash)
	for {
		block = block.GetPreviousBlock(nodeID)
		blocks = append(blocks, block.Hash)

		if bytes.Compare(block.PrevHash, make([]byte, 32)) == 0 {
			break
		}
	}

	return blocks
}

func (c *BlockChain) MineBlock(transactions []*Transaction, nodeID string) *Block {
	for _, tx := range transactions {
		if c.VerifyTransaction(tx) != true {
			handlers.HandleErrors(errors.New("invalid transaction"))
		}
	}

	lastBlock := GetLastBlock(nodeID)
	lastHeight := lastBlock.Height

	newBlock := CreateBlock(transactions, lastBlock.Hash, lastHeight+1)
	c.AddBlock(newBlock, nodeID)

	return newBlock
}

// todo
func (c *BlockChain) FindUTXO(nodeID string) map[string]TxOutputs {
	UTXO := make(map[string]TxOutputs)
	spentTXOs := make(map[string][]int)

	block := GetLastBlock(nodeID)

	for {
		block = block.GetPreviousBlock(nodeID)

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				outs := UTXO[txID]
				outs.Outputs = append(outs.Outputs, out)
				UTXO[txID] = outs
			}
			if tx.IsCoinbase() == false {
				for _, in := range tx.Inputs {
					inTxID := hex.EncodeToString(in.ID)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
				}
			}
		}

		if bytes.Compare(block.PrevHash, make([]byte, 32)) == 0 {
			break
		}
	}
	return UTXO
}

func (c *BlockChain) FindTransaction(ID []byte, nodeID string) (Transaction, error) {
	block := GetLastBlock(nodeID)
	for _, tx := range block.Transactions {
		if bytes.Compare(tx.ID, ID) == 0 {
			return *tx, nil
		}
	}
	if !block.IsGenesis() {
		for {
			block = block.GetPreviousBlock(nodeID)

			for _, tx := range block.Transactions {
				if bytes.Compare(tx.ID, ID) == 0 {
					return *tx, nil
				}
			}

			if bytes.Compare(block.PrevHash, make([]byte, 32)) == 0 {
				break
			}
		}
	}

	return Transaction{}, errors.New("transaction does not exist")
}

func (c *BlockChain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey, nodeID string) {
	prevTXs := make(map[string]Transaction)

	for _, in := range tx.Inputs {
		prevTX, err := c.FindTransaction(in.ID, nodeID)
		handlers.HandleErrors(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	tx.Sign(privKey, prevTXs)
}

func (c *BlockChain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}
	prevTXs := make(map[string]Transaction)

	for _, in := range tx.Inputs {
		prevTX, err := c.FindTransaction(in.ID, "3000")
		handlers.HandleErrors(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	return tx.Verify(prevTXs)
}
