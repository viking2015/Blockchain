package cli

import (
	"fmt"
	"strconv"

	"github.com/NlaakStudios/Blockchain/api/core"
)

//PrintChain prints the blockchain out to stdout
func (cli *CLI) PrintChain() {
	bc := core.NewBlockchain(cli.NodePort)
	defer bc.DB.Close()

	bci := bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("============ Block %x ============\n", block.Hash)
		fmt.Printf("Height: %d\n", block.Height)
		fmt.Printf("Prev. block: %x\n", block.PrevBlockHash)
		pow := core.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n\n", strconv.FormatBool(pow.Validate()))
		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		fmt.Printf("\n\n")

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}