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
	defaultHTTPHost          = "127.0.0.1"
	defaultHTTPPort          = "7778"
	defaultLogLevel          = "info"
	defaultCurrencyConverter = "none"
)

var (
	defaultAppDataDir          = dcrutil.AppDataDir("godcr", false)
	DefaultDcrwalletAppDataDir = dcrutil.AppDataDir("dcrwallet", false)
	defaultRPCCertFile         = filepath.Join(DefaultDcrwalletAppDataDir, "rpc.cert")
	LogFile                    = filepath.Join(defaultAppDataDir, "logs/godcr.log")
)

// Config holds the top-level options/flags for the application
type Config struct {
	ConfFileOptions
	CommandLineOptions
}

// CommandLineOptions holds the top-level options/flags that are displayed on the command-line menu
type CommandLineOptions struct {
	CliOptions
}

type CliOptions struct {
	SyncBlockchain bool `long:"sync" description:"Syncs blockchain when running in cli mode. If used with a command, command is executed after blockchain syncs"`
}

// defaultConfig an instance of Config with the defaults set.
func defaultConfig() Config {
	return Config{
		ConfFileOptions:    defaultFileOptions(),
		CommandLineOptions: defaultCommandLineOptions(),
	}
}

func defaultCommandLineOptions() CommandLineOptions {
	return CommandLineOptions{}
}

// LoadConfig parses program configuration from both command-line args and godcr config file, ignoring unknown options and the help flag
// While unknown options seen on the command line are ignored, unknown options in the configuration file return an error.
// Returns the parsed config object, any command-line args that could not be parsed and any error encountered
func LoadConfig() (*Config, []string, error) {
	// check if config file does not exist and create it before proceeding
	var configFileExists bool
	if _, err := os.Stat(AppConfigFilePath); os.IsNotExist(err) {
		configFileExists = createConfigFile()
	} else if !os.IsNotExist(err) {
		configFileExists = true
	}

	// check if any of the command-line args belong in the config file and alert user to set such values in config file only
	if hasConfigFileOption(os.Args) {
		return nil, nil, fmt.Errorf("Unexpected command-line flag/option, "+
			"see godcr -h for supported command-line flags/options"+
			"\nSet other flags/options in %s", AppConfigFilePath)
	}

	// load default config values and create parser object with it
	config := defaultConfig()
	parser := flags.NewParser(&config, flags.IgnoreUnknown)

	// parse command-line args and return any error encountered
	unknownArgs, err := parser.Parse()
	if err != nil {
		return nil, nil, err
	}

	// if config file doesn't exist, no need to attempt to parse and then re-parse command-line args
	if !configFileExists {
		return &config, unknownArgs, nil
	}

	// Load additional config from file
	parser.Options = flags.None
	err = flags.NewIniParser(parser).ParseFile(AppConfigFilePath)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok {
			err = fmt.Errorf("Error parsing configuration file: %s", err.Error())
		}
		return nil, nil, err
	}

	// Parse command line options again to ensure they take precedence.
	parser.Options = flags.IgnoreUnknown
	unknownArgs, err = parser.Parse()

	// return parsed config, unknown args encountered and any error that occurred during last parsing
	return &config, unknownArgs, err
}

// hasConfigFileOption checks if an unknown arg found in command-line is a config file option that should only be set in the config file
func hasConfigFileOption(commandLineArgs []string) bool {
	configFileOptions := configFileOptions()
	isConfigFileOption := func(option string) bool {
		for _, configFileOption := range configFileOptions {
			if strings.EqualFold(configFileOption, option) {
				return true
			}
		}
		return false
	}

	for _, argAndValue := range commandLineArgs {
		arg := strings.Split(argAndValue, "=")[0]
		if isConfigFileOption(strings.TrimSpace(arg)) {
			return true
		}
	}

	return false
}
