package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/raedahgroup/dcrcli/cli"

	flags "github.com/btcsuite/go-flags"
	"github.com/decred/dcrd/dcrutil"
)

const (
	defaultConfigFilename    = "dcrcli.conf"
	defaultLogDirname        = "log"
	defaultLogFilename       = "dcrcli.log"
	defaultHTTPServerAddress = "127.0.0.1:1234"
)

var (
	defaultAppDataDir          = dcrutil.AppDataDir("dcrcli", false)
	defaultDcrwalletAppDataDir = dcrutil.AppDataDir("dcrwallet", false)
	defaultRPCCertFile         = filepath.Join(defaultDcrwalletAppDataDir, "rpc.cert")
	defaultConfigFile          = filepath.Join(defaultAppDataDir, defaultConfigFilename)
	defaultLogDir              = filepath.Join(defaultAppDataDir, defaultLogDirname)
)

type Config struct {
	ShowVersion       bool   `short:"v" long:"version" description:"Display version information and exit"`
	ConfigFile        string `short:"C" long:"configfile" description:"Path to configuration file"`
	RPCUser           string `short:"u" long:"rpcuser" description:"RPC username"`
	RPCPassword       string `short:"p" long:"rpcpass" default-mask:"-" description:"RPC password"`
	WalletRPCServer   string `short:"w" long:"walletrpcserver" description:"Wallet RPC server to connect to"`
	RPCCert           string `short:"c" long:"rpccert" description:"RPC server certificate chain for validation"`
	HTTPServerAddress string `short:"s" long:"serveraddress" description:"Address and port of the HTTP server."`
	HTTPMode          bool   `long:"http" description:"Run in HTTP mode."`
	NoDaemonTLS       bool   `long:"nodaemontls" description:"Disable TLS"`
}

type AppCommands struct {
	Config
	Balance    cli.BalanceCommand    `command:"balance" description:"show your balance"`
	Send       cli.SendCommand       `command:"send" description:"send a transaction"`
	SendCustom cli.SendCustomCommand `command:"send-custom" description:"send a transaction, manually selecting inputs from unspent outputs"`
	Receive    cli.ReceiveCommand    `command:"receive" description:"show your address to receive funds"`
	History    cli.HistoryCommand    `command:"history" description:"show your transaction history"`
}

var appCommands AppCommands

func loadConfig(appName string) (*Config, []string, error) {
	// load defaults first
	cfg := Config{
		ConfigFile:        defaultConfigFile,
		RPCCert:           defaultRPCCertFile,
		HTTPServerAddress: defaultHTTPServerAddress,
	}

	// Pre-parse command line arguments.
	//
	// separate help parser (used for displaying help) from config parser.
	// This is to prevent triggering the execution of any command encountered: the application is not
	// fully initialized at this point.
	preParser := flags.NewParser(&cfg, flags.HelpFlag)
	helpParser := flags.NewParser(&AppCommands{Config: cfg}, flags.HelpFlag)

	_, err := preParser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); ok && e.Type != flags.ErrHelp {
			os.Exit(1)
		} else if ok && e.Type == flags.ErrHelp {
			helpParser.WriteHelp(os.Stderr)
			os.Exit(0)
		}
	}

	// Show version and exit if the version flag was specified
	if cfg.ShowVersion {
		fmt.Println(appName, "version", Ver.String())
		os.Exit(0)
	}

	// Load additional config from file
	parser := flags.NewParser(&cfg, flags.Default)

	err = flags.NewIniParser(parser).ParseFile(cfg.ConfigFile)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok {
			return nil, nil, fmt.Errorf("Error parsing configuration file: %v", err.Error())
		}
		return nil, nil, err
	}

	// Parse command line options again to ensure they take precedence.
	remainingArgs, err := parser.Parse()
	if err != nil {
		return nil, nil, err
	}

	return &cfg, remainingArgs, nil
}
