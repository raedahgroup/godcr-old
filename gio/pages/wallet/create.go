package wallet 

import (
	"gioui.org/text"
	"gioui.org/layout"
	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/gio/helper"
	"github.com/raedahgroup/godcr/gio/widgets"
)

type (
	CreateWalletPage struct {
		multiWallet *dcrlibwallet.MultiWallet 
	}
)

func NewCreateWalletPage(multiWallet *dcrlibwallet.MultiWallet) *CreateWalletPage {
	return &CreateWalletPage{
		multiWallet: multiWallet,
	}
}

func (w *CreateWalletPage) Render(ctx *layout.Context, refreshWindowFunc func(), changePageFunc func(page string)) {
	stack := layout.Stack{}
	child := stack.Expand(ctx, func(){
		widgets.NewLabel("Page not yet implemented", 4).
			SetWeight(text.Bold).SetStyle(text.Italic).
			SetColor(helper.GrayColor).
			Draw(ctx, widgets.AlignLeft)
	})
	stack.Layout(ctx, child)
}