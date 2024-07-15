package command

import (
	"blockchain-app/blockchain"
	"blockchain-app/database"
	"blockchain-app/handlers"
	"blockchain-app/wallet"
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime"
)

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) setNodeID() {
	cli.nodeID = os.Getenv(nodeID)
	if cli.nodeID == "" {
		fmt.Println(nodeID + " env is not set!")
		runtime.Goexit()
	}
}

func (cli *CommandLine) createWallet() {
	address := wallet.CreateAndSaveNewWallet(cli.nodeID)
	fmt.Printf("New address is: %s\n", address)
	runtime.Goexit()
}

func (cli *CommandLine) listAddresses() {
	wallets := wallet.GetWallets(cli.nodeID)
	addresses := wallets.GetAllWalletAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
	runtime.Goexit()
}

func (cli *CommandLine) createBlockchain() {
	rewardAddress := getFlagValue(createBlockchainFlag, "address", "The address to send genesis block reward to")

	if !wallet.ValidateAddress(rewardAddress) {
		handlers.HandleErrors(errors.New("address is not valid"))
	}

	database.OpenDB(cli.nodeID)
	defer database.CloseDB(cli.nodeID)

	_ = blockchain.InitBlockChain(rewardAddress, cli.nodeID)
	fmt.Println("Genesis created")

	//UTXOSet := blockchain.UTXOSet{blockChain}
	//UTXOSet.Reindex()

	fmt.Println("Finished!")
}

func (cli *CommandLine) printBlockchain() {
	database.OpenDB(cli.nodeID)
	defer database.CloseDB(cli.nodeID)

	block := blockchain.GetLastBlock(cli.nodeID)
	block.PrintBlockDetails()
	if block.IsGenesis() {
		return
	}
	for {
		block = block.GetPreviousBlock(cli.nodeID)
		block.PrintBlockDetails()
		if block.IsGenesis() {
			break
		}
	}
}

func getFlagValue(commandName string, flagName string, usageMessage string) string {
	blockchainCmd := flag.NewFlagSet(commandName, flag.ExitOnError)
	blockchainCmdValue := blockchainCmd.String(flagName, "", usageMessage)
	blockchainCmd.Parse(os.Args[2:])
	if *blockchainCmdValue == "" {
		blockchainCmd.Usage()
		runtime.Goexit()
	}
	return *blockchainCmdValue
}
