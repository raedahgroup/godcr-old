package nuklear

import (
	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/pagehandlers"
	"github.com/raedahgroup/godcr/nuklear/widgets"
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

type notImplementedNavPageHandler struct{}

func (_ *notImplementedNavPageHandler) BeforeRender() {}
func (_ *notImplementedNavPageHandler) Render(window *nucular.Window, _ walletcore.Wallet) {
	w := widgets.Window{window}
	w.DisplayErrorMessage("Page not yet implemented")
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
		{
			name:    "accounts",
			label:   "Accounts",
			handler: &notImplementedNavPageHandler{},
		},
		{
			name:    "security",
			label:   "Security",
			handler: &notImplementedNavPageHandler{},
		},
		{
			name:    "settings",
			label:   "Settings",
			handler: &notImplementedNavPageHandler{},
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
