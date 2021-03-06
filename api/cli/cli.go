package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/NlaakStudios/Blockchain/api/utils"

	"github.com/NlaakStudios/Blockchain/api/core"

	"github.com/NlaakStudios/Blockchain/api/config"
	"github.com/NlaakStudios/gowaf/logger"
	//"os"
)

//Client is the main client structure
type Client struct {
	NodePort string
	//CoinInfo     config.CoinStruct
	Version      string
	GoWAFVersion string
	Config       *config.Config
	Log          logger.Logger
	DataFolder   string
	ConfigName   string
	ConfigFolder string
	isInit       bool
}

// NewMVC creates a new MVC gowaf app. If dir is passed, it should be a directory to look for
// all project folders (config, static, views, models, controllers, etc). The App returned is initialized.
func NewClient(cfgDir, cfgName string) (*Client, error) {
	cli := &Client{
		Version: Version(),
		Log:     logger.NewDefaultLogger(os.Stdout),
	}

	//Prepare Config Folder
	if len(cfgDir) > 0 {
		cli.setDataPath(cfgDir)
	} else {
		cli.setDataPath("./data")
	}

	//Prepare Config Name (without extension)
	if len(cfgName) > 0 {
		cli.ConfigName = cfgName
	} else {
		cli.ConfigName = "blockchain"
	}

	cli.setConfigFolder(fmt.Sprintf("%s/config", cli.DataFolder))

	if cli.ConfigFolder == "" {
		cli.setConfigFolder("config")
	}

	cli.init()

	cli.printHeader()
	return cli, nil
}

// loadConfig loads the configuration file. If cfg is provided, then it is used as the directory
// for searching the configuration files. It defaults to the directory named config in the current
// working directory.
func loadConfig(name string, cfg ...string) (*config.Config, error) {
	cfgDir := "config"
	if len(cfg) > 0 {
		cfgDir = cfg[0]
	}

	// Load configurations.
	cfgFile, err := findConfigFile(cfgDir, name)
	if err != nil {
		return nil, err
	}
	return config.NewConfig(cfgFile)
}

// findConfigFile finds the configuration file name in the directory dir.
func findConfigFile(dir string, name string) (file string, err error) {
	extensions := []string{".json", ".toml", ".yml", ".hcl"}

	for _, ext := range extensions {
		file = filepath.Join(dir, name)
		if info, serr := os.Stat(file); serr == nil && !info.IsDir() {
			//TODO: Coverage -  Need to hit here
			return
		}
		file = file + ext
		if info, serr := os.Stat(file); serr == nil && !info.IsDir() {
			return
		}
	}
	return "", fmt.Errorf("gowaf: can't find configuration file %s in %s", name, dir)
}

// init initializes values to the app components.
func (cli *Client) init() error {
	appConfig, err := loadConfig(cli.ConfigName, cli.ConfigFolder)
	if err != nil {
		return err
	}
	cli.Config = appConfig

	//Defaul node port (CONST)
	cli.NodePort = config.NodePort
	utils.CreateDirIfNotExist(config.FilePathData)

	//Validate the command line arguments
	cli.validateArgs()

	cli.initMasterWallet()

	cli.isInit = true

	return nil
}

// SetFixturePath sets the directory path as a base to all other folders (config, views, etc).
func (cli *Client) setDataPath(dir string) {
	cli.DataFolder = dir
}

// SetConfigFolder sets the directory path to search for the config files.
func (cli *Client) setConfigFolder(dir string) {
	cli.ConfigFolder = dir
}

//printHeader diplay commandline application header to the user.
func (cli *Client) printHeader() {
	fmt.Printf("%s Blockchain Version %s\n", config.CoinName, config.Version())
	fmt.Printf("%s (%s)\n", config.CoinCompany, config.CoinLandingPage)
	fmt.Println("-----------------------------------------------------------------------------------------")
}

//printUsage diplay commandline usage information to the user.
func (cli *Client) printUsage() {
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
func (cli *Client) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		cli.terminate()
	}
}

//validateArgs validates the parameters passsed in via commandline
func (cli *Client) terminate() {
	println("Exited.")
	os.Exit(1)
}

func (cli *Client) initMasterWallet() {
	//See if blockchain file exists including data folder. If not create folder, display notice
	walletsFile := core.GetWalletsFile(cli.NodePort)
	if _, err := os.Stat(walletsFile); os.IsNotExist(err) {
		println("Wallets file does not exist in ", walletsFile, ", use `blockchain createwallet` to create at least one wallet")
		//cli.terminate()
		cli.cmdCreateWallet()
	}

}

func (cli *Client) initBlockchain() {

	//See if blockchain file exists including data folder. If not crete folder, display notice
	blockchainFile := core.GetBlockChainFile(cli.NodePort)
	if _, err := os.Stat(blockchainFile); os.IsNotExist(err) {
		println("Blockchain does not exist at ", blockchainFile, ", use `blockchain createblockchain` to create one")
		cli.terminate()
	}

}
func (cli *Client) cmdCreateWallet() {
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	err := createWalletCmd.Parse(os.Args[2:])
	if err != nil {
		cli.Log.Errors(err)
	}
}

// Run parses command line arguments and processes commands
func (cli *Client) Run() {

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
	//case "createwallet":
	//	err := createWalletCmd.Parse(os.Args[2:])
	//	if err != nil {
	//		log.Panic(err)
	//	}
	case "createwallet":
		cli.cmdCreateWallet()
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
		cli.terminate()
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
