package gio

import (
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/gio/helper"
	"github.com/raedahgroup/godcr/gio/pages/wallet"
)

type standalonePageHandler interface {
	//BeforeRender(*dcrlibwallet.MultiWallet)
	Render(ctx *layout.Context, refreshWindowFunc func(), changePageFunc func(string))
}

type navPageHandler interface {
	BeforeRender(*dcrlibwallet.MultiWallet)
	Render(ctx *layout.Context, refreshWindowFunc func())
}

type navPage struct {
	name      string
	label     string
	button    *widget.Button
	handler   navPageHandler
}

func getStandalonePages(multiWallet *dcrlibwallet.MultiWallet, theme *helper.Theme) map[string]standalonePageHandler {
	return map[string]standalonePageHandler{
		"welcome": wallet.NewWelcomePage(multiWallet, theme),
	}
}

func getNavPages() []navPage {
	return []navPage{
		{
			name:      "overview",
			label:     "Overview",
			button:    new(widget.Button),
			handler:   &notImplementedNavPageHandler{"Overview"},
		},
		{
			name:      "history",
			label:     "History",
			button:    new(widget.Button),
			handler:   &notImplementedNavPageHandler{"History"},
		},
		{
			name:      "send",
			label:     "Send",
			button:    new(widget.Button),
			handler:   &notImplementedNavPageHandler{"Send"},
		},
		{
			name:      "receive",
			label:     "Receive",
			button:    new(widget.Button),
			handler:   &notImplementedNavPageHandler{"Receive"},
		},
		{
			name:      "staking",
			label:     "Staking",
			button:    new(widget.Button),
			handler:   &notImplementedNavPageHandler{"Staking"},
		},
		{
			name:      "security",
			label:     "Security",
			button:    new(widget.Button),
			handler:   &notImplementedNavPageHandler{"Security"},
		},
		{
			name:      "settings",
			label:     "Settings",
			button:    new(widget.Button),
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
