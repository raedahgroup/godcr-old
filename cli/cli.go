package cli

import (
	"os"

	"github.com/raedahgroup/dcrcli/walletrpcclient"
)

type response struct {
	columns []string
	result  [][]interface{}
}

var (
	walletClient *walletrpcclient.Client
	stdoutWriter = tabWriter(os.Stdout)
)

// Setup initializes the states of variables used by this package.
func Setup(walletrpcclient *walletrpcclient.Client) {
	walletClient = walletrpcclient
}
