package command

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"runtime"
)

type CommandLine struct {
	nodeID string
}

func (cli *CommandLine) Run() {
	go closeApp()
	cli.setNodeID()
	cli.validateArgs()

	switch os.Args[1] {
	case createWalletFlag:
		cli.createWallet()
	case listAddressesFlag:
		cli.listAddresses()
	case createBlockchainFlag:
		cli.createBlockchain()
	case updateBlockchainFlag:
		cli.updateBlockchain()
	case printBlockchainFlag:
		cli.printBlockchain()
	case printLastBlockFlag:
		cli.printLastBlock()
	case removeLastBlockFlag:
		cli.removeLastBlock()
	case getBalanceFlag:
		cli.getBalance()
	case sendFlag:
		cli.send()
	case reindexUTXOFlag:
		cli.reindexUTXO()
	case startNodeFlag:
		cli.startNode()
	default:
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) printUsage() {
	c := color.New(color.FgHiCyan)
	space := "  "

	fmt.Println("Usage:")

	c.Print(space + createWalletFlag)
	fmt.Println(" -----------------------------------> Create a new wallet")

	c.Print(space + listAddressesFlag)
	fmt.Println(" ----------------------------------> List wallet addressses")

	c.Print(space + createBlockchainFlag + " -address ADDRESS")
	fmt.Println(" --------------> Create a blockchain and sends genesis reward to address")

	c.Print(space + updateBlockchainFlag)
	fmt.Println(" -------------------------------> Update blockchain via other nodes")

	c.Print(space + printBlockchainFlag)
	fmt.Println(" --------------------------------> Print the blocks in the blockchain")

	c.Print(space + printLastBlockFlag)
	fmt.Println(" ---------------------------------> Print last block")

	c.Print(space + removeLastBlockFlag)
	fmt.Println(" --------------------------------> Remove last block from the blockchain")

	c.Print(space + getBalanceFlag + " -address ADDRESS")
	fmt.Println(" --------------------> Get the balance for an address")

	c.Print(space + sendFlag + " -from FROM -to TO -amount AMOUNT -mine")
	fmt.Println(" ----> Send amount of coins. Then -mine flag is set, mine off of this node")

	c.Print(space + reindexUTXOFlag)
	fmt.Println(" ------------------------------------> Rebuilds the UTXO set")

	c.Print(space + startNodeFlag + " -miner ADDRESS")
	fmt.Println(" -----------------------> Start a node with TxID specified in " + nodeID + " env. var. -miner enables mining")
}
