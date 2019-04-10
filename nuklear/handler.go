package nuklear

import (
	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/pagehandlers"
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
	Render(*nucular.Window, app.WalletMiddleware, func(*nucular.Window, string))
}

func getNavPages() []navPage {
	return []navPage{
		{
			name:    "overview",
			label:   "Overview",
			handler: &pagehandlers.OverviewHandler{},
		},
		{
			name:    "history",
			label:   "History",
			handler: &pagehandlers.HistoryHandler{},
		},
		{
			name:    "send",
			label:   "Send",
			handler: &pagehandlers.SendHandler{},
		},
		{
			name:    "receive",
			label:   "Receive",
			handler: &pagehandlers.ReceiveHandler{},
		},
		{
			name:    "staking",
			label:   "Staking",
			handler: &pagehandlers.StakingHandler{},
		},
	}
}

func getStandalonePages() []standalonePage {
	return []standalonePage{
		{
			name:    "sync",
			handler: &pagehandlers.SyncHandler{},
		},
		{
			name:    "createwallet",
			handler: &pagehandlers.CreateWalletHandler{},
		},
	}
}
