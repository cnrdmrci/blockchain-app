package main

import (
	"blockchain-app/blockchain"
	"fmt"
	"strconv"
)

func main() {
	fmt.Printf("-------------------- Blockchain Started ------------------------\n")

	chain := blockchain.InitBlockChain()
	chain.AddBlock("First Block after Genesis")
	chain.AddBlock("Second Block after Genesis")
	chain.AddBlock("Third Block after Genesis")

	fmt.Printf("--------------------- Mining Completed  ------------------------\n")

	for _, block := range chain.Blocks {
		fmt.Printf("Data in Block  : %s\n", block.Data)
		fmt.Printf("Previous Hash  : %x\n", block.PrevHash)
		fmt.Printf("Hash           : %x\n", block.Hash)
		fmt.Printf("Nonce          : %x\n", block.Nonce)
		fmt.Printf("PoW Validation : %s\n", strconv.FormatBool(blockchain.NewProof(block).Validate()))
		fmt.Printf("---------------------------------------------------------------------------------\n")
	}
}
