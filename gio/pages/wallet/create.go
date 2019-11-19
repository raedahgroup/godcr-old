package wallet 

import (
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/layout"
	"github.com/raedahgroup/dcrlibwallet"

	//"github.com/raedahgroup/godcr/gio/helper"
	"github.com/raedahgroup/godcr/gio/widgets"
	"github.com/raedahgroup/godcr/gio/widgets/editor"
)

type (
	passwordTab struct {
		passwordInput        *editor.Input 
		confirmPasswordInput *editor.Input
		cancelLabel          *widgets.ClickableLabel
		createButton         *widgets.Button
	}

	pinTab struct {

	}

	CreateWalletPage struct {
		multiWallet          *dcrlibwallet.MultiWallet 
		changePageFunc		  func(string)
		formTabContainer     *widgets.TabContainer
		passwordTab          *passwordTab
	}
)

func NewCreateWalletPage(multiWallet *dcrlibwallet.MultiWallet) *CreateWalletPage {
	page := &CreateWalletPage{
		multiWallet  : multiWallet,
		passwordTab  : &passwordTab{
			passwordInput       : editor.NewInput("Spending Password").SetMask("*"),
			confirmPasswordInput: editor.NewInput("Confirm Spending Password").SetMask("*"),
			cancelLabel         : widgets.NewClickableLabel("Cancel").SetSize(8),
			createButton        : widgets.NewButton("Create", nil),
		},
	}

	page.formTabContainer = widgets.NewTabContainer().AddTab("Password").AddTab("PIN")
	return page
}

func (w *CreateWalletPage) Render(ctx *layout.Context, refreshWindowFunc func(), changePageFunc func(page string)) {
	w.changePageFunc = changePageFunc
	
	stack := layout.Stack{}
	header := stack.Rigid(ctx, func(){
		widgets.NewLabel("Create a Spending Password", 5).
			SetWeight(text.Bold).
			Draw(ctx, widgets.AlignLeft)
	})
	form := stack.Expand(ctx, func(){
		inset := layout.Inset{
			Top: unit.Dp(25),
			Left: unit.Dp(0),
		}
		inset.Layout(ctx, func(){
			w.formTabContainer.Draw(ctx, w.passwordRenderFunc, w.pinRenderFunc)
		})
	})
	stack.Layout(ctx, header, form)
}

func (w *CreateWalletPage) passwordRenderFunc(ctx *layout.Context) {
	inset := layout.UniformInset(unit.Dp(20))
	inset.Layout(ctx, func(){
		stack := layout.Stack{}
		passwordSection := stack.Rigid(ctx, func(){
			inset := layout.Inset{}
			inset.Layout(ctx, func(){
				w.passwordTab.passwordInput.Draw(ctx)
			})
		})
		confirmPasswordSection := stack.Rigid(ctx, func(){
			inset := layout.Inset{
				Top:  unit.Dp(40),
			}
			inset.Layout(ctx, func(){
				w.passwordTab.confirmPasswordInput.Draw(ctx)
			})
		})
		buttonsSection := stack.Rigid(ctx, func(){
			inset := layout.Inset{
				Top: unit.Dp(80),
			}
			inset.Layout(ctx, func(){
			
				flex := layout.Flex{
					Axis: layout.Horizontal,
				}
				createButtonSection := flex.Rigid(ctx, func(){
					ctx.Constraints.Height.Max = 50
					w.passwordTab.createButton.Draw(ctx, widgets.AlignMiddle, func(){

					})
					w.passwordTab.cancelLabel.Draw(ctx, widgets.AlignLeft, func(){
						w.changePageFunc("welcome")
					})
				})
				
				cancelButtonSection := flex.Rigid(ctx, func(){
					w.passwordTab.cancelLabel.Draw(ctx, widgets.AlignLeft, func(){
						w.changePageFunc("welcome")
					})
				})
				flex.Layout(ctx, cancelButtonSection, createButtonSection)
			})
		})
		stack.Layout(ctx, passwordSection, confirmPasswordSection, buttonsSection)
	})
}

func (w *CreateWalletPage) pinRenderFunc(ctx *layout.Context) {
	inset := layout.UniformInset(unit.Dp(20))
	inset.Layout(ctx, func(){
		stack := layout.Stack{}
		passwordSection := stack.Rigid(ctx, func(){
			w.passwordTab.passwordInput.Draw(ctx)
		})
		confirmPasswordSection := stack.Rigid(ctx, func(){

		})
		stack.Layout(ctx, passwordSection, confirmPasswordSection)
	})
}