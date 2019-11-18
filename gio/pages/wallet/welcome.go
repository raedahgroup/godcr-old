package wallet

import (
	"github.com/raedahgroup/dcrlibwallet"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/text"

	"github.com/raedahgroup/godcr/gio/widgets"
	"github.com/raedahgroup/godcr/gio/helper"
)

type (
	WelcomePage struct {
		multiWallet *dcrlibwallet.MultiWallet
		theme *helper.Theme

		createWalletButton *widget.Button 
		restoreWalletButton *widget.Button
	}
)

func NewWelcomePage(multiWallet *dcrlibwallet.MultiWallet, theme *helper.Theme) *WelcomePage {
	return &WelcomePage{
		multiWallet: multiWallet,
		theme: theme,
		createWalletButton: new(widget.Button),
		restoreWalletButton: new(widget.Button),
	}
}

func (w *WelcomePage) Render(ctx *layout.Context, refreshWindowFunc func(), changePageFunc func(page string)) {
	stack := layout.Stack{}

	welcomeTextSection := stack.Rigid(ctx, func(){
		inset := layout.Inset{
			Top: unit.Dp(0),
		}
		inset.Layout(ctx, func(){
			widgets.BoldText("Welcome to Decred Wallet", w.theme, ctx)
		})
	})

	createButtonSection := stack.Rigid(ctx, func(){
		inset := layout.Inset{
			Top: unit.Dp(220),
		}
		c := ctx.Constraints
		inset.Layout(ctx, func(){
			ctx.Constraints.Width.Min = ctx.Constraints.Width.Max

			btn := w.theme.Button("Create a new wallet")
			btn.Font = text.Font{
				Size: unit.Dp(14),
			}
			btn.Background = helper.DecredLightBlueColor

			for w.createWalletButton.Clicked(ctx) {
				changePageFunc("createwallet")
			}

			btn.Layout(ctx, w.createWalletButton)
		})
		ctx.Constraints = c
	})

	restoreButtonSection := stack.Rigid(ctx, func(){
		inset := layout.Inset{
			Top: unit.Dp(273),
		}
		c := ctx.Constraints
		inset.Layout(ctx, func(){
			ctx.Constraints.Width.Min = ctx.Constraints.Width.Max
			
			btn := w.theme.Button("Restore an existing wallet")
			btn.Font = text.Font{
				Size: unit.Dp(14),
			}
			btn.Background = helper.DecredGreenColor

			for w.restoreWalletButton.Clicked(ctx) {
				changePageFunc("restorewallet")
			}

			btn.Layout(ctx, w.restoreWalletButton)
		})
		ctx.Constraints = c
	})


	stack.Layout(ctx, welcomeTextSection, createButtonSection, restoreButtonSection)
}