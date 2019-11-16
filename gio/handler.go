package gio

import (
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/raedahgroup/dcrlibwallet"
)

type handler interface {
	BeforeRender(*dcrlibwallet.LibWallet)
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

func (_ *notImplementedNavPageHandler) BeforeRender(_ *dcrlibwallet.LibWallet) {

}

func (p *notImplementedNavPageHandler) Render(_ *layout.Context, _ func()) {

}
