package network

import (
	"blockchain-app/blockchain"
	"blockchain-app/network/blockchain_network"
)

func mapBlockToGrpcBlock(block blockchain.Block) *blockchain_network.Block {
	grpcTransactions := make([]*blockchain_network.Transaction, len(block.Transactions))
	for i, tx := range block.Transactions {
		grpcTransactions[i] = mapTransactionToGrpcTransaction(tx)
	}
	return &blockchain_network.Block{
		Timestamp:    block.Timestamp,
		Hash:         block.Hash,
		Transactions: grpcTransactions,
		PrevHash:     block.PrevHash,
		Nonce:        block.Nonce,
		Height:       block.Height,
	}
}

func mapTransactionToGrpcTransaction(tx *blockchain.Transaction) *blockchain_network.Transaction {
	pbInputs := make([]*blockchain_network.TxInput, len(tx.Inputs))
	for i, in := range tx.Inputs {
		pbInputs[i] = mapTxInputToPB(&in)
	}
	pbOutputs := make([]*blockchain_network.TxOutput, len(tx.Outputs))
	for i, out := range tx.Outputs {
		pbOutputs[i] = mapTxOutputToPB(&out)
	}
	return &blockchain_network.Transaction{
		Id:      tx.ID,
		Inputs:  pbInputs,
		Outputs: pbOutputs,
	}
}

func mapTxInputToPB(in *blockchain.TxInput) *blockchain_network.TxInput {
	return &blockchain_network.TxInput{
		TxId:      in.TxID,
		Out:       int32(in.Out),
		Signature: in.Signature,
		PubKey:    in.PubKey,
	}
}

func mapTxOutputToPB(out *blockchain.TxOutput) *blockchain_network.TxOutput {
	return &blockchain_network.TxOutput{
		Value:      int32(out.Value),
		PubKeyHash: out.PubKeyHash,
	}
}

func mapGrpcBlockToBlock(pbBlock *blockchain_network.Block) blockchain.Block {
	transactions := make([]*blockchain.Transaction, len(pbBlock.Transactions))
	for i, pbTx := range pbBlock.Transactions {
		transactions[i] = mapGrpcTransactionToTransaction(pbTx)
	}
	return blockchain.Block{
		Timestamp:    pbBlock.Timestamp,
		Hash:         pbBlock.Hash,
		Transactions: transactions,
		PrevHash:     pbBlock.PrevHash,
		Nonce:        pbBlock.Nonce,
		Height:       pbBlock.Height,
	}
}

func mapGrpcTransactionToTransaction(pbTx *blockchain_network.Transaction) *blockchain.Transaction {
	inputs := make([]blockchain.TxInput, len(pbTx.Inputs))
	for i, pbIn := range pbTx.Inputs {
		inputs[i] = mapPBToTxInput(pbIn)
	}
	outputs := make([]blockchain.TxOutput, len(pbTx.Outputs))
	for i, pbOut := range pbTx.Outputs {
		outputs[i] = mapPBToTxOutput(pbOut)
	}
	return &blockchain.Transaction{
		ID:      pbTx.Id,
		Inputs:  inputs,
		Outputs: outputs,
	}
}

func mapPBToTxInput(pbIn *blockchain_network.TxInput) blockchain.TxInput {
	return blockchain.TxInput{
		TxID:      pbIn.TxId,
		Out:       int(pbIn.Out),
		Signature: pbIn.Signature,
		PubKey:    pbIn.PubKey,
	}
}

func mapPBToTxOutput(pbOut *blockchain_network.TxOutput) blockchain.TxOutput {
	return blockchain.TxOutput{
		Value:      int(pbOut.Value),
		PubKeyHash: pbOut.PubKeyHash,
	}
}
