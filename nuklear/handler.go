package nuklear

import (
	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/nuklear/handlers"
)

type NavPageHandler interface {
	BeforeRender()
	Render(*nucular.Window, app.WalletMiddleware, func(page string))
}

type StandalonePageHandler interface {
	BeforeRender()
	Render(*nucular.Window, app.WalletMiddleware, func(string))
}

type navPageData struct {
	name    string
	label   string
	handler NavPageHandler
}

type standalonePageData struct {
	name    string
	handler StandalonePageHandler
}

func getNavPagesData() []navPageData {
	return []navPageData{
		{
			name:    "balance",
			label:   "Balance",
			handler: &handlers.BalanceHandler{},
		},
		{
			name:    "receive",
			label:   "Receive",
			handler: &handlers.ReceiveHandler{},
		},
		{
			name:    "send",
			label:   "Send (WIP)",
			handler: &handlers.SendHandler{},
		},
		{
			name:    "history",
			label:   "History",
			handler: &handlers.TransactionsHandler{},
		},
		{
			name:    "stakeinfo",
			label:   "Stake Info",
			handler: &handlers.StakeInfoHandler{},
		},
		{
			name:    "purchasetickets",
			label:   "Purchase Tickets",
			handler: &handlers.PurchaseTicketsHandler{},
		},
	}
}

func getStandalonePageData() []standalonePageData {
	return []standalonePageData{
		{
			name:    "sync",
			handler: &handlers.SyncHandler{},
		},
	}
}
