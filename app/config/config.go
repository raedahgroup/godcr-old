package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/decred/dcrd/dcrutil"
	flags "github.com/jessevdk/go-flags"
)

const (
	defaultConfigFilename = "godcr.conf"
	defaultHTTPHost       = "127.0.0.1"
	defaultHTTPPort       = "7778"
)

var (
	defaultAppDataDir          = dcrutil.AppDataDir("godcr", false)
	defaultDcrwalletAppDataDir = dcrutil.AppDataDir("dcrwallet", false)
	defaultRPCCertFile         = filepath.Join(defaultDcrwalletAppDataDir, "rpc.cert")

	AppConfigFilePath = filepath.Join(defaultAppDataDir, defaultConfigFilename)
)

// Config holds the top-level options/flags for the application
type Config struct {
	ConfFileOptions
	CommandLineOptions
}

// ConfFileOptions holds the top-level options/flags that are best set in config file rather than in command-line
type ConfFileOptions struct {
	AppDataDir      string `short:"A" long:"appdata" description:"Path to application data directory"`
	UseTestNet      bool   `short:"t" long:"testnet" description:"Connects to testnet wallet instead of mainnet"`
	UseWalletRPC    bool   `short:"w" long:"usewalletrpc" description:"Connect to a running drcwallet daemon over rpc to perform wallet operations"`
	WalletRPCServer string `long:"walletrpcserver" description:"Wallet RPC server address to connect to"`
	WalletRPCCert   string `long:"walletrpccert" description:"Path to dcrwallet certificate file"`
	NoWalletRPCTLS  bool   `long:"nowalletrpctls" description:"Disable TLS when connecting to dcrwallet daemon via RPC"`
	HTTPHost        string `long:"httphost" description:"HTTP server host address or IP"`
	HTTPPort        string `long:"httpport" description:"HTTP server port"`
}

// CommandLineOptions holds the top-level options/flags that are displayed on the command-line menu
type CommandLineOptions struct {
	InterfaceMode string `long:"mode" description:"Interface mode to run" choice:"cli" choice:"http" choice:"nuklear" choice:"terminal"`
	CliOptions
}

type CliOptions struct {
	SyncBlockchain bool `long:"sync" description:"Syncs blockchain when running in cli mode. If used with a command, command is executed after blockchain syncs"`
}

func defaultFileOptions() ConfFileOptions {
	return ConfFileOptions{
		AppDataDir:    defaultAppDataDir,
		WalletRPCCert: defaultRPCCertFile,
		HTTPHost:      defaultHTTPHost,
		HTTPPort:      defaultHTTPPort,
	}
}

// defaultConfig an instance of Config with the defaults set.
func defaultConfig() Config {
	return Config{
		ConfFileOptions: defaultFileOptions(),
	}
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

	for _, arg := range os.Args {
		if !strings.HasPrefix(arg, "-") {
			continue
		}
		var optionName string
		if strings.HasPrefix(arg, "--") {
			optionName = arg[2:]
		} else {
			optionName = arg[1:]
		}
		if isFileOption := isConfigFileOption(optionName); isFileOption {
			return args, config, parser, fmt.Errorf("Unexpected command-line flag/option, "+
				"see godcr -h for supported command-line flags/options"+
				"\nSet other flags/options in %s", AppConfigFilePath)
		}
	}

	// Load additional config from file
	err = parseConfigFile(parser)
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

func parseConfigFile(parser *flags.Parser) error {
	if (parser.Options & flags.IgnoreUnknown) != flags.None {
		options := parser.Options
		parser.Options = flags.None
		defer func() { parser.Options = options }()
	}
	err := flags.NewIniParser(parser).ParseFile(AppConfigFilePath)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok {
			return fmt.Errorf("Error parsing configuration file: %v", err.Error())
		}
		return err
	}
	return nil
}

func isConfigFileOption(name string) (isFileOption bool) {
	if name == "" {
		return
	}
	tConfFileOptions := reflect.TypeOf(ConfFileOptions{})
	for i := 0; i < tConfFileOptions.NumField(); i++ {
		fieldTag := tConfFileOptions.Field(i).Tag
		shortName := fieldTag.Get("short")
		longName := fieldTag.Get("long")
		isFileOption = longName == name || shortName == name
		if isFileOption {
			return
		}
	}
	return
}
