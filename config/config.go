package config

import (
	"fmt"
	"github.com/raedahgroup/godcr/cli/termio/terminalprompt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/decred/dcrd/dcrutil"
	flags "github.com/jessevdk/go-flags"
)

const (
	defaultConfigFilename    = "godcr.conf"
	defaultHTTPServerAddress = "127.0.0.1:1234"

	defaultRpcAddress  = "localhost:19111"
	defaultRpcUsername = "rpcuser"
	defaultRpcPassword = "rpcpass"
	ServerAddress      = "127.0.0.1:7778"
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
	ConfigFile        string `short:"C" long:"configfile" description:"Path to configuration file"`
	TestNet           bool   `short:"t" long:"testnet" description:"Connects to testnet wallet instead of mainnet"`
	RPCUser           string `short:"u" long:"rpcuser" description:"RPC username"`
	RPCPassword       string `short:"p" long:"rpcpass" default-mask:"-" description:"RPC password"`
	WalletRPCServer   string `short:"w" long:"walletrpcserver" description:"Wallet RPC server to connect to"`
	RPCCert           string `short:"c" long:"rpccert" description:"RPC server certificate chain for validation"`
	HTTPServerAddress string `short:"s" long:"serveraddress" description:"Address and port of the HTTP server."`
	HTTPMode          bool   `long:"http" description:"Run in HTTP mode. This flag cannot be used with a command."`
	NoDaemonTLS       bool   `long:"nodaemontls" description:"Disable TLS"`
	Init              bool   `long:"init" description:"Create the config file"`
}

// defaultConfig an instance of Config with the defaults set.
func defaultConfig() Config {
	return Config{
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
	}

	//config
	if config.Init {
		err = configure(config)
		if err != nil {
			fmt.Println(err)
		}
		os.Exit(0)
	}

	// Load additional config from file
	err = parseConfigFile(parser, config.ConfigFile)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok {
			return args, config, parser, fmt.Errorf("Error parsing configuration file: %v", err.Error())
		}
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
}


func configure(cfg Config) (err error) {
	inputErr := func(error) error {
		return fmt.Errorf("error receiving input: %s", err.Error())
	}

	var rpcAddress, rpcUsername, rpcPassword string

	// Todo load config from dcrwallet to auto collect the values below instead of requesting user input
	rpcAddress, err = terminalprompt.RequestInput(fmt.Sprintf("Wallet RPC server (default = %s)", defaultRpcAddress), nil)
	if err != nil {
		return inputErr(err)
	}

	if rpcAddress == "" {
		rpcAddress = defaultRpcAddress
	}

	rpcUsername, err = terminalprompt.RequestInput(fmt.Sprintf("RPC username (default = %s)", defaultRpcUsername), nil)
	if err != nil {
		return inputErr(err)
	}
	if rpcUsername == "" {
		rpcUsername = defaultRpcUsername
	}

	rpcPassword, err = terminalprompt.RequestInputSecure(fmt.Sprintf("RPC password (default = %s)", defaultRpcPassword), nil)
	if err != nil {
		return inputErr(err)
	}
	if rpcPassword == "" {
		rpcPassword = defaultRpcPassword
	}

	config := struct {
		RpcAddress    string
		RpcUsername   string
		RpcPassword   string
		ServerAddress string
	}{
		RpcAddress: rpcAddress, RpcUsername: rpcUsername, RpcPassword: rpcPassword,
		ServerAddress: ServerAddress,
	}

	f, err := os.Create(cfg.ConfigFile)
	if err != nil {
		return fmt.Errorf("error in creating config file: %v", err)
	}
	defer f.Close()

	tmpl := template.New("config")

	tmpl, err = tmpl.Parse(configText())
	if err != nil {
		return fmt.Errorf("error in parsing the config template: %v", err.Error())
	}

	err = tmpl.Execute(f, config)
	if err != nil {
		return err
	}

	fmt.Println("\n\nThe config file has been set successfully\n")
	return nil
}
