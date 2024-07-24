package network

import (
	"blockchain-app/blockchain"
	"bytes"
	"encoding/hex"
)

var (
	nodeAddress                     string
	mineAddress                     string
	serverNodeID                    string
	protocol                        = "tcp"
	KnownNodeAddresses              = []string{"localhost:3000"}
	MaxKnownAddressCount            = 5
	ControlBlockForTransactionCount = 3
	MemPoolMineCount                = 1
	MemPool                         = make(map[string]blockchain.Transaction)
)

func getTransactionIdForMemPool(transaction *blockchain.Transaction) string {
	txId := hex.EncodeToString(transaction.ID)
	return serverNodeID + "_" + txId
}

func getTransactionFromMemPool(transaction *blockchain.Transaction) (blockchain.Transaction, bool) {
	tx, exists := MemPool[getTransactionIdForMemPool(transaction)]
	return tx, exists
}

func isTransactionExistInLastAFewBlocks(transaction *blockchain.Transaction) bool {
	block := blockchain.GetLastBlock(serverNodeID)
	for blockCount := 0; blockCount < ControlBlockForTransactionCount; blockCount++ {
		for _, blockTransaction := range block.Transactions {
			if bytes.Compare(blockTransaction.ID, transaction.ID) == 0 {
				return true
			}
		}
		if block.IsGenesis() {
			break
		}
		block = block.GetPreviousBlock(serverNodeID)
	}

	return false
}
