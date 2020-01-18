package wallet

import (
	"strconv"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/gio/helper"
	"github.com/raedahgroup/godcr/gio/widgets"
	"github.com/raedahgroup/godcr/gio/widgets/editor"
	"github.com/raedahgroup/godcr/gio/widgets/security"
)

type (
	restoreScreen struct {
		backButton    *widgets.Button
		restoreButton *widgets.Button
		inputs        []*editor.Input
	}

	RestoreWalletPage struct {
		multiWallet    *helper.MultiWallet
		changePageFunc func(string)
		currentScreen  string

		errStr               string
		pinAndPasswordWidget *security.PinAndPasswordWidget
		backToWalletsButton  *widgets.ClickableLabel

		isRestoring   bool
		restoreScreen *restoreScreen
	}
)

func NewRestoreWalletPage(multiWallet *helper.MultiWallet) *RestoreWalletPage {
	w := &RestoreWalletPage{
		multiWallet:         multiWallet,
		currentScreen:       "verifySeedScreen",
		backToWalletsButton: widgets.NewClickableLabel("Get Started").SetAlignment(widgets.AlignMiddle).SetSize(5).SetColor(helper.DecredLightBlueColor).SetWeight(text.Bold),
	}

	w.pinAndPasswordWidget = security.NewPinAndPasswordWidget(w.cancel, w.restore)

	// restore screen widgets
	w.restoreScreen = &restoreScreen{
		inputs:        make([]*editor.Input, 33),
		restoreButton: widgets.NewButton("Continue", nil),
		backButton:    widgets.NewButton("", widgets.NavigationArrowBackIcon).SetBackgroundColor(helper.BackgroundColor).SetColor(helper.BlackColor).MakeRound(),
	}
	for i := 0; i < 33; i++ {
		w.restoreScreen.inputs[i] = editor.NewInput("")
	}

	return w
}

func (w *RestoreWalletPage) cancel() {
	w.errStr = ""

	for i := range w.restoreScreen.inputs {
		w.restoreScreen.inputs[i].SetText("")
	}
	w.currentScreen = "verifySeedScreen"
}

func (w *RestoreWalletPage) restore() {
	w.isRestoring = true
	w.errStr = ""
	w.currentScreen = "verifySeedScreen"

	doneChan := make(chan bool)

	// do restoring here
	go func() {
		defer func() {
			doneChan <- true
		}()

		seed := ""
		for i := range w.restoreScreen.inputs {
			seed += w.restoreScreen.inputs[i].Text() + " "
		}

		password := w.pinAndPasswordWidget.Value()

		wallet, err := w.multiWallet.RestoreWallet("public", seed, password, 0)
		if err != nil {
			w.errStr = err.Error()
			return
		}

		w.errStr = wallet.UnlockWallet([]byte(password)).Error()
		if w.errStr == "" {
			w.multiWallet.RegisterWalletID(wallet.ID)
		}
	}()

	<-doneChan
	w.isRestoring = false
	if w.errStr == "" {
		w.currentScreen = "restoreSuccessScreen"
	}
}

func (w *RestoreWalletPage) GetWidgets(ctx *layout.Context, changePageFunc func(page string)) []func() {
	if w.changePageFunc == nil {
		w.changePageFunc = changePageFunc
	}

	if w.currentScreen == "verifySeedScreen" {
		return w.renderVerifySeedScreen(ctx, changePageFunc)
	} else if w.currentScreen == "passwordScreen" {
		w.renderPasswordScreen(ctx)
	} else if w.currentScreen == "restoreSuccessScreen" {
		w.restoreSuccessScreen(ctx, changePageFunc)
	}

	return nil
}

func (w *RestoreWalletPage) restoreSuccessScreen(ctx *layout.Context, changePageFunc func(string)) {
	ctx.Constraints.Width.Min = ctx.Constraints.Width.Max
	layout.Stack{}.Layout(ctx,
		layout.Expanded(func() {
			ctx.Constraints.Width.Min = ctx.Constraints.Width.Max
			layout.Align(layout.Center).Layout(ctx, func() {
				inset := layout.Inset{
					Top: unit.Dp(120),
				}
				inset.Layout(ctx, func() {
					ctx.Constraints.Width.Min = 50
					widgets.NewCheckbox().SetSize(80).MakeAsIcon().Draw(ctx)
				})
			})
		}),

		layout.Expanded(func() {
			inset := layout.Inset{
				Top: unit.Dp(220),
			}
			inset.Layout(ctx, func() {
				widgets.NewLabel("Your wallet is successfully").
					SetSize(6).
					SetWeight(text.Bold).
					SetAlignment(widgets.AlignMiddle).
					SetColor(helper.BlackColor).
					Draw(ctx)
			})

			inset = layout.Inset{
				Top: unit.Dp(245),
			}
			inset.Layout(ctx, func() {
				widgets.NewLabel("restored").
					SetSize(6).
					SetWeight(text.Bold).
					SetAlignment(widgets.AlignMiddle).
					SetColor(helper.BlackColor).
					Draw(ctx)
			})

			inset = layout.Inset{
				Top: unit.Dp(450),
			}
			inset.Layout(ctx, func() {
				w.backToWalletsButton.SetWidth(ctx.Constraints.Width.Max).Draw(ctx, func() {
					changePageFunc("overview")
				})
			})
		}),
	)
}

func (w *RestoreWalletPage) renderVerifySeedScreen(ctx *layout.Context, changePageFunc func(string)) []func() {
	return []func(){
		func(){
			drawHeader(ctx, func(){
				w.restoreScreen.backButton.Draw(ctx, func(){
					w.resetAndGotoPage("welcome")
				})
			}, func(){
				widgets.NewLabel("Restore from seed phrase").SetWeight(text.Bold).SetSize(6).Draw(ctx)
			})
		},
		func(){
			topInset := float32(10)
			if w.errStr != "" {
				inset := layout.Inset{
					Top: unit.Dp(topInset),
				}
				inset.Layout(ctx, func() {
					helper.PaintArea(ctx, helper.DangerColor, ctx.Constraints.Width.Max, 30)

					ctx.Constraints.Width.Min = ctx.Constraints.Width.Max
					widgets.NewLabel(w.errStr).
						SetSize(5).
						SetColor(helper.WhiteColor).
						SetAlignment(widgets.AlignMiddle).
						Draw(ctx)
				})
				topInset += 30
			}

			inset := layout.Inset{
				Top: unit.Dp(topInset),
			}
			inset.Layout(ctx, func() {
				(&layout.List{Axis: layout.Vertical}).Layout(ctx, 33, func(i int) {
					inset := layout.Inset{
						Top:   unit.Dp(10),
						Left:  unit.Dp(15),
						Right: unit.Dp(15),
					}
					inset.Layout(ctx, func() {
						layout.Flex{Axis: layout.Horizontal}.Layout(ctx,
							layout.Rigid(func() {
								inset := layout.Inset{
									Top:   unit.Dp(20),
									Right: unit.Dp(10),
								}
								inset.Layout(ctx, func() {
									widgets.NewLabel(strconv.Itoa(i + 1)).Draw(ctx)
								})
							}),
							layout.Rigid(func() {
								w.restoreScreen.inputs[i].Draw(ctx)
							}),
						)
					})
				})
			})
		},
		func(){
			drawFooter(ctx,10, func() {
				ctx.Constraints.Height.Min = 50
		
				bgCol := helper.GrayColor
				txt := "Continue"
				if w.hasEnteredAllSeedWords() {
					bgCol = helper.DecredLightBlueColor
				}
		
				if w.isRestoring {
					bgCol = helper.GrayColor
					txt = "Restoring..."
				}
				w.restoreScreen.
					restoreButton.
					SetText(txt).
					SetBackgroundColor(bgCol).
					Draw(ctx, func() {
						if w.hasEnteredAllSeedWords() && w.doVerify() {
							w.currentScreen = "passwordScreen"
						}
					})
			})
		},
	}
}

func (w *RestoreWalletPage) renderPasswordScreen(ctx *layout.Context) {
	inset := layout.Inset{
		Top:   unit.Dp(30),
		Left:  unit.Dp(helper.StandaloneScreenPadding),
		Right: unit.Dp(helper.StandaloneScreenPadding),
	}
	inset.Layout(ctx, func() {
		w.pinAndPasswordWidget.Render(ctx)
	})
}

func (w *RestoreWalletPage) hasEnteredAllSeedWords() bool {
	for i := range w.restoreScreen.inputs {
		if w.restoreScreen.inputs[i].Text() == "" {
			return false
		}
	}
	return true
}

func (w *RestoreWalletPage) doVerify() bool {
	w.errStr = ""

	seedString := ""
	for _,v := range w.restoreScreen.inputs {
		seedString += v.Text() + " "
	}
	
	if dcrlibwallet.VerifySeed(seedString) {
		w.errStr = ""
		return true
	}

	w.errStr = "Failed to restore. Please verify all words and try again"
	return false
}

func (w *RestoreWalletPage) resetAndGotoPage(page string) {
	w.errStr = ""
	w.changePageFunc(page)
}
