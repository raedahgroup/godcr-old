package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/decred/dcrd/dcrutil"
	"github.com/jessevdk/go-flags"
)

const (
	defaultHTTPHost 	 = "127.0.0.1"
	defaultHTTPPort		 = "7778"
	defaultLogLevel      = "info"
	defaultLogDirname    = "logs"
	defaultLogFilename   = "godcr.log"
)

var (
	defaultAppDataDir          = dcrutil.AppDataDir("godcr", false)
	DefaultDcrwalletAppDataDir = dcrutil.AppDataDir("dcrwallet", false)
	defaultRPCCertFile         = filepath.Join(DefaultDcrwalletAppDataDir, "rpc.cert")
	defaultLogDir      		   = filepath.Join(defaultAppDataDir, defaultLogDirname)
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
	DebugLevel      string `short:"d" long:"debuglevel" description:"Logging level {trace, debug, info, warn, error, critical}"`
	LogDir          string `long:"logdir" description:"Directory to log output."`
	LogFilename     string `long:"logfilename" description:"Name of Log File in log directory."`
}

// CommandLineOptions holds the top-level options/flags that are displayed on the command-line menu
type CommandLineOptions struct {
	InterfaceMode string `long:"mode" description:"Interface mode to run" choice:"cli" choice:"http" choice:"nuklear" choice:"terminal" default:"cli"`
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
		DebugLevel:    defaultLogLevel,
		LogDir:        defaultLogDir,
		LogFilename:   defaultLogFilename,
	}
}

// defaultConfig an instance of Config with the defaults set.
func defaultConfig() Config {
	return Config{
		ConfFileOptions: defaultFileOptions(),
	}
}

// LoadConfig parses program configuration from both command-line args and godcr config file, ignoring unknown options and the help flag
// While unknown options seen on the command line are ignored, unknown options in the configuration file return an error.
// Returns the parsed config object, any command-line args that could not be parsed and any error encountered
func LoadConfig() (Config, []string, error) {
	// check if config file does not exist and create it before proceeding
	var configFileExists bool
	if _, err := os.Stat(AppConfigFilePath); os.IsNotExist(err) {
		configFileExists = createConfigFile()
	} else if !os.IsNotExist(err) {
		configFileExists = true
	}

	// load default config values and create parser object with it
	config := defaultConfig()
	parser := flags.NewParser(&config, flags.IgnoreUnknown)

	// parse command-line args and return any error encountered
	unknownArgs, err := parser.Parse()
	if err != nil {
		return config, unknownArgs, err
	}

	// check if any of the unknown command-line args belong in the config file and alert user to set such values in config file only
	if hasConfigFileOption(unknownArgs) {
		return config, unknownArgs, fmt.Errorf("Unexpected command-line flag/option, "+
			"see godcr -h for supported command-line flags/options"+
			"\nSet other flags/options in %s", AppConfigFilePath)
	}

	// if config file doesn't exist, no need to attempt to parse and then re-parse command-line args
	if !configFileExists {
		return config, unknownArgs, nil
	}

	// Load additional config from file
	err = parseConfigFile(parser)
	if err != nil {
		return config, unknownArgs, err
	}

	// Parse command line options again to ensure they take precedence.
	unknownArgs, err = parser.Parse()

	// return parsed config, unknown args encountered and any error that occurred during last parsing
	return config, unknownArgs, err
}

// hasConfigFileOption checks if an unknown arg found in command-line is a config file option that should only be set in the config file
func hasConfigFileOption(unknownArgs []string) bool {
	configFileOptions := configFileOptions()
	isConfigFileOption := func(option string) bool {
		for _, configFileOption := range configFileOptions {
			if configFileOption == option {
				return true
			}
		}
		return false
	}

	for _, arg := range unknownArgs {
		if isConfigFileOption(strings.TrimSpace(arg)) {
			return true
		}
	}

	return false
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
