package wallet 

import (
	"image/color"

	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/layout"
	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/gio/helper"
	"github.com/raedahgroup/godcr/gio/widgets"
	"github.com/raedahgroup/godcr/gio/widgets/editor"
)

type (
	passwordTab struct {
		passwordInput        		*editor.Input 
		confirmPasswordInput 		*editor.Input
		passwordStrengthProgressBar	*widgets.ProgressBar
		passwordStrength 			float64
	}

	pinTab struct {
		pinInput         *editor.Input 
		confirmPinInput  *editor.Input
		pinStrength 	 int
	}

	CreateWalletPage struct {
		multiWallet          *dcrlibwallet.MultiWallet 
		changePageFunc		  func(string)
		formTabContainer     *widgets.TabContainer
		passwordTab          *passwordTab
		pinTab               *pinTab
		cancelLabel          *widgets.ClickableLabel
		createButton         *widgets.Button
	}
)

func NewCreateWalletPage(multiWallet *dcrlibwallet.MultiWallet) *CreateWalletPage {
	passwordTab := &passwordTab{
		passwordInput               : editor.NewInput("Spending Password").SetMask("*"),
		confirmPasswordInput        : editor.NewInput("Confirm Spending Password").SetMask("*"),
		passwordStrength            : 0,
		passwordStrengthProgressBar : widgets.NewProgressBar().SetHeight(6),
	}

	pinTab := &pinTab{
		pinInput         : editor.NewInput("Pin").SetMask("*").Numeric(),
		confirmPinInput  : editor.NewInput("Confirm Pin").SetMask("*").Numeric(),
	}
	
	formTabContainer := widgets.NewTabContainer().AddTab("Password").AddTab("PIN")

	return &CreateWalletPage{
		multiWallet     :  multiWallet,
		formTabContainer:  formTabContainer,
		passwordTab     :  passwordTab,
		pinTab          :  pinTab,
		cancelLabel     :  widgets.NewClickableLabel("Cancel").SetSize(4).SetWeight(text.Bold).SetColor(helper.DecredLightBlueColor),
		createButton    :  widgets.NewButton("Create", nil),
	}
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
				go func(){
					w.passwordTab.passwordStrength = (dcrlibwallet.ShannonEntropy(w.passwordTab.passwordInput.Text()) / 4) * 100
				}()
				w.passwordTab.passwordInput.Draw(ctx)
			})
		})

		passwordStrengthSection := stack.Rigid(ctx, func(){
			inset := layout.Inset{
				Top: unit.Dp(55),
			}
			inset.Layout(ctx, func(){
				var col color.RGBA
				if w.passwordTab.passwordStrength > 70 {
					col = helper.DecredGreenColor
				} else {
					col = helper.DecredOrangeColor 
				}
				w.passwordTab.passwordStrengthProgressBar.SetProgressColor(col).Draw(ctx, &w.passwordTab.passwordStrength)
			})
		})

		bothPasswordsMatch := true 
		if (w.passwordTab.confirmPasswordInput.Len() > 0) && (w.passwordTab.confirmPasswordInput.Text() != w.passwordTab.passwordInput.Text()) {
			bothPasswordsMatch = false
		}

		confirmPasswordSection := stack.Rigid(ctx, func(){
			inset := layout.Inset{
				Top:  unit.Dp(76),
			}
			inset.Layout(ctx, func(){ 
				borderColor := helper.GrayColor 
				focusBorderColor := helper.DecredLightBlueColor 
				if !bothPasswordsMatch {
					borderColor = helper.DangerColor 
					focusBorderColor = helper.DangerColor
				}
				w.passwordTab.confirmPasswordInput.SetBorderColor(borderColor).SetFocusedBorderColor(focusBorderColor).Draw(ctx)
			})
		})

		errorTextSection := stack.Rigid(ctx, func(){
			inset := layout.Inset{
				Top: unit.Dp(-20),
			}
			inset.Layout(ctx, func(){
				if !bothPasswordsMatch {
					widgets.NewLabel("Both passwords do not match").
						SetColor(helper.DangerColor).
						Draw(ctx, widgets.AlignLeft)
				}
			})
		})

		buttonsSection := stack.Expand(ctx, func(){
			inset := layout.Inset{
				Top: unit.Dp(140),
			}
			ctx.Constraints.Width.Min = ctx.Constraints.Width.Max
			inset.Layout(ctx, func(){
				stack := layout.Stack{}

				cancelButtonSection := stack.Rigid(ctx, func(){
					inset := layout.Inset{
						Left: unit.Dp(310),
						Top : unit.Dp(10),
					}
					inset.Layout(ctx, func(){
						w.cancelLabel.Draw(ctx, widgets.AlignLeft, func(){
							w.resetAndGotoPage("welcome")
						})
					})
				})

				createButtonSection := stack.Rigid(ctx, func(){
					inset := layout.Inset{
						Left: unit.Dp(368),
					}
					inset.Layout(ctx, func(){
						bgCol := helper.GrayColor 
						if bothPasswordsMatch && w.passwordTab.confirmPasswordInput.Len() > 0 {
							bgCol = helper.DecredLightBlueColor
						}
						w.createButton.SetBackgroundColor(bgCol).Draw(ctx, widgets.AlignLeft, func(){
							if bothPasswordsMatch && w.passwordTab.confirmPasswordInput.Len() > 0 {
								w.showSeedInformationPage()
							}
						})
					})
				})

				stack.Layout(ctx, cancelButtonSection, errorTextSection, createButtonSection)
			})
		})
		stack.Layout(ctx, passwordSection, passwordStrengthSection, confirmPasswordSection, buttonsSection)
	})
}

func (w *CreateWalletPage) pinRenderFunc(ctx *layout.Context) {
	inset := layout.UniformInset(unit.Dp(20))
	inset.Layout(ctx, func(){
		stack := layout.Stack{}
		pinSection := stack.Rigid(ctx, func(){
			inset := layout.Inset{}
			inset.Layout(ctx, func(){
				w.pinTab.pinInput.Draw(ctx)
			})
		})

		confirmPinSection := stack.Rigid(ctx, func(){
			inset := layout.Inset{
				Top:  unit.Dp(50),
			}
			inset.Layout(ctx, func(){
				w.pinTab.confirmPinInput.Draw(ctx)
			})
		})

		buttonsSection := stack.Expand(ctx, func(){
			inset := layout.Inset{
				Top: unit.Dp(100),
			}
			ctx.Constraints.Width.Min = ctx.Constraints.Width.Max
			inset.Layout(ctx, func(){
				stack := layout.Stack{}

				cancelButtonSection := stack.Rigid(ctx, func(){
					inset := layout.Inset{
						Left: unit.Dp(240),
						Top : unit.Dp(10),
					}
					inset.Layout(ctx, func(){
						w.cancelLabel.Draw(ctx, widgets.AlignLeft, func(){
							w.resetAndGotoPage("welcome")
						})
					})
				})

				createButtonSection := stack.Rigid(ctx, func(){
					inset := layout.Inset{
						Left: unit.Dp(300),
					}
					inset.Layout(ctx, func(){
						w.createButton.Draw(ctx, widgets.AlignLeft, func(){
							
						})
					})
				})
				stack.Layout(ctx, cancelButtonSection, createButtonSection)
			})
		})
		stack.Layout(ctx, pinSection, confirmPinSection, buttonsSection)
	})
}

func (w *CreateWalletPage) resetAndGotoPage(page string) {
	w.passwordTab.passwordInput.Clear()
	w.passwordTab.confirmPasswordInput.Clear()

	w.changePageFunc(page)
}