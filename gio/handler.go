package gio

import (
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/app/walletcore"
)

type handler interface {
	BeforeRender(walletcore.Wallet, *config.Settings)
	Render(ctx *layout.Context, refreshWindowFunc func())
}

type page struct {
	name      string
	label     string
	button    *widget.Button
	isNavPage bool
	handler   handler
}

func getPages() []page {
	return []page{
		{
			name:      "overview",
			label:     "Overview",
			button:    new(widget.Button),
			isNavPage: true,
			handler:   &notImplementedNavPageHandler{"Overview"},
		},
		{
			name:      "history",
			label:     "History",
			button:    new(widget.Button),
			isNavPage: true,
			handler:   &notImplementedNavPageHandler{"History"},
		},
		{
			name:      "send",
			label:     "Send",
			button:    new(widget.Button),
			isNavPage: true,
			handler:   &notImplementedNavPageHandler{"Send"},
		},
		{
			name:      "receive",
			label:     "Receive",
			button:    new(widget.Button),
			isNavPage: true,
			handler:   &notImplementedNavPageHandler{"Receive"},
		},
		{
			name:      "staking",
			label:     "Staking",
			button:    new(widget.Button),
			isNavPage: true,
			handler:   &notImplementedNavPageHandler{"Staking"},
		},
		{
			name:      "security",
			label:     "Security",
			button:    new(widget.Button),
			isNavPage: true,
			handler:   &notImplementedNavPageHandler{"Security"},
		},
		{
			name:      "settings",
			label:     "Settings",
			button:    new(widget.Button),
			isNavPage: true,
			handler:   &notImplementedNavPageHandler{"Settings"},
		},
	}
}

type notImplementedNavPageHandler struct {
	pageTitle string
}

func (_ *notImplementedNavPageHandler) BeforeRender(_ walletcore.Wallet, _ *config.Settings) {

}

func (p *notImplementedNavPageHandler) Render(_ *layout.Context, _ func()) {

}
