package cli

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/NlaakStudios/Blockchain/utils"

	"github.com/NlaakStudios/Blockchain/config"
	//"os"
)

//CLI is the main client structure
type CLI struct {
	NodePort string
	CoinInfo config.CoinStruct
}

//printUsage diplay commandline usage information to the user.
func (cli *CLI) printUsage() {
	fmt.Println("")
	fmt.Printf("%s Daemon Version %s\n", config.CoinName, config.Version())
	fmt.Printf("%s (%s)\n", config.CoinCompany, config.CoinLandingPage)
	fmt.Println("-----------------------------------------------------------------------------------------")
	fmt.Println("Usage:")
	fmt.Println("	createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("	createwallet - Generates a new key-pair and saves it into the wallet file")
	fmt.Println("	getbalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("	listaddresses - Lists all addresses from the wallet file")
	fmt.Println("	printchain - Print all the blocks of the blockchain")
	fmt.Println("	reindexutxo - Rebuilds the UTXO set")
	fmt.Println("	send -from FROM -to TO -amount AMOUNT -mine - Send AMOUNT of coins from FROM address to TO. Mine on the same node, when -mine is set.")
	fmt.Println("	startnode -miner ADDRESS - Start a node with ID specified in NODE_ID env. var. -miner enables mining")
	fmt.Println("	version - Display node version")
	fmt.Println("")
}

//validateArgs validates the parameters passsed in via commandline
func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) Version() string {
	return config.Version()
}

// Run parses command line arguments and processes commands
func (cli *CLI) Run() {
	//Defaul node port (CONST)
	cli.NodePort = config.NodePort
	utils.CreateDirIfNotExist("data")

	//See if blockchain file exists including data folder. If not crete folder, display notice
	walletsFile := fmt.Sprintf(config.FilePathWallets, config.NodePort)
	if _, err := os.Stat(walletsFile); os.IsNotExist(err) {
		println("Wallets file does not exist, use `blockchain createwallet` to create at least one wallet")
	}

	//See if blockchain file exists including data folder. If not crete folder, display notice
	blockchainFile := fmt.Sprintf(config.FilePathBlockchain, config.NodePort)
	if _, err := os.Stat(blockchainFile); os.IsNotExist(err) {
		println("Blockchain does not exist, use `blockchain createblockchain` to create one")
	}

	//Validate the command line arguments
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	reindexUTXOCmd := flag.NewFlagSet("reindexutxo", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	startNodeCmd := flag.NewFlagSet("startnode", flag.ExitOnError)
	versionCmd := flag.NewFlagSet("version", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")
	sendMine := sendCmd.Bool("mine", false, "Mine immediately on the same node")
	startNodeMiner := startNodeCmd.String("miner", "", "Enable mining mode and send reward to ADDRESS")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "listaddresses":
		err := listAddressesCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "reindexutxo":
		err := reindexUTXOCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "startnode":
		err := startNodeCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "version":
		err := versionCmd.Parse(os.Args[1:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.ShowBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.CreateBlockchain(*createBlockchainAddress)
		cli.PopulateWallets(*createBlockchainAddress)
	}

	if createWalletCmd.Parsed() {
		cli.CreateWallet()
	}

	if listAddressesCmd.Parsed() {
		cli.ListAddresses()
	}

	if printChainCmd.Parsed() {
		cli.PrintChain()
	}

	if reindexUTXOCmd.Parsed() {
		cli.ReIndexUTXO()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}

		cli.Send(*sendFrom, *sendTo, *sendAmount, *sendMine)
	}

	if startNodeCmd.Parsed() {
		cli.StartNode(cli.NodePort, *startNodeMiner)
	}

	if versionCmd.Parsed() {
		fmt.Println(config.Version())
	}

}