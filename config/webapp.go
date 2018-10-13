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
		Flash:            "_flash",
		WhiteList:        whitelist.NewBasic(),
		Enable2FA:        false,
		EnableWList:      false,
		TimeBeforeUnused: 120,
		AdminUsername:    "Admin",
		AdminPassword:    "admin123",
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
