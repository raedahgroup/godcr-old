package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	flags "github.com/jessevdk/go-flags"
)

var AppConfigFilePath = filepath.Join(DefaultAppDataDir, "godcr.conf")

// ConfFileOptions holds the top-level options/flags that should be set in config file rather than in command-line
type ConfFileOptions struct {
	AppDataDir      string        `long:"appdata" description:"Path to application data directory."`
	WalletRPCServer string        `long:"walletrpcserver" description:"RPC server address of running dcrwallet daemon. Required to connect to wallet via dcrwallet."`
	WalletRPCCert   string        `long:"walletrpccert" description:"Path to dcrwallet certificate file. Required if walletrpcserver is set."`
	NoWalletRPCTLS  bool          `long:"nowalletrpctls" description:"Disable TLS when connecting to dcrwallet daemon via RPC."`
	HTTPHost        string        `long:"httphost" description:"HTTP server host address or IP when running godcr in http mode."`
	HTTPPort        string        `long:"httpport" description:"HTTP server port when running godcr in http mode."`
	DebugLevel      string        `long:"debuglevel" description:"Logging level {trace, debug, info, warn, error, critical}"`
	Wallets         []*WalletInfo `long:"wallets" description:"Auto detected wallets information"`

	Settings `group:"Settings"`
}

type Settings struct {
	SpendUnconfirmed                    bool   `long:"spendunconfirmed" description:"Spend unconfirmed funds"`
	ShowIncomingTransactionNotification bool   `long:"incomingtxnotification" description:"Show incoming transaction notification"`
	ShowNewBlockNotification            bool   `long:"newblocknotification" description:"Show new block notification"`
	CurrencyConverter                   string `long:"currencyconverter" description:"Currency Converter {none, bitrex}" choice:"none" choice:"bitrex" default:"none"`
}

func defaultFileOptions() ConfFileOptions {
	return ConfFileOptions{
		AppDataDir:    DefaultAppDataDir,
		WalletRPCCert: defaultRPCCertFile,
		HTTPHost:      defaultHTTPHost,
		HTTPPort:      defaultHTTPPort,
		DebugLevel:    defaultLogLevel,
		Settings: Settings{
			CurrencyConverter: defaultCurrencyConverter,
		},
	}
}

// createConfigFile create the configuration file in AppConfigFilePath using the default values
func createConfigFile() (successful bool) {
	configFile, err := os.Create(AppConfigFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "error in creating config file: %s\n", err.Error())
			return
		}
		err = os.Mkdir(DefaultAppDataDir, os.ModePerm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error in creating config file directory: %s\n", err.Error())
			return
		}
		// we were unable to create the file because the dir was not found.
		// we shall attempt to recreate the file now that we have successfully created the dir
		configFile, err = os.Create(AppConfigFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error in creating config file: %s\n", err.Error())
			return
		}
	}
	defer configFile.Close()

	err = saveConfigToFile(defaultFileOptions())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error saving default configuration to file: %s\n", err.Error())
		return
	}

	fmt.Println("Config file created with default values at", AppConfigFilePath)
	return true
}

func ReadConfigFile() (config ConfFileOptions, err error) {
	// load default config values and create parser object with it
	config = defaultFileOptions()
	parser := flags.NewParser(&config, flags.None)

	// read current config file content into config object
	fileParser := flags.NewIniParser(parser)
	err = fileParser.ParseFile(AppConfigFilePath)
	return
}

// UpdateConfigFile reads the config file into a pointer object
// Calls the update function to update the config object as needed
// And saves the updated config object back to file
func UpdateConfigFile(updateConfig func(config *ConfFileOptions)) error {
	config, err := ReadConfigFile()
	if err != nil {
		return fmt.Errorf("error reading config file: %s", err.Error())
	}

	updateConfig(&config)

	// write config object to file
	err = saveConfigToFile(config)
	if err != nil {
		return fmt.Errorf("error saving changes to config file: %s", err.Error())
	}

	return nil
}

func saveConfigToFile(config ConfFileOptions) error {
	parser := flags.NewParser(&config, flags.None)
	fileParser := flags.NewIniParser(parser)
	return fileParser.WriteFile(AppConfigFilePath, flags.IniIncludeComments|flags.IniIncludeDefaults|flags.IniCommentDefaults)
}

// configFileOptions returns a slice of the short names and long names of all config file options
func configFileOptions() (options []string) {
	tConfFileOptions := reflect.TypeOf(ConfFileOptions{})
	for i := 0; i < tConfFileOptions.NumField(); i++ {
		fieldTag := tConfFileOptions.Field(i).Tag

		if shortName, ok := fieldTag.Lookup("short"); ok {
			options = append(options, "-"+shortName)
		}

		if longName, ok := fieldTag.Lookup("long"); ok {
			options = append(options, "--"+longName)
		}
	}
	return
}
