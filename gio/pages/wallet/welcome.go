package wallet

import (
	"gioui.org/layout"

	"github.com/raedahgroup/godcr/gio/helper"
	"github.com/raedahgroup/godcr/gio/widgets"
)

type (
	WelcomePage struct {
		multiWallet         *helper.MultiWallet
		createWalletButton  *widgets.Button
		restoreWalletButton *widgets.Button
	}
)

func NewWelcomePage(multiWallet *helper.MultiWallet) *WelcomePage {
	return &WelcomePage{
		multiWallet:         multiWallet,
		createWalletButton:  widgets.NewButton("Create a new wallet", widgets.AddIcon).SetBackgroundColor(helper.DecredLightBlueColor),
		restoreWalletButton: widgets.NewButton("Restore an existing wallet", widgets.ReturnIcon).SetBackgroundColor(helper.DecredGreenColor),
	}
}

func (w *WelcomePage) GetWidgets(ctx *layout.Context, changePageFunc func(page string)) []func() {
	widgets := []func(){
		// logo row 
		func(){
			helper.DrawLogo(ctx)
		},

		// welcome text first row
		func(){ 
			helper.Inset(ctx, 45, helper.StandaloneScreenPadding, 0, 0, func(){
				widgets.NewLabel("Welcome to", 6).Draw(ctx)
			})
		},

		// welcome text second row
		func(){
			helper.Inset(ctx, 10, helper.StandaloneScreenPadding, 0, 0, func(){
				widgets.NewLabel("Decred Desktop Wallet", 6).Draw(ctx)
			})
		},

		// create wallet button row
		func() {
			helper.Inset(ctx, 190, helper.StandaloneScreenPadding, 0,  helper.StandaloneScreenPadding, func(){
				ctx.Constraints.Height.Min = 50
				w.createWalletButton.Draw(ctx, func() {
					changePageFunc("createwallet")
				})
			})
		},

		// restore wallet button row 
		func() {
			helper.Inset(ctx, 10, helper.StandaloneScreenPadding, 0,  helper.StandaloneScreenPadding, func(){
				ctx.Constraints.Height.Min = 50
				w.restoreWalletButton.Draw(ctx, func() {
					changePageFunc("restorewallet")
				})
			})
		},
	}
	
	return widgets
}
