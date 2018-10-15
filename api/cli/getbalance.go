package cli

import (
	"fmt"
	"log"

	"github.com/NlaakStudios/Blockchain/api/config"
	"github.com/NlaakStudios/Blockchain/api/core"
	"github.com/NlaakStudios/Blockchain/api/utils"
)

//ShowBalance shows the balance of the given wallet in the console
func (cli *Client) ShowBalance(address string) {
	if !core.ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := core.NewBlockchain(cli.NodePort)
	UTXOSet := core.UTXOSet{bc}
	defer bc.DB.Close()

	balance := 0
	pubKeyHash := utils.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := UTXOSet.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d %s\n", address, balance, config.CoinSymbol)
}

// GetBalance given a valid address returns the current balance in coins
func (cli *Client) GetBalance(address string) int {

	//fmt.Println("getBalance(%s, %s)", address, nodeID)
	if !core.ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}

	bc := core.NewBlockchain(cli.NodePort)
	UTXOSet := core.UTXOSet{bc}
	defer bc.DB.Close()

	balance := 0
	pubKeyHash := utils.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := UTXOSet.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	//fmt.Printf("Balance of '%s': %d\n", address, balance)
	return balance
}
