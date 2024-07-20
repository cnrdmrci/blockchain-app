package network

import (
	"blockchain-app/database"
	"blockchain-app/handlers"
	"blockchain-app/network/blockchain_network"
	"google.golang.org/grpc"
	"net"
)

func StartGrpcServer(nodeID, minerAddress string) {
	setNetworkVariablesToCommon(nodeID, minerAddress)
	netListen, netListenErr := net.Listen(protocol, nodeAddress)
	handlers.HandleErrors(netListenErr)

	database.OpenDB(serverNodeID)
	defer database.CloseDB(serverNodeID)
	panicIfBlockchainNotExist()

	go checkMaxHeight()

	grpcServer := grpc.NewServer()
	blockchain_network.RegisterBlockchainServiceServer(grpcServer, &BlockchainServer{})
	grpcServerErr := grpcServer.Serve(netListen)
	handlers.HandleErrors(grpcServerErr)
}
