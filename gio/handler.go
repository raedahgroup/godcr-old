package gio

import (
	"github.com/raedahgroup/godcr/gio/handlers"
	"gioui.org/widget"

)

type handler interface {
	BeforeRender()
	Render()
}

type page struct {
	name    string
	label   string
	button *widget.Button
	handler handler
}

func getPages() []page {
	return []page{
		{
			name:    "overview",
			label:   "Overview",
			button: new(widget.Button),
			handler: handlers.NewOverviewHandler(),
		},
		{
			name:    "history",
			label:   "History",
			button: new(widget.Button),
			handler: handlers.NewOverviewHandler(),
		},
		{
			name:    "send",
			label:   "Send",
			button: new(widget.Button),
			handler: handlers.NewOverviewHandler(),
		},
		{
			name:    "receive",
			label:   "Receive",
			button: new(widget.Button),
			handler: handlers.NewOverviewHandler(),
		},
		{
			name:    "staking",
			label:   "Staking",
			button: new(widget.Button),
			handler: handlers.NewOverviewHandler(),
		},
		{
			name:    "security",
			label:   "Security",
			button: new(widget.Button),
			handler: handlers.NewOverviewHandler(),
		},
		{
			name:    "settings",
			label:   "Settings",
			button: new(widget.Button),
			handler: handlers.NewOverviewHandler(),
		},
	}
}
