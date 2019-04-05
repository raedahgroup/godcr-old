package pages

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/godcr/app/walletcore"
)

type Page struct {
	Title    string
	Content func(wallet walletcore.Wallet) fyne.CanvasObject
}

func pageNotImplementedContent(_ walletcore.Wallet) fyne.CanvasObject {
	return widget.NewLabelWithStyle("Page is not implemented yet.", fyne.TextAlignLeading, fyne.TextStyle{Italic:true})
}

func NavPages() []*Page {
	return []*Page{
		{
			"Overview",
			overviewPageContent,
		},
		{
			"History",
			pageNotImplementedContent,
		},
		{
			"Send",
			pageNotImplementedContent,
		},
		{
			"Receive",
			pageNotImplementedContent,
		},
		{
			"Staking",
			pageNotImplementedContent,
		},
		{
			"Accounts",
			pageNotImplementedContent,
		},
		{
			"Security",
			pageNotImplementedContent,
		},
		{
			"Settings",
			pageNotImplementedContent,
		},
	}
}
