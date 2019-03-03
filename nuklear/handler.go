package nuklear

import (
	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/nuklear/handlers"
)

type navPageHandler interface {
	BeforeRender()
	Render(*nucular.Window, app.WalletMiddleware)
}

type standalonePageHandler interface {
	BeforeRender()
	Render(*nucular.Window, app.WalletMiddleware, func(string))
}

type navPage struct {
	name    string
	label   string
	handler navPageHandler
}

type standalonePage struct {
	name    string
	handler standalonePageHandler
}

func getNavPages() []navPage {
	return []navPage{
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

func getStandalonePages() []standalonePage {
	return []standalonePage{
		{
			name:    "sync",
			handler: &handlers.SyncHandler{},
		},
		{
			name: "createwallet",
			handler: &handlers.CreateWalletHandler{},
		}
	}
}
