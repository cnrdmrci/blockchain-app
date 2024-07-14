package command

import (
	"blockchain-app/wallet"
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
	cli.nodeID = os.Getenv("NODE_ID")
	if cli.nodeID == "" {
		fmt.Println("NODE_ID env is not set!")
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
