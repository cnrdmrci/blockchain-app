package command

import (
	"fmt"
	"os"
	"runtime"
)

type CommandLine struct {
	nodeID string
}

func (cli *CommandLine) Run() {
	cli.setNodeID()
	cli.validateArgs()

	switch os.Args[1] {
	case createWalletFlag:
		cli.createWallet()
	case listAddressesFlag:
		cli.listAddresses()
	case createBlockchainFlag:
		cli.createBlockchain()
	case printBlockchainFlag:
		cli.printBlockchain()
	default:
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" " + createWalletFlag + " -> Create a new wallet")
	fmt.Println(" " + listAddressesFlag + " -> List wallet addressses")
	fmt.Println(" " + createBlockchainFlag + " -address ADDRESS -> Create a blockchain and sends genesis reward to address")
	fmt.Println(" " + printBlockchainFlag + " -> Print the blocks in the blockchain")
}
