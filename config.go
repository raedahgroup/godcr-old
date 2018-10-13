package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	flags "github.com/btcsuite/go-flags"
	"github.com/decred/dcrd/dcrutil"
)

const (
	defaultConfigFilename = "dcrcli.conf"
	defaultLogDirname     = "log"
	defaultLogFilename    = "dcrcli.log"
)

var (
	defaultAppDataDir = dcrutil.AppDataDir("dcrcli", false)
	defaultConfigFile = filepath.Join(defaultAppDataDir, defaultConfigFilename)
	defaultLogDir     = filepath.Join(defaultAppDataDir, defaultLogDirname)
)

type config struct {
	ShowVersion       bool   `short:"v" long:"version" description:"Display version information and exit"`
	ListCommands      bool   `short:"l" long:"listcommands" description:"List all of the supported commands and exit"`
	ConfigFile        string `short:"C" long:"configfile" description:"Path to configuration file"`
	RPCUser           string `short:"u" long:"rpcuser" description:"RPC username"`
	RPCPassword       string `short:"P" long:"rpcpass" default-mask:"-" description:"RPC password"`
	WalletRPCServer   string `short:"w" long:"walletrpcserver" description:"Wallet RPC server to connect to"`
	RPCCert           string `short:"c" long:"rpccert" description:"RPC server certificate chain for validation"`
	HTTPServerAddress string `short:"h" long:"httpserveraddress" description:"Serve via http using this address if not empty"`
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

func loadConfig() (*config, []string, error) {
	// Default config
	cfg := config{
		ConfigFile: defaultConfigFile,
	}

	// Pre-parse the command line options to see if an alternative config
	// file, the version flag, or the list commands flag was specified.  Any
	// errors aside from the help message error can be ignored here since
	// they will be caught by the final parse below.
	preCfg := cfg
	preParser := flags.NewParser(&preCfg, flags.Default)
	_, err := preParser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); ok && e.Type != flags.ErrHelp {
			fmt.Fprintln(os.Stderr, err)
			fmt.Fprintln(os.Stderr, "")
			fmt.Fprintln(os.Stderr, "The special parameter `-` "+
				"indicates that a parameter should be read "+
				"from the\nnext unread line from standard input.")
			os.Exit(1)
		} else if ok && e.Type == flags.ErrHelp {
			fmt.Fprintln(os.Stdout, err)
			fmt.Fprintln(os.Stdout, "")
			fmt.Fprintln(os.Stdout, "The special parameter `-` "+
				"indicates that a parameter should be read "+
				"from the\nnext unread line from standard input.")
			os.Exit(0)
		}
	}

	// Show version and exit if the version flag was specified.
	appName := filepath.Base(os.Args[0])
	appName = strings.TrimSuffix(appName, filepath.Ext(appName))
	usageMessage := fmt.Sprintf("Use %s -h to show options", appName)
	if preCfg.ShowVersion {
		fmt.Println(appName, "version", Ver.String())
		os.Exit(0)
	}

	// Show available commands and exit if list commands flag was specified.
	if preCfg.ListCommands {
		return &cfg, []string{"listcommands"}, nil
	}

	// Load additional config from file
	parser := flags.NewParser(&cfg, flags.Default)
	err = flags.NewIniParser(parser).ParseFile(preCfg.ConfigFile)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok {
			fmt.Fprintf(os.Stderr, "Error parsing config file: %v\n",
				err)
			fmt.Fprintln(os.Stderr, usageMessage)
			return nil, nil, err
		}
	}

	// Parse command line options again to ensure they take precedence.
	remainingArgs, err := parser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); !ok || e.Type != flags.ErrHelp {
			fmt.Fprintln(os.Stderr, usageMessage)
		}
		return nil, nil, err
	}

	return &cfg, remainingArgs, nil
}
