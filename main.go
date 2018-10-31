package main

import (
	"fmt"
	"os"

	"github.com/raedahgroup/dcrcli/cli"
	"github.com/raedahgroup/dcrcli/server"
	"github.com/raedahgroup/dcrcli/walletrpcclient"
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

const AppName string = "dcrcli"

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
	config, args, err := loadConfig()
	if err != nil {
		os.Exit(1)
	}

	if config.Mode == "http" {
		fmt.Println("Running in http mode")
		enterHttpMode(config)
	} else {
		enterCliMode(config, args)
	}
}

func enterHttpMode(config *config) {
	if config.HTTPServerAddress == "" {
		fmt.Println("Cannot start http server. Server address not set")
		os.Exit(1)
	}

	client, err := walletrpcclient.New(config.WalletRPCServer, config.RPCCert, config.NoDaemonTLS)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error connecting to RPC server")
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	server.StartHttpServer(config.HTTPServerAddress, client)
}

func enterCliMode(config *config, args []string) {
	client, err := walletrpcclient.New(config.WalletRPCServer, config.RPCCert, config.NoDaemonTLS)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error connecting to RPC server")
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	cli := cli.New(client, AppName)
	cli.RunCommand(args)
}
