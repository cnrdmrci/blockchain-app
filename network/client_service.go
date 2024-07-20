package network

import (
	"blockchain-app/blockchain"
	"blockchain-app/handlers"
	"blockchain-app/network/blockchain_network"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"time"
)

func setNetworkVariablesToCommon(nodeID, minerAddress string) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	mineAddress = minerAddress
	serverNodeID = nodeID
}

func panicIfBlockchainNotExist() {
	if !blockchain.IsBlockchainExist(serverNodeID) {
		for _, knownNode := range KnownNodes {
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
	}
	return nil
}
