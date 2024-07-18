package wallet

import (
	"blockchain-app/handlers"
	"encoding/json"
	"fmt"
	"os"
)

const walletFile = "./database/wallets_%s.data"

type Wallets struct {
	Wallets map[string]*Wallet
}

func GetWallets(nodeId string) *Wallets {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	wallets.LoadFromFile(nodeId)

	return &wallets
}

func (ws *Wallets) AddNewWallet() string {
	wallet := createNewWallet()
	address := wallet.GetAddress()

	ws.Wallets[address] = wallet

	return address
}

func (ws *Wallets) GetAllWalletAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

func (ws *Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

func (ws *Wallets) LoadFromFile(nodeId string) {
	walletFile := fmt.Sprintf(walletFile, nodeId)
	if _, err := os.Stat(walletFile); err == nil {
		fileContent, err := os.ReadFile(walletFile)
		handlers.HandleErrors(err)
		_ = json.Unmarshal(fileContent, ws)
	}
}

func (ws *Wallets) SaveFile(nodeId string) {
	jsonData, err := json.MarshalIndent(ws, "", "  ")
	handlers.HandleErrors(err)

	walletFile := fmt.Sprintf(walletFile, nodeId)

	err = os.WriteFile(walletFile, jsonData, 0666)
	handlers.HandleErrors(err)
}
