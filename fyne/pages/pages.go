package pages

import (
	"context"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/godcr/app/walletcore"
)

type Page struct {
	Title string
	PageLoader interface{}
}

type SimplePageLoader interface {
	Load(updatePageOnMainWindow func(object fyne.CanvasObject))
}

type WalletPageLoader interface {
	Load(ctx context.Context, wallet walletcore.Wallet, updatePageOnMainWindow func(object fyne.CanvasObject))
}

type pageNotImplemented struct {}

func (_ *pageNotImplemented) Load(updatePageOnMainWindow func(object fyne.CanvasObject)) {
	notice := widget.NewLabelWithStyle("Page is not implemented yet.", fyne.TextAlignLeading, fyne.TextStyle{Italic:true})
	updatePageOnMainWindow(notice)
}

var defaultPageNotImplemented = &pageNotImplemented{}

func NavPages() []*Page {
	return []*Page{
		{
			"Overview",
			&overviewPageLoader{},
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
