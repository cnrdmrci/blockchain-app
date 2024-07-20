package network

import (
	"blockchain-app/blockchain"
	"blockchain-app/network/blockchain_network"
	"context"
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
