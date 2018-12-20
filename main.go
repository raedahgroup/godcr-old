package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jessevdk/go-flags"

	"github.com/raedahgroup/dcrcli/cli"
	"github.com/raedahgroup/dcrcli/walletrpcclient"
	"github.com/raedahgroup/dcrcli/web"
)

type Version struct {
	Major, Minor, Patch int
	Label               string
	Nick                string
}

var Ver = Version{
	Major: 0,
	Minor: 0,
	Patch: 1,
	Label: "",
}

// CommitHash may be set on the build command line:
// go build -ldflags "-X github.com/decred/dcrdata/version.CommitHash=`git describe --abbrev=8 --long | awk -F "-" '{print $(NF-1)"-"$NF}'`"
var CommitHash string

func (v *Version) String() string {
	var hashStr string
	if CommitHash != "" {
		hashStr = "+" + CommitHash
	}
	if v.Label != "" {
		return fmt.Sprintf("%d.%d.%d-%s%s",
			v.Major, v.Minor, v.Patch, v.Label, hashStr)
	}
	return fmt.Sprintf("%d.%d.%d%s",
		v.Major, v.Minor, v.Patch, hashStr)
}

func main() {
	appName := filepath.Base(os.Args[0])
	appName = strings.TrimSuffix(appName, filepath.Ext(appName))

	config, _, err := loadConfig(appName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if config == nil {
		os.Exit(0)
	}

	// Show version and exit if the version flag was specified
	if config.ShowVersion {
		fmt.Fprintf(os.Stdout, "%s version: %s\n", cli.AppName(), Ver.String())
		os.Exit(0)
	}

	if config.HTTPMode {
		fmt.Println("Running in http mode")
		enterHTTPMode(config)
	} else {
		enterCliMode(appName, config)
	}
}

func enterHTTPMode(config *cli.Config) {
	client, err := walletrpcclient.New(config.WalletRPCServer, config.RPCCert, config.NoDaemonTLS)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error connecting to RPC server")
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	web.StartHttpServer(config.HTTPServerAddress, client)
}

func enterCliMode(appName string, config *cli.Config) {
	client, err := walletrpcclient.New(config.WalletRPCServer, config.RPCCert, config.NoDaemonTLS)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error connecting to RPC server")
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	cli.Setup(client)
	parser := flags.NewParser(&cli.DcrcliCommands, flags.Default)
	if _, err := parser.Parse(); err != nil {
		handleParseError(err, parser)
		os.Exit(1)
	}
}

func loadConfig(appName string) (*cli.Config, []string, error) {
	// load defaults first
	cfg := cli.DefaultConfig()

	// Load additional config from file
	parser := flags.NewParser(&cfg, flags.IgnoreUnknown)

	// Errors here will be owing to unknown flags and commands, which are not of interest at this point.
	parser.Parse()

	err := flags.NewIniParser(parser).ParseFile(cfg.ConfigFile)
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

func handleParseError(err error, helpParser *flags.Parser) {
	if (helpParser.Options & flags.PrintErrors) != flags.None {
		// error printing is already handled by go-flags.
		return
	}
	if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrHelp {
		helpParser.WriteHelp(os.Stderr)
	} else {
		fmt.Println(err)
	}
}
