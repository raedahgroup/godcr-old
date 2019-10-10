package gio

import (
	"github.com/raedahgroup/godcr/gio/handlers"
	"github.com/raedahgroup/godcr/gio/widget"
)

type handler interface {
	BeforeRender()
	Render()
}

type page struct {
	name    string
	label   string
	clicker int
	handler handler
}

func getPages() []page {
	return []page{
		{
			name:    "overview",
			label:   "Overview",
			clicker: widget.OverviewNavClicker,
			handler: handlers.NewOverviewHandler(),
		},
		{
			name:    "history",
			label:   "History",
			clicker: widget.HistoryNavClicker,
			handler: handlers.NewOverviewHandler(),
		},
		{
			name:    "send",
			label:   "Send",
			clicker: widget.SendNavClicker,
			handler: handlers.NewOverviewHandler(),
		},
		{
			name:    "receive",
			label:   "Receive",
			clicker: widget.ReceiveNavClicker,
			handler: handlers.NewOverviewHandler(),
		},
		{
			name:    "staking",
			label:   "Staking",
			clicker: widget.StakingNavClicker,
			handler: handlers.NewOverviewHandler(),
		},
		{
			name:    "security",
			label:   "Security",
			clicker: widget.SecurityNavClicker,
			handler: handlers.NewOverviewHandler(),
		},
		{
			name:    "settings",
			label:   "Settings",
			clicker: widget.SettingsNavClicker,
			handler: handlers.NewOverviewHandler(),
		},
	}
}
