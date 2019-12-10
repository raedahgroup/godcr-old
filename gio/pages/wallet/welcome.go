package wallet

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/gio/helper"
	"github.com/raedahgroup/godcr/gio/widgets"
)

type (
	WelcomePage struct {
		multiWallet *dcrlibwallet.MultiWallet
		createWalletButton *widgets.Button 
		restoreWalletButton *widgets.Button
	}
)

func NewWelcomePage(multiWallet *dcrlibwallet.MultiWallet) *WelcomePage {
	return &WelcomePage{
		multiWallet        :  multiWallet,
		createWalletButton :  widgets.NewButton("Create a new wallet", widgets.AddIcon).SetBackgroundColor(helper.DecredLightBlueColor),
		restoreWalletButton:  widgets.NewButton("Restore an existing wallet", widgets.ReturnIcon).SetBackgroundColor(helper.DecredGreenColor),
	}
}

func (w *WelcomePage) Render(ctx *layout.Context, refreshWindowFunc func(), changePageFunc func(page string)) {
	stack := layout.Stack{}

	welcomeTextSection := stack.Rigid(ctx, func() {
		inset := layout.Inset{
			Top: unit.Dp(0),
		}
		inset.Layout(ctx, func() {
			widgets.NewLabel("Welcome to", 6).Draw(ctx, widgets.AlignLeft)
		})

		inset = layout.Inset{
			Top: unit.Dp(29),
		}
		inset.Layout(ctx, func() {
			widgets.NewLabel("Decred Desktop Wallet", 6).Draw(ctx, widgets.AlignLeft)
		})
	})

	createButtonSection := stack.Rigid(ctx, func() {
		inset := layout.Inset{
			Top: unit.Dp(280),
		}
		inset.Layout(ctx, func() {
			w.createWalletButton.Draw(ctx, widgets.AlignLeft, func(){
				changePageFunc("createwallet")
			})
		})
	})

	restoreButtonSection := stack.Rigid(ctx, func() {
		inset := layout.Inset{
			Top: unit.Dp(330),
		}
		inset.Layout(ctx, func() {
			w.restoreWalletButton.Draw(ctx, widgets.AlignLeft, func(){
				changePageFunc("restorewallet")
			})
		})
	})

	stack.Layout(ctx, welcomeTextSection, createButtonSection, restoreButtonSection)
}
