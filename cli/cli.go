package cli

import (
	"os"

	"github.com/raedahgroup/dcrcli/walletrpcclient"
)

type (
	response struct {
		columns []string
		result  [][]interface{}
	}
	// handler carries out the action required by a command.
	// commandArgs holds the arguments passed to the command.
	handler func(walletrpcclient *walletrpcclient.Client, commandArgs []string) (*response, error)

	// cli holds data needed to run the program.
	cli struct {
		funcMap         map[string]handler
		appName         string
		walletrpcclient *walletrpcclient.Client
	}
)

var (
	walletClient    *walletrpcclient.Client
	applicationName string
	stdoutWriter    = tabWriter(os.Stdout)
)

// Setup initializes the states of variables used by this package.
func Setup(walletrpcclient *walletrpcclient.Client, appName string) {
	walletClient = walletrpcclient
	applicationName = appName
}
