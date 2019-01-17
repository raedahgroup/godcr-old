package nuklear

import (
	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/handlers"
)

type Handler interface {
	BeforeRender()
	Render(*nucular.Window)
	SetWalletMiddleware(walletcore.Wallet)
}

type handlersData struct {
	name     string
	navLabel string
	handler  Handler
}

func getHandlers() []handlersData {
	return []handlersData{
		{
			name:     "balance",
			navLabel: "Balance",
			handler:  &handlers.BalanceHandler{},
		},
		{
			name:     "receive",
			navLabel: "Receive",
			handler:  &handlers.ReceiveHandler{},
		},
		{
			name:     "send",
			navLabel: "Send (WIP)",
			handler:  &handlers.SendHandler{},
		},
		{
			name:     "transactions",
			navLabel: "Transactions",
			handler:  &handlers.TransactionsHandler{},
		},
	}
}
