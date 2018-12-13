package cli

import (
	"os"

	ws "github.com/raedahgroup/dcrcli/walletsource"
)

type Response struct {
	Columns []string
	Result  [][]interface{}
}

var (
	WalletSource ws.WalletSource
	StdoutWriter = tabWriter(os.Stdout)
)
