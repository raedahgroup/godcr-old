package utils

import (
	"os"

	"github.com/raedahgroup/dcrcli/core"
)

type Response struct {
	Columns []string
	Result  [][]interface{}
}

var (
	Wallet core.Wallet
	StdoutTabWriter = tabWriter(os.Stdout)
)
