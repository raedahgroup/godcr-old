package nuklear

import (
	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/handlers"
)

type navPageHandler interface {
	BeforeRender()
	Render(*nucular.Window, walletcore.Wallet)
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
			name:    "send",
			label:   "Send",
			handler: &handlers.SendHandler{},
		},
		{
			name:    "receive",
			label:   "Receive",
			handler: &handlers.ReceiveHandler{},
		},
		{
			name:    "history",
			label:   "History",
			handler: &handlers.HistoryHandler{},
		},
		{
			name:    "staking",
			label:   "Staking",
			handler: &handlers.StakingHandler{},
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
			name:    "createwallet",
			handler: &handlers.CreateWalletHandler{},
		},
	}
}
