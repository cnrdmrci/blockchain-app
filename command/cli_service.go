package command

import (
	"blockchain-app/blockchain"
	"blockchain-app/database"
	"blockchain-app/encoders"
	"blockchain-app/handlers"
	"blockchain-app/network"
	"blockchain-app/wallet"
	"errors"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/vrecan/death/v3"
	"os"
	"runtime"
	"syscall"
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
	rewardAddress := getFlagValue(createBlockchainFlag, addressFlag, "The address to send genesis block reward to")

	if !wallet.ValidateAddress(rewardAddress) {
		handlers.HandleErrors(errors.New("address is not valid"))
	}

	database.OpenDB(cli.nodeID)
	defer database.CloseDB(cli.nodeID)

	blockchain.InitBlockChain(rewardAddress, cli.nodeID)
	blockchain.Reindex(cli.nodeID)

	color.HiCyan("Genesis block created.")
}

func (cli *CommandLine) updateBlockchain() {
	database.OpenDB(cli.nodeID)
	defer database.CloseDB(cli.nodeID)

	network.UpdateBlockchainViaOtherNodes(cli.nodeID)
}

func (cli *CommandLine) printBlockchain() {
	database.OpenDB(cli.nodeID)
	defer database.CloseDB(cli.nodeID)

	for block := blockchain.GetLastBlock(cli.nodeID); block != nil; block = block.GetPreviousBlock(cli.nodeID) {
		block.PrintBlockDetails()
	}
}

func (cli *CommandLine) printLastBlock() {
	database.OpenDB(cli.nodeID)
	defer database.CloseDB(cli.nodeID)

	block := blockchain.GetLastBlock(cli.nodeID)
	block.PrintBlockDetails()
}

func (cli *CommandLine) removeLastBlock() {
	database.OpenDB(cli.nodeID)
	defer database.CloseDB(cli.nodeID)
	blockchain.RemoveLastBlock(cli.nodeID)
	blockchain.Reindex(cli.nodeID)
}

func (cli *CommandLine) getBalance() {
	balanceAddress := getFlagValue(getBalanceFlag, addressFlag, "The address to get balance for")

	if !wallet.ValidateAddress(balanceAddress) {
		handlers.HandleErrors(errors.New("address is not valid"))
	}

	pubKeyHash := encoders.Base58Decode(balanceAddress)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]

	database.OpenDB(cli.nodeID)
	defer database.CloseDB(cli.nodeID)
	UTXOs := blockchain.FindUnspentTransactions(pubKeyHash, cli.nodeID)

	balance := 0
	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", balanceAddress, balance)
}

func (cli *CommandLine) send() {
	sendCmd := flag.NewFlagSet(sendFlag, flag.ExitOnError)

	sendFrom := sendCmd.String(fromFlag, "", "Source wallet address")
	sendTo := sendCmd.String(toFlag, "", "Destination wallet address")
	sendAmount := sendCmd.Int(amountFlag, 0, "Amount to send")
	sendMine := sendCmd.Bool(mineFlag, false, "Mine immediately on the same node")

	_ = sendCmd.Parse(os.Args[2:])
	if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
		sendCmd.Usage()
		runtime.Goexit()
	}

	if !wallet.ValidateAddress(*sendTo) {
		handlers.HandleErrors(errors.New("address is not valid"))
	}
	if !wallet.ValidateAddress(*sendFrom) {
		handlers.HandleErrors(errors.New("address is not valid"))
	}

	wallets := wallet.GetWallets(cli.nodeID)
	walletSendFrom := wallets.GetWallet(*sendFrom)

	database.OpenDB(cli.nodeID)
	defer database.CloseDB(cli.nodeID)
	tx := blockchain.NewTransaction(&walletSendFrom, *sendTo, *sendAmount, cli.nodeID)
	if *sendMine {
		mineRewardTx := blockchain.CreateMineRewardTx(*sendFrom)
		txs := []*blockchain.Transaction{mineRewardTx, tx}
		block := blockchain.MineBlock(txs, cli.nodeID)
		blockchain.UpdateIndex(block, cli.nodeID)
	} else {
		network.SendTransaction(tx, cli.nodeID)
	}

	fmt.Println("Success!")
}

func (cli *CommandLine) reindexUTXO() {
	reindexUTXOCmd := flag.NewFlagSet(reindexUTXOFlag, flag.ExitOnError)
	err := reindexUTXOCmd.Parse(os.Args[2:])
	handlers.HandleErrors(err)

	database.OpenDB(cli.nodeID)
	defer database.CloseDB(cli.nodeID)
	blockchain.Reindex(cli.nodeID)
	count := blockchain.CountTransactions(cli.nodeID)
	fmt.Printf("Done! There are %d transactions in the UTXO set.\n", count)
}

func (cli *CommandLine) startNode() {
	startNodeCmd := flag.NewFlagSet(startNodeFlag, flag.ExitOnError)
	startNodeMiner := startNodeCmd.String(minerFlag, "", "Enable mining mode and send reward to ADDRESS")
	err := startNodeCmd.Parse(os.Args[2:])
	handlers.HandleErrors(err)

	fmt.Printf("Starting Node %s\n", cli.nodeID)

	if len(*startNodeMiner) > 0 {
		if wallet.ValidateAddress(*startNodeMiner) {
			color.HiCyan("Mining is on. Address to receive rewards: %s", *startNodeMiner)
		} else {
			handlers.HandleErrors(errors.New("miner address not valid"))
		}
	}
	network.StartGrpcServer(cli.nodeID, *startNodeMiner)
}

func getFlagValue(commandName string, flagName string, usageMessage string) string {
	blockchainCmd := flag.NewFlagSet(commandName, flag.ExitOnError)
	blockchainCmdValue := blockchainCmd.String(flagName, "", usageMessage)
	_ = blockchainCmd.Parse(os.Args[2:])
	if *blockchainCmdValue == "" {
		blockchainCmd.Usage()
		runtime.Goexit()
	}
	return *blockchainCmdValue
}

func closeApp() {
	d := death.NewDeath(syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	d.WaitForDeathWithFunc(func() {
		color.Red("\napp closing!\n")
		defer os.Exit(1)
		defer runtime.Goexit()
	})
}
