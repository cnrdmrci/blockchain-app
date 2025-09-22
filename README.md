# Multi Node Blockchain Simulation

[Medium Article Series](https://medium.com/@caner.demirci/multi-node-blockchain-simulation-article-series-5beebeae510f)

## The Posts

[1-) Blocks Article](https://medium.com/@caner.demirci/multi-node-blockchain-simulation-part-1-blocks-en-37c530a9c887) : What is a block? A block is a “box” of transactions that also points to the previous block. This creates the chain.

[2-) Transactions and UTXO Article](https://medium.com/@caner.demirci/multi-node-blockchain-simulation-part-2-transactions-and-utxo-en-9b1a6310f58b) : How do coins move? A transaction spends old outputs and creates new ones. We track spendable coins with a simple UTXO set.

[3-) Proof of Work Article](https://medium.com/@caner.demirci/multi-node-blockchain-simulation-part-3-proof-of-work-en-8d90a34e2151) : How is a block mined? The miner keeps guessing a number (nonce) until the block’s hash meets a target. That’s Proof of Work.

[4-) Wallet Article](https://medium.com/@caner.demirci/multi-node-blockchain-simulation-part-4-wallet-en-286d9e21cb54) : How do you own coins? You create a key pair and an address. You sign transactions to prove the coins are yours.

[5-) Database Article](https://medium.com/@caner.demirci/multi-node-blockchain-simulation-part-5-database-en-1feac91549f9) : How do nodes remember state? We store the chain and UTXO data on disk so a node can restart safely and quickly.

[6-) Network Article](https://medium.com/@caner.demirci/multi-node-blockchain-simulation-part-6-network-en-73f39ceff952) : How do nodes talk? With gRPC. One node runs a server; others connect as clients. They share new blocks and transactions and stay in sync.

[7-) Build Simulation Article](https://medium.com/@caner.demirci/multi-node-blockchain-simulation-part-7-build-simulation-en-4be7866c5ffc) : How do you see it all working? Start multiple nodes, mine a block, send coins, and watch the network update in real time.

# Usage

```bash
$ go run main.go
```
```text
Usage:
  createwallet -----------------------------------> Create a new wallet
  listaddresses ----------------------------------> List wallet addressses
  createblockchain -address ADDRESS --------------> Create a blockchain and sends genesis reward to address
  updateblockchain -------------------------------> Update blockchain via other nodes
  printblockchain --------------------------------> Print the blocks in the blockchain
  printlastblock ---------------------------------> Print last block  
  removelastblock --------------------------------> Remove last block from the blockchain  
  getbalance -address ADDRESS --------------------> Get the balance for an address
  send -from FROM -to TO -amount AMOUNT -mine ----> Send amount of coins. Then -mine flag is set, mine off of this node
  reindexutxo ------------------------------------> Rebuilds the UTXO set
  startnode -miner ADDRESS -----------------------> Start a node with TxID specified in NODE_ID env. var. -miner enables mining
```

# Create Wallet

### Node 1
```bash
$ export NODE_ID=3000
$ go run main.go createwallet
```
```text
New address is: 15uQXdzyY8Bki98UQgeXcLtXaDhW4szBjp
```

### Node 2
```bash
$ export NODE_ID=4000
$ go run main.go createwallet
```
```text
New address is: 14sLJH78gy27wSEaWDTHNrhUxLa3yXYVzS
```

### Node 3
```bash
$ export NODE_ID=5000
$ go run main.go createwallet
```
```text
New address is: 1NEnVj4RVVjjZjRqoXBBqoCV5BXneMo97S
```

# Initialize Blockchain

### Node 1
```bash
$ go run main.go createblockchain -address 14sLJH78gy27wSEaWDTHNrhUxLa3yXYVzS
```
```text
Nonce found: 563481
Block Hash : 000000421470170d61f247b3b9e587debc37387b936fac018a474950a4ab8ad2
Genesis block created.
```

# Print Blockchain

### Node 1
```bash
$ go run main.go printblockchain
```
<img width="880" alt="printblockchain" src="https://github.com/user-attachments/assets/1a811e7e-9d27-4072-abd1-cde1bf758423">

# Get Wallet Balance and Check Reward

### Node 1
```bash
$ go run main.go getbalance -address 14sLJH78gy27wSEaWDTHNrhUxLa3yXYVzS
```
```text
Balance of 14sLJH78gy27wSEaWDTHNrhUxLa3yXYVzS: 20 
```

# Start Full Node and Miner Node

### Node 1
```bash
$ go run main.go startnode
```
<img width="470" alt="fullnode" src="https://github.com/user-attachments/assets/a9530473-410f-4a6d-b7ef-56164379bc50">

### Node 3
```bash
$ go run main.go startnode -miner 1NEnVj4RVVjjZjRqoXBBqoCV5BXneMo97S
```
<img width="1016" alt="miner" src="https://github.com/user-attachments/assets/99daf141-8965-47b3-b9a0-da8870df1dfd">

# Send Funds

### Node 2
```bash
$ go run main.go updateblockchain
$ go run main.go send -from 14sLJH78gy27wSEaWDTHNrhUxLa3yXYVzS -to 15uQXdzyY8Bki98UQgeXcLtXaDhW4szBjp -amount 5
```
<img width="1161" alt="send" src="https://github.com/user-attachments/assets/84c8db47-c508-4419-b0a9-cd8c294d1213">

### Node 1 (Full Node)
<img width="1235" alt="Screenshot 2024-07-25 at 22 26 46" src="https://github.com/user-attachments/assets/ffbb8a9f-04c7-4b80-99ab-5941d6d684fa">

### Node 2 (Miner Node)
<img width="1129" alt="Screenshot 2024-07-25 at 22 27 51" src="https://github.com/user-attachments/assets/eef1ff8f-8249-4015-b828-0c346e24c5e4">

# Get Final Balances

### Node 2
```bash
$ go run main.go updateblockchain
$ go run main.go getbalance -address 14sLJH78gy27wSEaWDTHNrhUxLa3yXYVzS
$ go run main.go getbalance -address 1NEnVj4RVVjjZjRqoXBBqoCV5BXneMo97S
$ go run main.go getbalance -address 15uQXdzyY8Bki98UQgeXcLtXaDhW4szBjp
```
```text
Balance of 14sLJH78gy27wSEaWDTHNrhUxLa3yXYVzS: 15
Balance of 1NEnVj4RVVjjZjRqoXBBqoCV5BXneMo97S: 20 //Mining Reward
Balance of 15uQXdzyY8Bki98UQgeXcLtXaDhW4szBjp: 5
```

# Get Final Blockchain
### Node 2
```bash
$ go run main.go printblockchain
```
<img width="1514" alt="finalblockchain" src="https://github.com/user-attachments/assets/cae17a5b-663d-4e39-873d-e8fcdc349ad0">



