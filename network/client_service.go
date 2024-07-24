package network

import (
	"blockchain-app/blockchain"
	"blockchain-app/handlers"
	"blockchain-app/network/blockchain_network"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"math/rand"
	"time"
)

func setNetworkVariablesToCommon(nodeID, minerAddress string) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	mineAddress = minerAddress
	serverNodeID = nodeID
}

func addUniqueKnownNodeAddress(address string) {
	exists := false

	if len(KnownNodeAddresses) >= MaxKnownAddressCount {
		return
	}

	if address == "" {
		return
	}

	for _, addr := range KnownNodeAddresses {
		if addr == address {
			exists = true
			break
		}
	}

	if !exists {
		KnownNodeAddresses = append(KnownNodeAddresses, address)
		color.Yellow("New node found. Address: %s", address)
	}
}

func removeKnownNodeAddress(addressToRemove string) {
	for i, address := range KnownNodeAddresses {
		if address == addressToRemove {
			KnownNodeAddresses = append(KnownNodeAddresses[:i], KnownNodeAddresses[i+1:]...)
			color.Red("Server is down, so removed from known node addresses: %s", addressToRemove)
			return
		}
	}
}

func getRandomKnownAddress(exceptAddress string) string {
	if len(KnownNodeAddresses) == 0 {
		return ""
	}

	filteredAddresses := make([]string, 0, len(KnownNodeAddresses))
	for _, address := range KnownNodeAddresses {
		if address != exceptAddress {
			filteredAddresses = append(filteredAddresses, address)
		}
	}

	if len(filteredAddresses) == 0 {
		return ""
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	randomIndex := r.Intn(len(filteredAddresses))
	return filteredAddresses[randomIndex]
}

func panicIfBlockchainNotExist() {
	if !blockchain.IsBlockchainExist(serverNodeID) {
		for _, knownNode := range KnownNodeAddresses {
			if nodeAddress != knownNode {
				if err := createBlockchainViaOtherNode(knownNode); err == nil {
					break
				}
			}
		}
	}
	if !blockchain.IsBlockchainExist(serverNodeID) {
		handlers.HandleErrors(errors.New("blockchain not exist"))
	}
}

func createBlockchainViaOtherNode(nodeAddress string) error {
	conn, clientErr := grpc.NewClient(nodeAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if clientErr != nil {
		return clientErr
	}
	defer conn.Close()
	client := blockchain_network.NewBlockchainServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	stream, requestErr := client.StreamGetAllBlockchain(ctx, &blockchain_network.GetAllBlockchainRequest{})
	if requestErr != nil {
		return requestErr
	}
	var lastBlock *blockchain.Block
	for {
		grpcBlock, streamErr := stream.Recv()
		if streamErr == io.EOF {
			break
		}
		if streamErr != nil {
			return streamErr
		}
		block := mapGrpcBlockToBlock(grpcBlock)
		blockchain.AddBlock(&block, serverNodeID)
		lastBlock = &block
	}
	if lastBlock != nil {
		color.Green("Blockchain created. Last Block Hash: %x\n", lastBlock.Hash)
		blockchain.Reindex(serverNodeID)
	}
	return nil
}

func checkMaxHeight() {
	printServerMaxHeight()
	for {
		time.Sleep(3 * time.Second)
		ourServerMaxHeight := blockchain.GetMaxHeight(serverNodeID)
		for _, knownNode := range KnownNodeAddresses {
			if nodeAddress != knownNode {
				otherServerHeight := getMaxHeightFromOtherServer(knownNode, nodeAddress)
				if otherServerHeight > ourServerMaxHeight {
					createBlocksUntilHeightFromOtherNode(knownNode, ourServerMaxHeight)
				}
			}
		}
	}
}

func UpdateBlockchainViaOtherNodes(nodeID string) {
	setNetworkVariablesToCommon(nodeID, "")
	ourServerMaxHeight := blockchain.GetMaxHeight(serverNodeID)
	for _, knownNode := range KnownNodeAddresses {
		if nodeAddress != knownNode {
			otherServerHeight := getMaxHeightFromOtherServer(knownNode, "")
			if otherServerHeight > ourServerMaxHeight {
				createBlocksUntilHeightFromOtherNode(knownNode, ourServerMaxHeight)
				color.Green("Blockchain updated.")
				return
			}
		}
	}
	color.Green("Blockchain up to date.")
}

func takeTransactions() {
	for {
		time.Sleep(5 * time.Second)
		for _, knownNode := range KnownNodeAddresses {
			if nodeAddress != knownNode {
				takeTransactionsFromOtherNode(knownNode)
			}
		}
	}
}

func takeTransactionsFromOtherNode(serverAddress string) {
	conn, clientErr := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if clientErr != nil {
		removeKnownNodeAddress(serverAddress)
		return
	}
	defer conn.Close()
	client := blockchain_network.NewBlockchainServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Minute)
	defer cancel()
	stream, requestErr := client.StreamGetAllTransactions(ctx, &blockchain_network.GetAllTransactionsRequest{})
	if requestErr != nil {
		removeKnownNodeAddress(serverAddress)
		return
	}
	for {
		grpcTransaction, streamErr := stream.Recv()
		if streamErr == io.EOF {
			break
		}
		handlers.HandleErrors(streamErr)
		transaction := mapGrpcTransactionToTransaction(grpcTransaction)
		if !blockchain.VerifyTransaction(transaction, serverNodeID) {
			color.Red("Invalid transaction. TxID: %x\n", hex.EncodeToString(transaction.ID))
			continue
		}
		if !isTransactionExistInLastAFewBlocks(transaction) {
			if _, isTransactionExistOnMemPool := getTransactionFromMemPool(transaction); !isTransactionExistOnMemPool {
				MemPool[getTransactionIdForMemPool(transaction)] = *transaction
				color.Green("New Transaction added. MemPool Count: %x , TxID: %x\n", len(MemPool), transaction.ID)
			}
		}
	}
}

func printServerMaxHeight() {
	ourServerMaxHeight := blockchain.GetMaxHeight(serverNodeID)
	color.HiCyan("Server Address: %s , Max height: %d", nodeAddress, ourServerMaxHeight)
}

func getMaxHeightFromOtherServer(serverAddress string, ourServerNodeAddress string) int64 {
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		removeKnownNodeAddress(serverAddress)
		return 0
	}
	defer conn.Close()
	client := blockchain_network.NewBlockchainServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	res, err := client.GetMaxHeight(ctx, &blockchain_network.GetMaxHeightRequest{RequesterNodeAddress: ourServerNodeAddress})
	if err != nil {
		removeKnownNodeAddress(serverAddress)
		return 0
	}
	addUniqueKnownNodeAddress(res.RandomActiveNodeAddress)
	return res.Height
}

func createBlocksUntilHeightFromOtherNode(serverAddress string, ourServerMaxHeight int64) {
	conn, clientErr := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if clientErr != nil {
		removeKnownNodeAddress(serverAddress)
		return
	}
	defer conn.Close()
	client := blockchain_network.NewBlockchainServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	stream, requestErr := client.StreamGetBlocksUntilHeight(ctx, &blockchain_network.GetBlocksUntilHeightRequest{Height: ourServerMaxHeight})
	if requestErr != nil {
		removeKnownNodeAddress(serverAddress)
		return
	}
	for {
		grpcBlock, streamErr := stream.Recv()
		if streamErr == io.EOF {
			break
		}
		handlers.HandleErrors(streamErr)
		color.HiCyan("New Block Found.")
		block := mapGrpcBlockToBlock(grpcBlock)
		blockchain.AddBlock(&block, serverNodeID)
		blockchain.UpdateIndex(&block, serverNodeID)
		color.Green("New Block added. Max Height: %x , Block Hash: %x\n", block.Height, block.Hash)
		removeTransactionsFromMemPoolIfBlockContains(block)
	}
	printServerMaxHeight()
}

func removeTransactionsFromMemPoolIfBlockContains(block blockchain.Block) {
	for _, transaction := range block.Transactions {
		if _, exist := getTransactionFromMemPool(transaction); exist {
			delete(MemPool, getTransactionIdForMemPool(transaction))
			color.Red("Transaction which in the block were removed. TxID: %x\n", transaction.ID)
		}
	}
}

func mineTxToBlockViaMemPool() {
	if mineAddress != "" {
		for {
			time.Sleep(5 * time.Second)
			if len(MemPool) >= MemPoolMineCount {
				color.HiCyan("%d transactions found. Mining started.", MemPoolMineCount)
				transactionsForCreateBlock := []*blockchain.Transaction{}
				takenMemPoolCount := 0
				for _, transaction := range MemPool {
					if takenMemPoolCount == MemPoolMineCount {
						break
					}
					transactionsForCreateBlock = append(transactionsForCreateBlock, &transaction)
					takenMemPoolCount++
				}
				transactionsForCreateBlock = append(transactionsForCreateBlock, blockchain.CreateMineRewardTx(mineAddress))
				block := blockchain.MineBlock(transactionsForCreateBlock, serverNodeID)
				blockchain.UpdateIndex(block, serverNodeID)

				color.Green("New Block added. Mining Completed.")
				for _, transaction := range transactionsForCreateBlock {
					delete(MemPool, getTransactionIdForMemPool(transaction))
				}
				printServerMaxHeight()
			}
		}
	}
}

func SendTransaction(transaction *blockchain.Transaction, nodeID string) {
	randomAddress := getRandomKnownAddress("localhost:" + nodeID)
	conn, err := grpc.NewClient(randomAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		removeKnownNodeAddress(randomAddress)
		return
	}
	defer conn.Close()
	client := blockchain_network.NewBlockchainServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	res, err := client.CreateTransaction(ctx, mapTransactionToGrpcTransaction(transaction))
	if err != nil {
		removeKnownNodeAddress(randomAddress)
		return
	}
	if res.Success {
		color.Green("Transaction send to the server successfully. TxID: %x\n", transaction.ID)
	} else {
		color.Red(res.Message)
	}
}
