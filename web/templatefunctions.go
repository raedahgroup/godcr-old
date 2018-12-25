package web

import (
	"fmt"
	"text/template"

	"github.com/raedahgroup/dcrcli/app/walletcore"
)

func templateFuncMap() template.FuncMap {
	return template.FuncMap{
		"txExplorerLink": func(tx walletcore.Transaction) string {
			if tx.Testnet {
				return fmt.Sprintf("https://testnet.dcrdata.org/tx/%s", tx.Hash)
			} else {
				return fmt.Sprintf("https://mainnet.dcrdata.org/tx/%s", tx.Hash)
			}
		},
	}
}
