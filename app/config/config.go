package config

import (
	"path/filepath"

	"github.com/decred/dcrd/dcrutil"
)

const (
	defaultLogLevel          = "info"
	defaultCurrencyConverter = "none"
)

var (
	defaultAppDataDir = dcrutil.AppDataDir("godcr", false)
	LogFile           = filepath.Join(defaultAppDataDir, "logs/godcr.log")
)

// Config holds app-wide configuration values.
// Struct tags present on each field are used for
// parsing/reading config values from a config file or from db.
type Config struct {
	AppDataDir string `long:"appdata" description:"Path to application data directory."`
	DebugLevel string `long:"debuglevel" description:"Logging level {trace, debug, info, warn, error, critical}."`
	UseTestnet bool   `long:"usetestnet" description:"Use testnet rather than mainnet."`

	Settings `group:"Settings"`
}

type Settings struct {
	SpendUnconfirmed                    bool     `long:"spendunconfirmed" description:"Spend unconfirmed funds"`
	ShowIncomingTransactionNotification bool     `long:"incomingtxnotification" description:"Show incoming transaction notification"`
	ShowNewBlockNotification            bool     `long:"newblocknotification" description:"Show new block notification"`
	CurrencyConverter                   string   `long:"currencyconverter" description:"Currency Converter {none, bitrex}" choice:"none" choice:"bitrex" default:"none"`
	HiddenAccounts                      []uint32 `long:"hiddenaccounts" description:"Accounts with ignored balances"`
	DefaultAccount                      uint32   `long:"defaultaccount" description:"Default account for incoming and outgoing transactions"`
}

func initConfigWithDefaultValues() *Config {
	return &Config{
		AppDataDir: defaultAppDataDir,
		DebugLevel: defaultLogLevel,
		Settings: Settings{
			CurrencyConverter: defaultCurrencyConverter,
		},
	}
}

// LoadConfigFromDb uses dcrlibwallet helper functions to read config values from a bolt db.
func LoadConfigFromDb() (*Config, error) {
	return initConfigWithDefaultValues(), nil
}
