package pages

import (
	"context"

	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type Page struct {
	Title   string
	Handler pageHandler
}

type pageHandler interface {
	Render(ctx context.Context, wallet walletcore.Wallet, container *widgets.Box)
}

func GetPages() []*Page {
	return []*Page{
		{
			"Overview",
			&OverviewHandler{},
		},
		{
			"History",
			defaultPageNotImplemented,
		},
		{
			"Send",
			defaultPageNotImplemented,
		},
		{
			"Receive",
			defaultPageNotImplemented,
		},
		{
			"Staking",
			defaultPageNotImplemented,
		},
		{
			"Accounts",
			defaultPageNotImplemented,
		},
		{
			"Security",
			defaultPageNotImplemented,
		},
		{
			"Settings",
			defaultPageNotImplemented,
		},
	}
}

type pageNotImplemented struct{}

func (_ *pageNotImplemented) Render(ctx context.Context, wallet walletcore.Wallet, container *widgets.Box) {
	container.AddLabel("Coming Soon")
}

var defaultPageNotImplemented = &pageNotImplemented{}
