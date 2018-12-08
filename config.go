package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/raedahgroup/dcrcli/cli"

	flags "github.com/btcsuite/go-flags"
	"github.com/decred/dcrd/dcrutil"
)

const (
	defaultConfigFilename = "dcrcli.conf"
	defaultLogDirname     = "log"
	defaultLogFilename    = "dcrcli.log"
)

var (
	defaultAppDataDir          = dcrutil.AppDataDir("dcrcli", false)
	defaultDcrwalletAppDataDir = dcrutil.AppDataDir("dcrwallet", false)
	defaultRPCCertFile         = filepath.Join(defaultDcrwalletAppDataDir, "rpc.cert")
	defaultConfigFile          = filepath.Join(defaultAppDataDir, defaultConfigFilename)
	defaultLogDir              = filepath.Join(defaultAppDataDir, defaultLogDirname)
)

type config struct {
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

func cleanAndExpandPath(path string) string {
	// Do not try to clean the empty string
	if path == "" {
		return ""
	}

	// NOTE: The os.ExpandEnv doesn't work with Windows cmd.exe-style
	// %VARIABLE%, but they variables can still be expanded via POSIX-style
	// $VARIABLE.
	path = os.ExpandEnv(path)

	if !strings.HasPrefix(path, "~") {
		return filepath.Clean(path)
	}

	// Expand initial ~ to the current user's home directory, or ~otheruser
	// to otheruser's home directory.  On Windows, both forward and backward
	// slashes can be used.
	path = path[1:]

	var pathSeparators string
	if runtime.GOOS == "windows" {
		pathSeparators = string(os.PathSeparator) + "/"
	} else {
		pathSeparators = string(os.PathSeparator)
	}

	userName := ""
	if i := strings.IndexAny(path, pathSeparators); i != -1 {
		userName = path[:i]
		path = path[i:]
	}

	homeDir := ""
	var u *user.User
	var err error
	if userName == "" {
		u, err = user.Current()
	} else {
		u, err = user.Lookup(userName)
	}
	if err == nil {
		homeDir = u.HomeDir
	}
	// Fallback to CWD if user lookup fails or user has no home directory.
	if homeDir == "" {
		homeDir = "."
	}

	return filepath.Join(homeDir, path)
}

func addParserSettings(parser *flags.Parser) {
	parser.Usage = cli.HelpMessage()
	parser.UnknownOptionHandler = func(option string, arg flags.SplitArgument, args []string) ([]string, error) {
		return nil, fmt.Errorf("unknown option %s", option)
	}
}

func loadConfig() (*config, []string, error) {
	cfg := config{
		ConfigFile: defaultConfigFile,
		RPCCert:    defaultRPCCertFile,
	}
	// Pre-parse command line arguments
	preCfg := cfg
	preParser := flags.NewParser(&preCfg, flags.HelpFlag)
	addParserSettings(preParser)

	_, err := preParser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); ok && e.Type != flags.ErrHelp {
			cli.PrintHelp()
			os.Exit(1)
		} else if ok && e.Type == flags.ErrHelp {
			preParser.WriteHelp(os.Stderr)
			os.Exit(0)
		}
	}

	// Show version and exit if the version flag was specified
	appName := filepath.Base(os.Args[0])
	appName = strings.TrimSuffix(appName, filepath.Ext(appName))
	if preCfg.ShowVersion {
		fmt.Println(appName, "version", Ver.String())
		os.Exit(0)
	}

	// Load additional config from file
	parser := flags.NewParser(&cfg, flags.Default)
	addParserSettings(parser)

	err = flags.NewIniParser(parser).ParseFile(preCfg.ConfigFile)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok {
			fmt.Fprintf(os.Stderr, "Error parsing config file: %v\n",
				err)
			parser.WriteHelp(os.Stderr)
			return nil, nil, err
		}
	}

	// Parse command line options again to ensure they take precedence.
	remainingArgs, err := parser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); !ok || e.Type != flags.ErrHelp {
			cli.PrintHelp()
		}
		return nil, nil, err
	}

	return &cfg, remainingArgs, nil
}
