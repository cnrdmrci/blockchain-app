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

func InitBlockChain(rewardAddress, nodeID string) {
	lastBlockHash := database.Get([]byte(lastBlockHashKey), nodeID)
	if lastBlockHash != nil {
		fmt.Println("Blockchain already exists")
		runtime.Goexit()
	}

	cbtx := CreateCoinbaseTx(rewardAddress)
	genesis := CreateGenesisBlock(cbtx)
	database.Set(genesis.Hash, genesis.Serialize(), nodeID)
	database.Set([]byte(lastBlockHashKey), genesis.Hash, nodeID)
}

func GetLastBlock(nodeID string) *Block {
	lastHash := database.Get([]byte(lastBlockHashKey), nodeID)
	if lastHash == nil {
		fmt.Println("Blockchain not created.")
		runtime.Goexit()
	}
	blockByte := database.Get(lastHash, nodeID)
	return Deserialize(blockByte)
}

func IsBlockchainExist(nodeID string) bool {
	lastHash := database.Get([]byte(lastBlockHashKey), nodeID)
	if lastHash == nil {
		return false
	}
	return true
}

func (b *Block) GetPreviousBlock(nodeID string) *Block {
	if b.IsGenesis() {
		return nil
	}
	prevBlock := database.Get(b.PrevHash, nodeID)
	return Deserialize(prevBlock)
}

func (b *Block) GetSTXO(currentSTXO map[string][]int) map[string][]int {
	STXO := make(map[string][]int)
	if currentSTXO != nil {
		STXO = currentSTXO
	}

	for _, tx := range b.Transactions {

		if !tx.IsCoinbase() {
			for _, in := range tx.Inputs {
				inTxID := hex.EncodeToString(in.TxID)
				STXO[inTxID] = append(STXO[inTxID], in.Out)
			}
		}
	}
	return STXO
}

func (b *Block) GetUTXO(STXO map[string][]int, currentUTXO map[string]TxOutputs) map[string]TxOutputs {
	UTXO := make(map[string]TxOutputs)
	if currentUTXO != nil {
		UTXO = currentUTXO
	}

	for _, tx := range b.Transactions {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Outputs {
			if STXO[txID] != nil {
				for _, spentOut := range STXO[txID] {
					if spentOut != outIdx {
						outs := UTXO[txID]
						outs.Outputs = append(outs.Outputs, out)
						UTXO[txID] = outs
					}
				}
			} else {
				outs := UTXO[txID]
				outs.Outputs = append(outs.Outputs, out)
				UTXO[txID] = outs
			}
		}
	}
	return UTXO
}

func AddBlock(block *Block, nodeID string) {
	blockData := block.Serialize()
	database.Set(block.Hash, blockData, nodeID)
	lastBlockHash := database.Get([]byte(lastBlockHashKey), nodeID)
	if lastBlockHash != nil {
		lastBlockData := database.Get(lastBlockHash, nodeID)
		lastBlock := Deserialize(lastBlockData)
		if block.Height > lastBlock.Height {
			database.Set([]byte(lastBlockHashKey), block.Hash, nodeID)
		}
	} else {
		database.Set([]byte(lastBlockHashKey), block.Hash, nodeID)
	}
}

func GetMaxHeight(nodeID string) int64 {
	lastBlockHash := database.Get([]byte(lastBlockHashKey), nodeID)
	if lastBlockHash == nil {
		return 0
	}
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

	for block := GetLastBlock(nodeID); block != nil; block = block.GetPreviousBlock(nodeID) {
		blocks = append(blocks, block.Hash)
	}

	return blocks
}

func MineBlock(transactions []*Transaction, nodeID string) *Block {
	for _, tx := range transactions {
		if VerifyTransaction(tx, nodeID) != true {
			handlers.HandleErrors(errors.New("invalid transaction"))
		}
	}

	lastBlock := GetLastBlock(nodeID)
	lastHeight := lastBlock.Height

	newBlock := CreateBlock(transactions, lastBlock.Hash, lastHeight+1)
	AddBlock(newBlock, nodeID)

	return newBlock
}

func FindUTXO(nodeID string) map[string]TxOutputs {
	UTXO := make(map[string]TxOutputs)
	STXO := make(map[string][]int)

	for block := GetLastBlock(nodeID); block != nil; block = block.GetPreviousBlock(nodeID) {
		STXO = block.GetSTXO(STXO)
		UTXO = block.GetUTXO(STXO, UTXO)
	}

	return UTXO
}

func FindTransaction(txID []byte, nodeID string) (Transaction, error) {
	for block := GetLastBlock(nodeID); block != nil; block = block.GetPreviousBlock(nodeID) {
		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, txID) == 0 {
				return *tx, nil
			}
		}
	}

	return Transaction{}, errors.New("transaction does not exist")
}

func SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey, nodeID string) {
	prevTXs := make(map[string]Transaction)

	for _, in := range tx.Inputs {
		prevTX, err := FindTransaction(in.TxID, nodeID)
		handlers.HandleErrors(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	tx.Sign(privKey, prevTXs)
}

func VerifyTransaction(tx *Transaction, nodeID string) bool {
	if tx.IsCoinbase() {
		return true
	}
	prevTXs := make(map[string]Transaction)

	for _, in := range tx.Inputs {
		prevTX, err := FindTransaction(in.TxID, nodeID)
		handlers.HandleErrors(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	return tx.Verify(prevTXs)
}
