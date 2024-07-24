package network

import (
	"blockchain-app/blockchain"
	"blockchain-app/network/blockchain_network"
	"blockchain-app/wallet"
	"context"
	"fmt"
	"github.com/fatih/color"
)

type BlockchainServer struct {
	blockchain_network.UnimplementedBlockchainServiceServer
}

func (s *BlockchainServer) GetMaxHeight(ctx context.Context, req *blockchain_network.GetMaxHeightRequest) (*blockchain_network.GetMaxHeightResponse, error) {
	addUniqueKnownNodeAddress(req.RequesterNodeAddress)
	heightResponse := &blockchain_network.GetMaxHeightResponse{Height: blockchain.GetMaxHeight(serverNodeID), RandomActiveNodeAddress: getRandomKnownAddress(req.RequesterNodeAddress)}
	return heightResponse, nil
}

func (s *BlockchainServer) StreamGetAllBlockchain(req *blockchain_network.GetAllBlockchainRequest, stream blockchain_network.BlockchainService_StreamGetAllBlockchainServer) error {
	for block := blockchain.GetLastBlock(serverNodeID); block != nil; block = block.GetPreviousBlock(serverNodeID) {
		if err := stream.Send(mapBlockToGrpcBlock(*block)); err != nil {
			return err
		}
	}
	return nil
}

func (s *BlockchainServer) StreamGetBlocksUntilHeight(req *blockchain_network.GetBlocksUntilHeightRequest, stream blockchain_network.BlockchainService_StreamGetBlocksUntilHeightServer) error {
	for block := blockchain.GetLastBlock(serverNodeID); block != nil; block = block.GetPreviousBlock(serverNodeID) {
		if block.Height == req.Height {
			return nil
		}
		if err := stream.Send(mapBlockToGrpcBlock(*block)); err != nil {
			return err
		}
	}
	return nil
}

func (s *BlockchainServer) CreateTransaction(ctx context.Context, grpcTransaction *blockchain_network.Transaction) (*blockchain_network.CreateTransactionResponse, error) {
	transaction := mapGrpcTransactionToTransaction(grpcTransaction)

	if len(transaction.Inputs) > 1 {
		return &blockchain_network.CreateTransactionResponse{
			Success: false,
			Message: "Transaction can contains only one input",
		}, nil
	}

	if isTransactionInputAlreadyExistInMemPool(transaction) {
		return &blockchain_network.CreateTransactionResponse{
			Success: false,
			Message: "Utxo already exist in MemPool",
		}, nil
	}

	if !isTransactionSenderWalletsHaveEnoughFunds(transaction) {
		return &blockchain_network.CreateTransactionResponse{
			Success: false,
			Message: "Not enough funds",
		}, nil
	}

	if isTransactionExistInLastAFewBlocks(transaction) {
		return &blockchain_network.CreateTransactionResponse{
			Success: false,
			Message: "Transaction already exists",
		}, nil
	}

	if !blockchain.VerifyTransaction(transaction, serverNodeID) {
		return &blockchain_network.CreateTransactionResponse{
			Success: false,
			Message: "invalid transaction",
		}, nil
	}

	MemPool[getTransactionIdForMemPool(transaction)] = *transaction
	color.Green("New transaction added to MemPool. MemPool Count: %d, TxID: %x", len(MemPool), transaction.ID)
	return &blockchain_network.CreateTransactionResponse{Success: true}, nil
}

func (s *BlockchainServer) StreamGetAllTransactions(req *blockchain_network.GetAllTransactionsRequest, stream blockchain_network.BlockchainService_StreamGetAllTransactionsServer) error {
	for _, transaction := range MemPool {
		if err := stream.Send(mapTransactionToGrpcTransaction(&transaction)); err != nil {
			return err
		}
	}
	return nil
}

func isTransactionSenderWalletsHaveEnoughFunds(transaction *blockchain.Transaction) bool {
	neededFunds := 0
	for _, output := range transaction.Outputs {
		neededFunds += output.Value
	}

	totalFunds := 0
	for _, input := range transaction.Inputs {
		pubHashKey := wallet.GetPublicKeyHash(input.PubKey)
		utxos := blockchain.FindUnspentTransactions(pubHashKey, serverNodeID)
		balance := 0
		for _, out := range utxos {
			balance += out.Value
		}
		totalFunds += balance
	}

	return totalFunds >= neededFunds
}

func isTransactionInputAlreadyExistInMemPool(transaction *blockchain.Transaction) bool {
	for _, input := range transaction.Inputs {
		utxoKey := fmt.Sprintf("%x:%d", input.TxID, input.Out)
		for _, memPoolTx := range MemPool {
			for _, memInput := range memPoolTx.Inputs {
				memPoolUtxoKey := fmt.Sprintf("%x:%d", memInput.TxID, memInput.Out)
				if memPoolUtxoKey == utxoKey {
					return true
				}
			}
		}
	}
	return false
}
