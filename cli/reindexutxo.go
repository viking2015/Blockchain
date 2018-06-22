package cli

import (
	"fmt"

	"github.com/NlaakStudios/Blockchain/core"
)

func (cli *CLI) reindexUTXO() {
	bc := core.NewBlockchain(cli.NodeID)
	UTXOSet := core.UTXOSet{bc}
	UTXOSet.Reindex()

	count := UTXOSet.CountTransactions()
	fmt.Printf("Done! There are %d transactions in the UTXO set.\n", count)
}
