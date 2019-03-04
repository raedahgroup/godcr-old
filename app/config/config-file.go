package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"text/template"

	flags "github.com/jessevdk/go-flags"
)

var AppConfigFilePath = filepath.Join(DefaultAppDataDir, "godcr.conf")

// ConfFileOptions holds the top-level options/flags that should be set in config file rather than in command-line
type ConfFileOptions struct {
	AppDataDir      string              `long:"appdata" description:"Path to application data directory"`
	UseWalletRPC    bool                `long:"usewalletrpc" description:"Connect to a running drcwallet daemon over rpc to perform wallet operations"`
	WalletRPCServer string              `long:"walletrpcserver" description:"Wallet RPC server address to connect to. Required if usewalletrpc=true"`
	WalletRPCCert   string              `long:"walletrpccert" description:"Path to dcrwallet certificate file. Required if usewalletrpc=true"`
	NoWalletRPCTLS  bool                `long:"nowalletrpctls" description:"Disable TLS when connecting to dcrwallet daemon via RPC"`
	HTTPHost        string              `long:"httphost" description:"HTTP server host address or IP"`
	HTTPPort        string              `long:"httpport" description:"HTTP server port"`
	DebugLevel      string              `long:"debuglevel" description:"Logging level {trace, debug, info, warn, error, critical}"`
	LogDir          string              `long:"logdir" description:"Directory to log output."`
	LogFilename     string              `long:"logfilename" description:"Name of Log File in log directory."`
	Wallets         []*WalletInfo `long:"wallets" description:"Auto detected wallets information"`
}

func defaultFileOptions() ConfFileOptions {
	return ConfFileOptions{
		AppDataDir:    DefaultAppDataDir,
		WalletRPCCert: defaultRPCCertFile,
		HTTPHost:      defaultHTTPHost,
		HTTPPort:      defaultHTTPPort,
		DebugLevel:    defaultLogLevel,
		LogDir:        defaultLogDir,
		LogFilename:   defaultLogFilename,
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

	tmpl := template.New("config")

	tmpl, err = tmpl.Parse(configTextTemplate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error preparing default config file content: %s", err.Error())
		return
	}

	err = tmpl.Execute(configFile, defaultFileOptions())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error saving default configuration to file: %s\n", err.Error())
		return
	}

	fmt.Println("Config file created with default values at", AppConfigFilePath)
	return true
}

// todo defer?
func parseConfigFile(parser *flags.Parser) error {
	if (parser.Options & flags.IgnoreUnknown) != flags.None {
		options := parser.Options
		parser.Options = flags.None
		defer func() { parser.Options = options }()
	}
	err := flags.NewIniParser(parser).ParseFile(AppConfigFilePath)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok {
			return fmt.Errorf("Error parsing configuration file: %s", err.Error())
		}
		return err
	}
	return nil
}

// UpdateConfigFile reads the config file into a pointer object
// Calls the update function to update the config object as needed
// And saves the updated config object back to file
func UpdateConfigFile(updateConfig func(config *Config)) error {
	// load default config values and create parser object with it
	config := defaultConfig()
	parser := flags.NewParser(&config, flags.None)

	// read current config file content into config object
	fileParser := flags.NewIniParser(parser)
	err := fileParser.ParseFile(AppConfigFilePath)
	if err != nil {
		return fmt.Errorf("error reading config file: %s", err.Error())
	}

	updateConfig(&config)

	// write config object to file
	err = fileParser.WriteFile(AppConfigFilePath, flags.IniIncludeComments|flags.IniIncludeDefaults|flags.IniCommentDefaults)
	if err != nil {
		return fmt.Errorf("error saving changes to config file: %s", err.Error())
	}

	return nil
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
