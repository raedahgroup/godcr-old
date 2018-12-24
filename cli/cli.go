package cli

import (
	"os"

	"github.com/jessevdk/go-flags"

	"github.com/raedahgroup/dcrcli/walletrpcclient"
)

type Response struct {
	Columns []string
	Result  [][]interface{}
}

var (
	WalletClient *walletrpcclient.Client
	StdoutWriter = tabWriter(os.Stdout)
)

// IsFlagErrorType determines whether a given error is of a given flags.ErrorType.
// It is safe to call IsFlagErrorType with err = nil.
func IsFlagErrorType(err error, errorType flags.ErrorType) bool {
	if err == nil {
		return false
	}
	if flagErr, ok := err.(*flags.Error); ok && flagErr.Type == errorType {
		return true
	}
	return false
}
