package nuklear

import (
	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/nuklear/handlers"
)

type Handler interface {
	BeforeRender()
	Render(*nucular.Window, app.WalletMiddleware, func(page string))
}

type handlersData struct {
	name       string
	navLabel   string
	standalone bool
	handler    Handler
}

func getHandlers() []handlersData {
	return []handlersData{
		{
			name:       "balance",
			navLabel:   "Balance",
			standalone: false,
			handler:    &handlers.BalanceHandler{},
		},
		{
			name:       "receive",
			navLabel:   "Receive",
			standalone: false,
			handler:    &handlers.ReceiveHandler{},
		},
		{
			name:       "send",
			navLabel:   "Send (WIP)",
			standalone: false,
			handler:    &handlers.SendHandler{},
		},
		{
			name:     "history",
			navLabel: "History",
			handler:  &handlers.TransactionsHandler{},
		},
		{
			name:       "createwallet",
			standalone: true,
			handler:    &handlers.CreateWalletHandler{},
		},
		{
			name:       "stakeinfo",
			navLabel:   "Stake Info",
			standalone: false,
			handler:    &handlers.StakeInfoHandler{},
		},
		{
			name:       "purchasetickets",
			navLabel:   "Purchase Tickets",
			standalone: false,
			handler:    &handlers.PurchaseTicketsHandler{},
		},
	}
}
