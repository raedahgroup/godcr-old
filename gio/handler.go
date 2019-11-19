package gio

import (
	"gioui.org/layout"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/gio/widgets"
	"github.com/raedahgroup/godcr/gio/pages/wallet"
)

type standalonePageHandler interface {
	Render(ctx *layout.Context, refreshWindowFunc func(), changePageFunc func(string))
}

type navPageHandler interface {
	BeforeRender(*dcrlibwallet.MultiWallet)
	Render(ctx *layout.Context, refreshWindowFunc func())
}

type navPage struct {
	name      string
	label     string
	button    *widgets.Button
	handler   navPageHandler
}

func getStandalonePages(multiWallet *dcrlibwallet.MultiWallet) map[string]standalonePageHandler {
	return map[string]standalonePageHandler{
		"welcome"      : wallet.NewWelcomePage(multiWallet),
		"createwallet" : wallet.NewCreateWalletPage(multiWallet),
		"restorewallet": wallet.NewRestoreWalletPage(multiWallet),
	}
}

func getNavPages() []navPage {
	return []navPage{
		{
			name:      "overview",
			label:     "Overview",
			button:     widgets.NewButton("Overview", nil),
			handler:   &notImplementedNavPageHandler{"Overview"},
		},
		{
			name:      "history",
			label:     "History",
			button:   	widgets.NewButton("History", nil),
			handler:   &notImplementedNavPageHandler{"History"},
		},
		{
			name:      "send",
			label:     "Send",
			button:   	widgets.NewButton("Send", nil),
			handler:   &notImplementedNavPageHandler{"Send"},
		},
		{
			name:      "receive",
			label:     "Receive",
			button:   	widgets.NewButton("Receive", nil),
			handler:   &notImplementedNavPageHandler{"Receive"},
		},
		{
			name:      "staking",
			label:     "Staking",
			button:   	widgets.NewButton("Staking",  nil),
			handler:   &notImplementedNavPageHandler{"Staking"},
		},
		{
			name:      "security",
			label:     "Security",
			button:   	widgets.NewButton("Security", nil),
			handler:   &notImplementedNavPageHandler{"Security"},
		},
		{
			name:      "settings",
			label:     "Settings",
			button:   	widgets.NewButton("Settings", nil),
			handler:   &notImplementedNavPageHandler{"Settings"},
		},
	}
}

type notImplementedNavPageHandler struct {
	pageTitle string
}

func (_ *notImplementedNavPageHandler) BeforeRender(_ *dcrlibwallet.MultiWallet) {

}

func (p *notImplementedNavPageHandler) Render(_ *layout.Context, _ func()) {

}
