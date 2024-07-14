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
	default:
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" " + createWalletFlag + " - Create a new wallet")
	fmt.Println(" " + listAddressesFlag + " - List wallet addressses")
}
