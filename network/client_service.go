package network

import (
	"blockchain-app/blockchain"
	"blockchain-app/handlers"
	"blockchain-app/network/blockchain_network"
	"context"
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
		color.Yellow("New node added. Address: %s", address)
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
		fmt.Printf("Blockchain created. Last Block Hash: %x\n", lastBlock.Hash)
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
				otherServerHeight := getMaxHeightFromOtherServer(knownNode)
				if otherServerHeight > ourServerMaxHeight {
					createBlocksUntilHeightFromOtherNode(knownNode, ourServerMaxHeight)
				}
			}
		}
	}
}

func printServerMaxHeight() {
	ourServerMaxHeight := blockchain.GetMaxHeight(serverNodeID)
	color.HiCyan("Server Address: %s , Max height: %d", nodeAddress, ourServerMaxHeight)
}

func getMaxHeightFromOtherServer(serverAddress string) int64 {
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		removeKnownNodeAddress(serverAddress)
		return 0
	}
	defer conn.Close()
	client := blockchain_network.NewBlockchainServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	res, err := client.GetMaxHeight(ctx, &blockchain_network.GetMaxHeightRequest{RequesterNodeAddress: nodeAddress})
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
		block := mapGrpcBlockToBlock(grpcBlock)
		blockchain.AddBlock(&block, serverNodeID)
		color.Green("New Block added. Max Height: %x , Block Hash: %x\n", block.Height, block.Hash)
	}
	printServerMaxHeight()
}
