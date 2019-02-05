package routes

import (
	"html/template"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/godcr/app/walletcore"
)

type templateData struct {
	name string
	path string
}

func templates() []templateData {
	return []templateData{
		{"error.html", "web/views/error.html"},
		{"createwallet.html", "web/views/createwallet.html"},
		{"balance.html", "web/views/balance.html"},
		{"send.html", "web/views/send.html"},
		{"receive.html", "web/views/receive.html"},
		{"history.html", "web/views/history.html"},
		{"transaction_details.html", "web/views/transaction_details.html"},
	}
}

func templateFuncMap() template.FuncMap {
	return template.FuncMap{
		"simpleBalance": walletcore.SimpleBalance,
		"amountDcr": func(amount int64) string {
			return dcrutil.Amount(amount).String()
		},
	}
}
