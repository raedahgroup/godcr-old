package routes

import (
	"fmt"
	"html/template"

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
	}
}

func templateFuncMap() template.FuncMap {
	return template.FuncMap{
		"simpleBalance": func(balance *walletcore.Balance, detailed bool) string {
			if detailed || balance.Total == balance.Spendable {
				return balance.Total.String()
			} else {
				return fmt.Sprintf("Total %s (Spendable %s)", balance.Total.String(), balance.Spendable.String())
			}
		},
	}
}
