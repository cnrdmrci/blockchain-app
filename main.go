package main

import (
	"blockchain-app/blockchain"
	"bytes"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
)

type CommandLine struct {
	blockchain *blockchain.BlockChain
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" add -block BLOCK_DATA - add a block to the chain")
	fmt.Println(" print - Prints the blocks in the chain")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) addBlock(data string) {
	cli.blockchain.AddBlock(data)
	fmt.Println("Added Block!")
}

func (cli *CommandLine) printChain() {
	iter := cli.blockchain.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if bytes.Equal(block.PrevHash, make([]byte, 32)) {
			break
		}
	}
}

func (cli *CommandLine) run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

func main() {
	defer os.Exit(0)
	chain := blockchain.InitBlockChain()
	defer chain.Database.Close()

	cli := CommandLine{chain}
	cli.run()
	//fmt.Printf("-------------------- Blockchain Started ------------------------\n")
	//
	//chain := blockchain.InitBlockChain()
	//chain.AddBlock("First Block after Genesis")
	//chain.AddBlock("Second Block after Genesis")
	//chain.AddBlock("Third Block after Genesis")
	//
	//fmt.Printf("--------------------- Mining Completed  ------------------------\n")
	//
	//for _, block := range chain.Blocks {
	//	fmt.Printf("Data in Block  : %s\n", block.Data)
	//	fmt.Printf("Previous Hash  : %x\n", block.PrevHash)
	//	fmt.Printf("Hash           : %x\n", block.Hash)
	//	fmt.Printf("Nonce          : %x\n", block.Nonce)
	//	fmt.Printf("PoW Validation : %s\n", strconv.FormatBool(blockchain.NewProof(block).Validate()))
	//	fmt.Printf("---------------------------------------------------------------------------------\n")
	//}
}
