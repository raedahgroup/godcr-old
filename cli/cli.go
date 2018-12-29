package cli

import (
	"os"

	"github.com/raedahgroup/godcr/walletrpcclient"
)

type Response struct {
	Columns []string
	Result  [][]interface{}
}

var (
	WalletClient *walletrpcclient.Client
	StdoutWriter = tabWriter(os.Stdout)
)
