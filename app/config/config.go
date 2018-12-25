package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/decred/dcrd/dcrutil"
	"github.com/jessevdk/go-flags"
	"github.com/raedahgroup/dcrcli/app"
	"github.com/raedahgroup/dcrcli/cli/commands"
)

const (
	defaultConfigFilename    = "dcrcli.conf"
	defaultHTTPServerAddress = "127.0.0.1:7778"
)

var (
	defaultAppDataDir          = dcrutil.AppDataDir("dcrcli", false)
	defaultDcrwalletAppDataDir = dcrutil.AppDataDir("dcrwallet", false)
	defaultRPCCertFile         = filepath.Join(defaultDcrwalletAppDataDir, "rpc.cert")
	defaultConfigFile          = filepath.Join(defaultAppDataDir, defaultConfigFilename)
)

// Config holds the top-level options for the application and cli-only command options/flags/args
type Config struct {
	commands.CliCommands
	AppDataDir        string `short:"A" long:"appdata" description:"Path to application data directory"`
	ConfigFile        string `short:"C" long:"configfile" description:"Path to configuration file"`
	ShowVersion       bool   `short:"v" long:"version" description:"Display version information and exit"`
	CreateWallet      bool   `long:"createwallet" description:"Create a new testnet or mainnet wallet if one doesn't already exist"`
	SyncBlockchain    bool   `long:"sync" description:"Syncs blockchain. If used with a cli command, command is executed after blockchain syncs"`
	UseTestNet        bool   `short:"t" long:"testnet" description:"Connects to testnet wallet instead of mainnet"`
	UseWalletRPC      bool   `short:"w" long:"usewalletrpc" description:"Connect to a running drcwallet daemon over rpc to perform wallet operations"`
	WalletRPCServer   string `long:"walletrpcserver" description:"Wallet RPC server address to connect to"`
	WalletRPCCert     string `long:"walletrpccert" description:"Path to dcrwallet certificate file"`
	NoWalletRPCTLS    bool   `long:"nowalletrpctls" description:"Disable TLS when connecting to dcrwallet daemon via RPC"`
	HTTPMode          bool   `long:"http" description:"Run in HTTP mode"`
	HTTPServerAddress string `long:"httpserveraddress" description:"Address and port for the HTTP server"`
}

// defaultConfig an instance of Config with the defaults set.
func defaultConfig() Config {
	return Config{
		AppDataDir:        defaultAppDataDir,
		ConfigFile:        defaultConfigFile,
		WalletRPCCert:     defaultRPCCertFile,
		HTTPServerAddress: defaultHTTPServerAddress,
	}
}

// LoadConfig loads program configuration by parsing options/flags from the command-line and the config file
func LoadConfig() *Config {
	// create parser with default config
	config := defaultConfig()
	parser := flags.NewParser(&config, flags.HelpFlag)

	// stub out the command handler so that the commands are not executed while loading configuration
	parser.CommandHandler = func(command flags.Commander, args []string) error {
		return nil
	}

	_, err := parser.Parse()
	if err != nil && !IsFlagErrorType(err, flags.ErrCommandRequired) {
		handleParseError(err, parser)
		return nil
	}

	if config.ShowVersion {
		fmt.Printf("%s version: %s\n", app.Name(), app.Version())
		return nil
	}

	// Load additional config from file
	err = flags.NewIniParser(parser).ParseFile(config.ConfigFile)
	if err != nil {
		fmt.Printf("Error parsing configuration file: %s", err.Error())
		return nil
	}

	// Parse command line options again to ensure they take precedence.
	_, err = parser.Parse()
	if err != nil && !IsFlagErrorType(err, flags.ErrCommandRequired) {
		handleParseError(err, parser)
		return nil
	}

	return &config
}

func handleParseError(err error, parser *flags.Parser) {
	if IsFlagErrorType(err, flags.ErrHelp) {
		printHelp(parser)
	} else {
		fmt.Println(err)
	}
}

func printHelp(parser *flags.Parser) {
	if parser.Active == nil {
		// Print help for the root command (general help with all the options and commands).
		parser.WriteHelp(os.Stderr)
	} else {
		// Print a concise command-specific help.
		printCommandHelp(parser.Name, parser.Active)
	}
}

func printCommandHelp(appName string, command *flags.Command) {
	helpParser := flags.NewParser(nil, flags.HelpFlag)
	helpParser.Name = appName
	helpParser.Active = command
	helpParser.WriteHelp(os.Stderr)
	fmt.Printf("To view application options, use '%s -h'\n", appName)
}
