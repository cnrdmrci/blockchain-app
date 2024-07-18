package blockchain

import (
	"blockchain-app/database"
	"blockchain-app/handlers"
	"bytes"
	"encoding/hex"
)

var (
	utxoPrefix = []byte("utxo-")
)

func FindSpendableOutputs(pubKeyHash []byte, amount int, nodeID string) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	accumulated := 0

	for iter := database.GetIteratorByPrefix(utxoPrefix, nodeID); iter.ValidForPrefix(utxoPrefix); iter.Next() {
		item := iter.Item()
		key := item.Key()
		value, err := item.ValueCopy(nil)
		handlers.HandleErrors(err)
		key = bytes.TrimPrefix(key, utxoPrefix)
		txID := hex.EncodeToString(key)
		outs := DeserializeOutputs(value)

		for outIdx, out := range outs.Outputs {
			if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)
			}
		}
		if accumulated >= amount {
			break
		}
	}

	return accumulated, unspentOuts
}

func FindUnspentTransactions(pubKeyHash []byte, nodeID string) []TxOutput {
	var UTXOs []TxOutput

	for iter := database.GetIteratorByPrefix(utxoPrefix, nodeID); iter.ValidForPrefix(utxoPrefix); iter.Next() {
		item := iter.Item()
		value, err := item.ValueCopy(nil)
		handlers.HandleErrors(err)
		outs := DeserializeOutputs(value)
		for _, out := range outs.Outputs {
			if out.IsLockedWithKey(pubKeyHash) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

func CountTransactions(nodeID string) int {
	return database.CountByPrefix(utxoPrefix, nodeID)
}

func Reindex(nodeID string) {
	database.DeleteByPrefix(utxoPrefix, nodeID)
	UTXO := FindUTXO(nodeID)
	for txId, outs := range UTXO {
		key, err := hex.DecodeString(txId)
		handlers.HandleErrors(err)
		key = append(utxoPrefix, key...)
		database.Set(key, outs.Serialize(), nodeID)
	}
}

func UpdateIndex(block *Block, nodeID string) {
	for _, tx := range block.Transactions {
		if !tx.IsCoinbase() {
			for _, in := range tx.Inputs {
				updatedOuts := TxOutputs{}
				inID := append(utxoPrefix, in.TxID...)
				v := database.Get(inID, nodeID)
				outs := DeserializeOutputs(v)

				for outIdx, out := range outs.Outputs {
					if outIdx != in.Out {
						updatedOuts.Outputs = append(updatedOuts.Outputs, out)
					}
				}

				if len(updatedOuts.Outputs) == 0 {
					database.Delete(inID, nodeID)
				} else {
					database.Set(inID, updatedOuts.Serialize(), nodeID)
				}
			}
		}
		newOutputs := TxOutputs{}
		for _, out := range tx.Outputs {
			newOutputs.Outputs = append(newOutputs.Outputs, out)
		}

		txID := append(utxoPrefix, tx.ID...)
		database.Set(txID, newOutputs.Serialize(), nodeID)
	}
}
