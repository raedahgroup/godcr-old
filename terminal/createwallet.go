package terminal

import (
	"github.com/gdamore/tcell"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/rivo/tview"
)

func createWallet(tviewApp *tview.Application, walletMiddleware app.WalletMiddleware) tview.Primitive {
	pages := tview.NewPages()
	setFocus := tviewApp.SetFocus

	layout := tview.NewFlex().SetDirection(tview.FlexRow)

	layout.AddItem(helpers.CenterAlignedTextView("Create Wallet"), 4, 1, false)

	// get seed and display to user
	seed, err := walletMiddleware.GenerateNewWalletSeed()
	if err != nil {
		return layout.AddItem(helpers.CenterAlignedTextView(err.Error()), 4, 1, false)
	}

	outputTextView := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	outputMessage := func(text string) {
		layout.RemoveItem(outputTextView)
		layout.AddItem(outputTextView.SetText(text), 0, 1, true)
	}

	var passphrase string
	passphraseField := tview.NewInputField().
		SetLabel("Wallet Passphrase:  ").
		SetMaskCharacter('*').
		SetFieldWidth(20).
		SetChangedFunc(func(text string) {
			passphrase = text
		})

	var confirmPassphrase string
	confirmPassphraseField := tview.NewInputField().
		SetLabel("Confirm Passphrase: ").
		SetMaskCharacter('*').
		SetFieldWidth(20).
		SetChangedFunc(func(text string) {
			confirmPassphrase = text
		})

	layout.AddItem(passphraseField, 2, 1, true)
	layout.AddItem(confirmPassphraseField, 2, 1, true)

	seedView := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(seed).SetRegions(true).SetWordWrap(true)
	seedView.SetBorder(true).SetTitle("Wallet Seed").SetTitleColor(helpers.SeedLabelColor)
	layout.AddItem(seedView, 7, 1, true)

	seedInfo := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(`Keep the seed in a safe place as you will NOT be able to restore your wallet without it. Please keep in mind that anyone who has access to the seed can also restore your wallet thereby giving them access to all your funds, so it is imperative that you keep it in a secure location.`).
		SetRegions(true).SetWordWrap(true)
	seedInfo.SetBorder(true).SetTitle("IMPORTANT NOTICE:").SetTitleColor(helpers.WarnColor)
	layout.AddItem(seedInfo, 7, 3, true)

	var hasStoredSeed bool
	checkbox := tview.NewCheckbox().SetLabel("I have stored wallet seed  ").SetChecked(false).SetChangedFunc(func(checked bool) {
		hasStoredSeed = checked
	})
	layout.AddItem(checkbox, 2, 1, true)

	createButton := tview.NewForm().
		AddButton("Create Wallet", func() {
			if len(passphrase) == 0 {
				outputMessage("Passphrase cannot empty")
				return
			}
			if passphrase != confirmPassphrase {
				outputMessage("passphrase does not match")
				return
			}
			if !hasStoredSeed {
				outputMessage("Please store seed in a safe location and check the box")
				return
			}

			err = walletMiddleware.CreateWallet(passphrase, seed)
			if err != nil {
				outputMessage(err.Error())
			}

			res := pageLoader(tviewApp, walletMiddleware)
			pages.AddAndSwitchToPage("sync", res, true)
			return
		})

	layout.AddItem(createButton, 0, 1, true)

	layout.SetFullScreen(true).SetBorder(true).SetBorderPadding(3, 1, 6, 4)

	passphraseField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			tviewApp.Stop()
			return nil
		}
		if event.Key() == tcell.KeyTAB {
			setFocus(confirmPassphraseField)
			return nil
		}

		return event
	})

	// listen to escape and left key press events on all form items and buttons
	confirmPassphraseField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			tviewApp.Stop()
			return nil
		}
		if event.Key() == tcell.KeyTAB {
			setFocus(checkbox)
			return nil
		}
		if event.Key() == tcell.KeyBacktab {
			setFocus(passphraseField)
			return nil
		}
		return event
	})

	// use different key press listener on first form item to watch for backtab press and restore focus to stake info
	checkbox.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			tviewApp.Stop()
			return nil
		}
		if event.Key() == tcell.KeyTAB {
			setFocus(createButton)
			return nil
		}
		if event.Key() == tcell.KeyBacktab {
			setFocus(confirmPassphraseField)
			return nil
		}
		return event
	})

	// use different key press listener on form button to watch for tab press and restore focus to stake info
	createButton.GetButton(0).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			tviewApp.Stop()
			return nil
		}
		if event.Key() == tcell.KeyTAB {
			setFocus(passphraseField)
			return nil
		}
		if event.Key() == tcell.KeyBacktab {
			setFocus(checkbox)
			return nil
		}
		return event
	})

	pages.AddPage("create", layout, true, true)

	return pages
}
