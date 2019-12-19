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

	layout.Stack{Alignment: layout.NW}.Layout(ctx, 
		layout.Expanded(func(){
			widgets.NewLabel("Create a Spending Password", 5).
				SetWeight(text.Bold).
				Draw(ctx)
		}),
		layout.Stacked(func(){
			inset := layout.Inset{
				Top: unit.Dp(25),
			}
			inset.Layout(ctx, func(){
				w.formTabContainer.Draw(ctx, w.passwordRenderFunc, w.pinRenderFunc)
			})
		}),
	)
}

func (w *CreateWalletPage) passwordRenderFunc(ctx *layout.Context) {
	var bothPasswordsMatch bool
	if (w.passwordTab.confirmPasswordInput.Len() > 0) && (w.passwordTab.confirmPasswordInput.Text() != w.passwordTab.passwordInput.Text()) {
		bothPasswordsMatch = false
	} else {
		bothPasswordsMatch = true
	}
	
	
	inset := layout.Inset{
		Left: unit.Dp(20),
		Right: unit.Dp(20),
	}
	inset.Layout(ctx, func(){
		// password section
		inset := layout.Inset{
			Top: unit.Dp(20),
		}
		inset.Layout(ctx, func(){
			go func(){
				w.passwordTab.passwordStrength = (dcrlibwallet.ShannonEntropy(w.passwordTab.passwordInput.Text()) / 4) * 100
			}()
			w.passwordTab.passwordInput.Draw(ctx)
		})

		// password strength section 
		inset = layout.Inset{
			Top: unit.Dp(85),
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

		// confirm password section 
		inset = layout.Inset{
			Top: unit.Dp(105),
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

		// error text section 
		inset = layout.Inset{
			Top: unit.Dp(165),
		}
		inset.Layout(ctx, func(){
			if !bothPasswordsMatch {
				widgets.NewLabel("Both passwords do not match").SetColor(helper.DangerColor).Draw(ctx)
			}
		})

		// buttons section
		inset = layout.Inset{
			Top: unit.Dp(185),
		}
		inset.Layout(ctx, func(){
			ctx.Constraints.Width.Min = ctx.Constraints.Width.Max
			layout.Stack{Alignment: layout.NE}.Layout(ctx, 
				layout.Stacked(func(){
					layout.Flex{Axis: layout.Horizontal}.Layout(ctx, 
						layout.Rigid(func(){
							inset := layout.Inset{
								Right: unit.Dp(10),
								Top: unit.Dp(10),
							}
							inset.Layout(ctx, func(){
								w.cancelLabel.Draw(ctx, func(){
									w.resetAndGotoPage("welcome")
								})
							})
						}),
						layout.Rigid(func(){
							bgCol := helper.GrayColor 
							if bothPasswordsMatch && w.passwordTab.confirmPasswordInput.Len() > 0 {
								bgCol = helper.DecredLightBlueColor
							}
							w.createButton.SetBackgroundColor(bgCol).Draw(ctx, func(){
								if bothPasswordsMatch && w.passwordTab.confirmPasswordInput.Len() > 0 {
									w.showSeedInformationPage()
								}
							})
						}),
					)
				}),
			)
		})
	})
}

func (w *CreateWalletPage) pinRenderFunc(ctx *layout.Context) {
	/**inset := layout.UniformInset(unit.Dp(20))
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
	})**/
}

func (w *CreateWalletPage) resetAndGotoPage(page string) {
	w.passwordTab.passwordInput.Clear()
	w.passwordTab.confirmPasswordInput.Clear()

	w.changePageFunc(page)
}