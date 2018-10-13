package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/fatih/camelcase"
	"github.com/hashicorp/hcl"
	"gopkg.in/yaml.v2"
)

/*
Required Directory Structure for WebApp
webapp-binary.exe
|
|- serve				`all static files served fromt nis folder and all subfolders`
     |- assets			`All WebApp specific files (css, js, images, audio and video)`
     |    |- js
     |    |- css
     |    |- img
     |    |- vid
     |    |- aud
     |- vendor			`All vendor assets - each in thier own folder`
     |- widgets			`HTML Snippets wrapped in a div with data-widget="WidgetName"`
     |- templates		`route template files, ie Landing.html, Dashboard.tpl, Login.html, ect.`

*/
var errBlockchainCfgUnsupported = errors.New("gowaf: Blockchain config file format not supported")

// Config stores configurations values
type BlockchainConfig struct {
	//////////////////////////////////////
	// Blockchainn Specific Configuration
	//////////////////////////////////////
	//CoinName is your Coin or Tokens ICO name. ie Bitcoin Copper ICO
	CoinName string `json:"coin_name" yaml:"coin_name" toml:"coin_name" hcl:"coin_name"`
	//CoinSymbol is the 3-4 uppercase letters used to identify your coin in the exchange/listings
	CoinSymbol string `json:"coin_symbol" yaml:"coin_symbol" toml:"coin_symbol" hcl:"coin_symbol"`
	//CoinPurpose is the reason for your new coin or token. What will it privide. Summarize here.
	CoinPurpose string `json:"coin_purpose" yaml:"coin_purpose" toml:"coin_purpose" hcl:"coin_purpose"`
	//CoinAMLCompliant shows if your ICO is AML compliant. Yes|No
	CoinAMLCompliant bool `json:"coin_aml_compliant" yaml:"coin_aml_compliant" toml:"coin_aml_compliant" hcl:"coin_aml_compliant"`
	//CoinKYCCompliant shows if your ICO is KYC compliant. Yes|No
	CoinKYCCompliant bool `json:"coin_kyc_compliant" yaml:"coin_kyc_compliant" toml:"coin_kyc_compliant" hcl:"coin_kyc_compliant"`
	//CoinCompany is the registered name for your company
	//CoinCompany // Use CompanyName
	//CoinCountry is the country in which you compan is registered
	//CoinCountry = "Cayman Islands"
	//CoinCEO is the CEO or Owner of the registered company
	//CoinCEO = "Andrew Donelson"
	//CoinContact = is the valid direct email to the {CoinCEO} and is also the Account Email
	//CoinContact = "gwf@nlaak.com"
	//CoinSubsidy is the Number of coins given to all new addresses
	CoinSubsidy uint `json:"coin_subsidy" yaml:"coin_subsidy" toml:"coin_subsidy" hcl:"coin_subsidy"`
	//CoinICOSupply is the total number of coins sold in ICO
	CoinICOSupply uint64 `json:"coin_ico_supply" yaml:"coin_ico_supply" toml:"coin_ico_supply" hcl:"coin_ico_supply"`
	//CoinDevSupply is the total numbner of coins reserved for developers (initial bonus)
	CoinDevSupply uint64 `json:"coin_dev_supply" yaml:"coin_dev_supply" toml:"coin_dev_supply" hcl:"coin_dev_supply"`
	//CoinOAMSupply is the total numnber of coins reserved for Operating and Marketing Costs (min 6 months)
	CoinOAMSupply uint64 `json:"coin_oam_supply" yaml:"coin_oam_supply" toml:"coin_oam_supply" hcl:"coin_oam_supply"`
	//CoinPLTSupply is the total number of coins reserved for working capital/misc expenses not forseen
	CoinPLTSupply uint64 `json:"coin_plt_supply" yaml:"coin_plt_supply" toml:"coin_plt_supply" hcl:"coin_plt_supply"`
	CoinDecimals  uint   `json:"coin_decimals" yaml:"coin_decimals" toml:"coin_decimals" hcl:"coin_decimals"`

	/* INTERNAL USE ONLY */
	//12M, 10M ICO, 400k Staff, 1.6M Operations & Marketing
	CoinTotalSupply uint64 `json:"coin_total_supply" yaml:"coin_total_supply" toml:"coin_total_supply" hcl:"coin_total_supply"`

	//CoinLandingPage is the URL to your landing page (We handle this)
	//format: https://www/gwf.io/ico/{CoinSymbol}/
	CoinLandingPage string `json:"coin_landing" yaml:"coin_landing" toml:"coin_landing" hcl:"coin_landing"`

	//CoinWhitePapere is the URL to your whitepaper (We handle this)
	//format: https://www/gwf.io/ico/{CoinSymbol}/
	CoinWhitePaper string `json:"coin_whitepaper" yaml:"coin_whitepaper" toml:"coin_whitepaper" hcl:"coin_whitepaper"`

	//CoinDashboardPage is the URL to your users dashboard page (We handle this)
	//format: https://www/gwf.io/ico/{CoinSymbol}/dashboard
	CoinDashboardPage string `json:"coin_dashboard" yaml:"coin_dashboard" toml:"coin_dashboard" hcl:"coin_dashboard"`

	//CoinAPI is the URL to your API endpoints (We handle this)
	//format: https://www/gwf.io/ico/api/{CoinSymbol}/{EndPoint}
	CoinAPI string `json:"coin_api" yaml:"coin_api" toml:"coin_api" hcl:"coin_api"`

	//FilePathData is the complete path to the master data directory
	//Format: data/blockchain-{CoinPort}.db
	FilePathData string `json:"filepath_data" yaml:"filepath_data" toml:"filepath_data" hcl:"filepath_data"`

	//FilePathBlockchain is the complete path to the ICO's blockchain database
	//Format: data/ico/{CoinSymbol}/blockchain-{CoinPort}.db
	FilePathBlockchain string `json:"filepath_blockchain" yaml:"filepath_blockchain" toml:"filepath_blockchain" hcl:"filepath_blockchain"`

	//FilePathWallets is the complete path to the ICO's wallets file
	//Format: data/ico/{CoinSymbol}/wallets-{CoinPort}.dat
	FilePathWallets string `json:"filepath_wallets" yaml:"filepath_wallets" toml:"filepath_wallets" hcl:"filepath_wallets"`
}

// DefaultBlockchainConfig returns the default configuration settings for the blockchain.
func DefaultBlockchainConfig() *BlockchainConfig {
	return &BlockchainConfig{
		CoinName:           "GWF ICO Coin",
		CoinSymbol:         "GWFI",
		CoinPurpose:        "This blockchain and coin supply is being used for the GWF ICO and will be swapped out 1:1 for Spartacus Token when the actual chains development is complete.",
		CoinAMLCompliant:   false,
		CoinKYCCompliant:   false,
		CoinSubsidy:        0,
		CoinICOSupply:      uint64(10000000), //10M ICO
		CoinDevSupply:      uint64(500000),   //500k Dev Team
		CoinOAMSupply:      uint64(1500000),  //1.5M Operations & Marketing
		CoinPLTSupply:      uint64(3000000),  //3M Platform, Misc, Petty Cash, Capital, etc.
		CoinTotalSupply:    uint64(15000000),
		CoinDecimals:       10,
		CoinLandingPage:    "https://www/gwf.io/%s",
		CoinWhitePaper:     "https://www/gwf.io/%s/whitepaper",
		CoinDashboardPage:  "https://www/gwf.io/%s/dashboard",
		CoinAPI:            "https://www/gwf.io/api/%s/%s",
		FilePathData:       "../data",
		FilePathBlockchain: "blockchain/blockchain-%s.db",
		FilePathWallets:    "blockchain/wallets-%s.dat",
	}
}

// NewConfig reads configuration from path. The format is deduced from the file extension
//	* .json    - is decoded as json
//	* .yml     - is decoded as yaml
//	* .toml    - is decoded as toml
//  * .hcl	   - is decoded as hcl
func NewBlockchainConfig(path string) (*BlockchainConfig, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &BlockchainConfig{}
	switch filepath.Ext(path) {
	case ".json":
		jerr := json.Unmarshal(data, cfg)
		if jerr != nil {
			return nil, jerr
		}
	case ".toml":
		_, terr := toml.Decode(string(data), cfg)
		if terr != nil {
			return nil, terr
		}
	case ".yml":
		yerr := yaml.Unmarshal(data, cfg)
		if yerr != nil {
			return nil, yerr
		}
	case ".hcl":
		obj, herr := hcl.Parse(string(data))
		if herr != nil {
			return nil, herr
		}
		if herr = hcl.DecodeObject(&cfg, obj); herr != nil {
			return nil, herr
		}
	default:
		return nil, errCfgUnsupported
	}
	err = cfg.SyncEnv()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// SyncEnv overrides c field's values that are set in the environment.
//
// The environment variable names are derived from config fields by underscoring, and uppercasing
// the name. E.g. AppName will have a corresponding environment variable APP_NAME
//
// NOTE only int, string and bool fields are supported and the corresponding values are set.
// when the field value is not supported it is ignored.
func (c *BlockchainConfig) SyncEnv() error {
	cfg := reflect.ValueOf(c).Elem()
	cTyp := cfg.Type()

	for k := range make([]struct{}, cTyp.NumField()) {
		field := cTyp.Field(k)

		cm := getEnvName(field.Name)
		env := os.Getenv(cm)
		if env == "" {
			continue
		}
		switch field.Type.Kind() {
		case reflect.String:
			cfg.FieldByName(field.Name).SetString(env)
		case reflect.Int:
			v, err := strconv.Atoi(env)
			if err != nil {
				return fmt.Errorf("gowaf: loading Blockchain config field %s %v", field.Name, err)
			}
			cfg.FieldByName(field.Name).Set(reflect.ValueOf(v))
		case reflect.Bool:
			b, err := strconv.ParseBool(env)
			if err != nil {
				return fmt.Errorf("gowaf: loading Blockchain config field %s %v", field.Name, err)
			}
			cfg.FieldByName(field.Name).SetBool(b)
		}

	}
	return nil
}

// getEnvName returns all upper case and underscore separated string, from field.
// field is a camel case string.
//
// example
//	AppName will change to APP_NAME
func getEnvName(field string) string {
	camSplit := camelcase.Split(field)
	var rst string
	for k, v := range camSplit {
		if k == 0 {
			rst = strings.ToUpper(v)
			continue
		}
		rst = rst + "_" + strings.ToUpper(v)
	}
	return rst
}
