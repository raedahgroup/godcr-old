package routes

import (
	"html/template"
	"time"

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
		{"staking.html", "web/views/staking.html"},
	}
}

func templateFuncMap() template.FuncMap {
	return template.FuncMap{
		"simpleBalance": func(balance *walletcore.Balance, detailed bool) string {
			result := walletcore.SimpleBalance(balance, detailed)
			return walletcore.NormalizeBalance(result)
		},
		"spendableBalance": func(balance *walletcore.Balance) string {
			result := walletcore.SpendableBalance(balance)
			return walletcore.NormalizeBalance(result)
		},
		"accountString": func(account *walletcore.Account) string {
			return account.String()
		},
		"amountDcr": func(amount int64) string {
			return dcrutil.Amount(amount).String()
		},
		"timestamp": func() int64 {
			return time.Now().Unix()
		},
	}
}
