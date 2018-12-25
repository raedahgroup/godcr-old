package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/decred/dcrd/dcrutil"
	flags "github.com/jessevdk/go-flags"
)

const (
	defaultConfigFilename    = "godcr.conf"
	defaultHTTPServerAddress = "127.0.0.1:1234"
)

var (
	defaultAppDataDir          = dcrutil.AppDataDir("godcr", false)
	defaultDcrwalletAppDataDir = dcrutil.AppDataDir("dcrwallet", false)
	defaultRPCCertFile         = filepath.Join(defaultDcrwalletAppDataDir, "rpc.cert")
	defaultConfigFile          = filepath.Join(defaultAppDataDir, defaultConfigFilename)
)

// Config holds the top-level options for the CLI program.
type Config struct {
	ShowVersion       bool   `short:"v" long:"version" description:"Display version information and exit. Any other flag or command is ignored."`
	AppDataDir        string `short:"A" long:"appdata" description:"Application data directory for wallet config, databases and logs"`
	ConfigFile        string `short:"C" long:"configfile" description:"Path to configuration file"`
	TestNet           bool   `short:"t" long:"testnet" description:"Connects to testnet wallet instead of mainnet"`
	UseWalletRPC      bool   `long:"usewalletrpc" description:"Connect to a running drcwallet rpc"`
	WalletRPCServer   string `short:"w" long:"walletrpcserver" description:"Wallet RPC server to connect to"`
	RPCUser           string `short:"u" long:"rpcuser" description:"RPC username"`
	RPCPassword       string `short:"p" long:"rpcpass" default-mask:"-" description:"RPC password"`
	RPCCert           string `short:"c" long:"rpccert" description:"RPC server certificate chain for validation"`
	HTTPServerAddress string `short:"s" long:"serveraddress" description:"Address and port of the HTTP server."`
	HTTPMode          bool   `long:"http" description:"Run in HTTP mode."`
	DesktopMode       bool   `long:"desktop" description:"Run in Desktop mode"`
	NoDaemonTLS       bool   `long:"nodaemontls" description:"Disable TLS"`
	CreateWallet      bool   `long:"createwallet" description:"Creates a new testnet or mainnet wallet if one doesn't already exist"`
	SyncBlockchain    bool   `long:"sync" description:"Syncs blockchain. If used with a command, command is executed after blockchain syncs"`
}

// defaultConfig an instance of Config with the defaults set.
func Default() *Config {
	return &Config{
		AppDataDir:        defaultAppDataDir,
		ConfigFile:        defaultConfigFile,
		RPCCert:           defaultRPCCertFile,
		HTTPServerAddress: defaultHTTPServerAddress,
	}
}

// AppName returns the name of the program binary file that started the process.
func AppName() string {
	appName := filepath.Base(os.Args[0])
	appName = strings.TrimSuffix(appName, filepath.Ext(appName))
	return appName
}

// LoadConfig parses program configuration from both the CLI flags and the config file.
// It returns any non-option arguments encountered, the Config parsed, the parser used, and any
// error, except errors of type flags.ErrHelp.
// If ignoreUnknownOptions is true, then unknown options seen on the command line are ignored.
// However, unknown options in the configuration file must return an error.
func LoadConfig(ignoreUnknownOptions bool) ([]string, Config, *flags.Parser, error) {
	// load defaults first
	config := defaultConfig()

	parser := flags.NewParser(&config, flags.HelpFlag)
	if ignoreUnknownOptions {
		parser.Options = parser.Options | flags.IgnoreUnknown
	}

	args, err := parser.Parse()
	if err != nil && !IsFlagErrorType(err, flags.ErrHelp) {
		return args, config, parser, err
	}

	if config.ShowVersion {
		return args, config, parser, fmt.Errorf(AppVersion())

// ParseConfig parses program configuration from both the CLI command flags and the config file.
// Returns false if an error occurs or version flag was specified
func ParseConfig(config *Config, parser *flags.Parser) bool {
	// stub out the command handler so that the commands are not executed while loading configuration
	parser.CommandHandler = func(command flags.Commander, args []string) error {
		return nil
	}

	_, err := parser.Parse()
	if err != nil && !IsFlagErrorType(err, flags.ErrCommandRequired) {
		handleParseError(err, parser)
		return false
	}

	if config.ShowVersion {
		displayAppVersion()
		return false
	}

	// Load additional config from file
	err = parseConfigFile(parser, config.ConfigFile)
	if err != nil {
		return args, config, parser, err
	}

	// Parse command line options again to ensure they take precedence.
	args, err = parser.Parse()
	if err != nil && !IsFlagErrorType(err, flags.ErrHelp) {
		return args, config, parser, err
	}

	return args, config, parser, nil
}

func parseConfigFile(parser *flags.Parser, file string) error {
	if (parser.Options & flags.IgnoreUnknown) != flags.None {
		options := parser.Options
		parser.Options = flags.None
		defer func() { parser.Options = options }()
	}
	err := flags.NewIniParser(parser).ParseFile(file)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok {
			return fmt.Errorf("Error parsing configuration file: %v", err.Error())
		}
		return err
	}
	return nil
		// error parsing from file
		fmt.Printf("Error parsing configuration file: %s", err.Error())
		return false
	}

	// Parse command line options again to ensure they take precedence.
	_, err = parser.Parse()
	if err != nil && !IsFlagErrorType(err, flags.ErrCommandRequired) {
		handleParseError(err, parser)
		return false
	}

	return true
}

func handleParseError(err error, parser *flags.Parser) {
	if err == nil {
		return
	}
	if (parser.Options & flags.PrintErrors) != flags.None {
		// error printing is already handled by go-flags.
		return
	}
	if IsFlagErrorType(err, flags.ErrHelp) {
		PrintHelp(parser)
	} else {
		fmt.Println(err)
	}
}


func PrintHelp(parser *flags.Parser) {
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

