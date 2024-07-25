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



