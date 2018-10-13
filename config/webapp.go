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

	"github.com/casbin/gorm-adapter"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/securecookie"
	"github.com/hashicorp/hcl"
	"github.com/kisom/whitelist"
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
var errCfgUnsupported = errors.New("gowaf: config file format not supported")

// Config stores configurations values
type Config struct {
	AppName          string `json:"app_name" yaml:"app_name" toml:"app_name" hcl:"app_name"`
	Domain           string `json:"domain" yaml:"domain" toml:"domain" hcl:"domain"`
	CompanyName      string `json:"company_name" yaml:"company_name" toml:"company_name" hcl:"company_name"`
	BaseURL          string `json:"base_url" yaml:"base_url" toml:"base_url" hcl:"base_url"`
	Port             int    `json:"port" yaml:"port" toml:"port" hcl:"port"`
	Verbose          bool   `json:"verbose" yaml:"verbose" toml:"verbose" hcl:"verbose"`
	StaticDir        string `json:"static_dir" yaml:"static_dir" toml:"static_dir" hcl:"static_dir"`
	Database         string `json:"database" yaml:"database" toml:"database" hcl:"database"`
	DatabaseConn     string `json:"database_conn" yaml:"database_conn" toml:"database_conn" hcl:"database_conn"`
	DatabasePassword string `json:"database_password" yaml:"database_password" toml:"database_password" hcl:"database_password"`
	Automigrate      bool   `json:"automigrate" yaml:"automigrate" toml:"automigrate" hcl:"automigrate"`
	DropTables       bool   `json:"droptables" yaml:"droptables" toml:"droptables" hcl:"droptables"`
	LoadTestData     bool   `json:"load_test_data" yaml:"load_test_data" toml:"load_test_data" hcl:"load_test_data"`
	NoModel          bool   `json:"no_model" yaml:"no_model" toml:"no_model" hcl:"no_model"`
	GoogleID         string `json:"googleid" yaml:"googleid" toml:"googleid" hcl:"googleid"`
	Notifications    bool   `json:"notifications" yaml:"notifications" toml:"notifications" hcl:"notifications"`
	WebAppOnboard    bool   `json:"webapp_onboard" yaml:"webapp_onboard" toml:"webapp_onboard" hcl:"webapp_onboard"`
	CompanyOnboard   bool   `json:"company_onboard" yaml:"company_onboard" toml:"company_onboard" hcl:"company_onboard"`
	UserOnboard      bool   `json:"user_onboard" yaml:"user_onboard" toml:"user_onboard" hcl:"user_onboard"`
	Mail             bool   `json:"mail" yaml:"mail" toml:"mail" hcl:"mail"`
	MailServer       string `json:"mail_server" yaml:"mail_server" toml:"mail_server" hcl:"mail_server"`
	MailFrom         string `json:"mail_user" yaml:"mail_user" toml:"mail_user" hcl:"mail_user"`
	MailUsername     string `json:"mail_username" yaml:"mail_username" toml:"mail_username" hcl:"mail_username"`
	MailPort         int    `json:"mail_port" yaml:"mail_port" toml:"_mailport" hcl:"_mail_port"`
	MailPassword     string `json:"mail_password" yaml:"mail_password" toml:"mail_password" hcl:"mail_password"`
	Profile          bool   `json:"profile" yaml:"profile" toml:"profile" hcl:"profile"`
	ThemeColor       string `json:"themecolor" yaml:"themecolor" toml:"themecolor" hcl:"themecolor"`
	Flash            string `json:"flash" yaml:"flash" toml:"flash" hcl:"flash"`
	FlashTime        uint   `json:"flash_time" yaml:"flash_time" toml:"flash_time" hcl:"flash_time"`
	FlashStack       uint   `json:"flash_stack" yaml:"flash_stack" toml:"flash_stack" hcl:"flash_stack"`
	FlashContextKey  string `json:"flash_context_key" yaml:"flash_context_key" toml:"flash_context_key" hcl:"flash_context_key"`
	SessionName      string `json:"session_name" yaml:"session_name" toml:"session_name" hcl:"session_name"`
	SessionPath      string `json:"session_path" yaml:"session_path" toml:"session_path" hcl:"session_path"`
	SessionDomain    string `json:"session_domain" yaml:"session_domain" toml:"session_domain" hcl:"session_domain"`
	SessionMaxAge    int    `json:"session_max_age" yaml:"session_max_age" toml:"session_max_age" hcl:"session_max_age"`
	SessionSecure    bool   `json:"session_secure" yaml:"session_secure" toml:"session_secure" hcl:"session_secure"`
	SessionHTTPOnly  bool   `json:"session_httponly" yaml:"session_httponly" toml:"session_httponly" hcl:"session_httponly"`
	Enable2FA        bool   `json:"enable2fa" yaml:"enable2fa" toml:"enable2fa" hcl:"enable2fa"`
	EnableWList      bool   `json:"enableWList" yaml:"enableWList" toml:"enableWList" hcl:"enableWList"`
	AdminUsername    string `json:"admin_username" yaml:"admin_username" toml:"admin_username" hcl:"admin_username"`
	AdminPassword    string `json:"admin_password" yaml:"admin_password" toml:"admin_password" hcl:"admin_password"`
	// The name of the session store to use. Options are: file , cookie ,ql
	SessionStore string `json:"session_store" yaml:"session_store" toml:"session_store" hcl:"session_store"`
	// KeyPair for secure cookie its a comma separates strings of keys.
	SessionKeyPair []string `json:"session_key_pair" yaml:"session_key_pair" toml:"session_key_pair" hcl:"session_key_pair"`
	WhiteList      *whitelist.Basic
	//Adapter for casbin
	Adapter *gormadapter.Adapter
	//Time in hours before model is active
	TimeBeforeUnused int `json:"time_before_unused" yaml:"time_before_unused" toml:"time_before_unused" hcl:"time_before_unused"`
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

// DefaultConfig returns the default configuration settings.
func DefaultConfig() *Config {
	a := securecookie.GenerateRandomKey(32)
	b := securecookie.GenerateRandomKey(32)
	return &Config{
		AppName:          "GoWAF Blockchain WebApp",
		Domain:           "nlaak.com",
		CompanyName:      "Nlaak Studios",
		BaseURL:          "http://localhost:8090",
		Port:             8090,
		ThemeColor:       "blue",
		GoogleID:         "",
		Verbose:          false,
		Database:         "",
		DatabaseConn:     "",
		DatabasePassword: "",
		Automigrate:      true,
		DropTables:       true,
		NoModel:          true,
		Notifications:    false,
		WebAppOnboard:    false,
		CompanyOnboard:   false,
		UserOnboard:      false,
		Mail:             false,
		MailServer:       "",
		MailFrom:         "noreply",
		MailUsername:     "",
		MailPort:         587,
		MailPassword:     "",
		Profile:          false,
		LoadTestData:     false,
		FlashTime:        3500,
		FlashStack:       6,
		SessionName:      "_gowaf",
		SessionPath:      "/",
		SessionMaxAge:    2592000,
		Adapter:          nil,
		SessionKeyPair: []string{
			string(a), string(b),
		},
		Flash:              "_flash",
		WhiteList:          whitelist.NewBasic(),
		Enable2FA:          false,
		EnableWList:        false,
		TimeBeforeUnused:   120,
		AdminUsername:      "Admin",
		AdminPassword:      "admin123",
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
func NewConfig(path string) (*Config, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
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

	// ensure the key pairs are set
	if cfg.SessionKeyPair == nil {
		a := securecookie.GenerateRandomKey(32)
		b := securecookie.GenerateRandomKey(32)
		cfg.SessionKeyPair = []string{
			string(a), string(b),
		}
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
func (c *Config) SyncEnv() error {
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
				return fmt.Errorf("gowaf: loading config field %s %v", field.Name, err)
			}
			cfg.FieldByName(field.Name).Set(reflect.ValueOf(v))
		case reflect.Bool:
			b, err := strconv.ParseBool(env)
			if err != nil {
				return fmt.Errorf("gowaf: loading config field %s %v", field.Name, err)
			}
			cfg.FieldByName(field.Name).SetBool(b)
		}

	}
	return nil
}
