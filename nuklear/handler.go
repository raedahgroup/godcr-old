package nuklear

import (
	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/handlers"
)

type navPage struct {
	name    string
	label   string
	handler navPageHandler
}

type navPageHandler interface {
	BeforeRender()
	Render(*nucular.Window, walletcore.Wallet)
}

type standalonePage struct {
	name    string
	handler standalonePageHandler
}

type standalonePageHandler interface {
	BeforeRender()
	Render(*nucular.Window, app.WalletMiddleware, func(string))
}

func getNavPages() []navPage {
	return []navPage{
		{
			name:    "overview",
			label:   "Overview",
			handler: &handlers.OverviewHandler{},
		},
		{
			name:    "history",
			label:   "History",
			handler: &handlers.HistoryHandler{},
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
