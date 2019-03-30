package pages

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

func CreateWalletPage(tviewApp *tview.Application, walletMiddleware app.WalletMiddleware) tview.Primitive {
	walletPage := tview.NewFlex().SetDirection(tview.FlexRow)
	walletPage.SetBorderPadding(1, 1, 2, 2).SetBackgroundColor(tcell.ColorBlack)

	// page title and hint
	walletPage.AddItem(primitives.NewCenterAlignedTextView("First Time? Create Wallet"), 1, 0, false)
	hintText := primitives.WordWrappedTextView("(Use Tab to move between fields, Arrow Keys to scroll each field, Esc to exit)")
	hintText.SetTextColor(tcell.ColorGray)
	walletPage.AddItem(hintText, 3, 0, false)

	// attempt to get seed and display any error to user
	seed, err := walletMiddleware.GenerateNewWalletSeed()
	if err != nil {
		return walletPage.AddItem(primitives.NewCenterAlignedTextView(err.Error()), 4, 1, false)
	}

	walletSeedTextView := primitives.WordWrappedTextView(seed)
	walletSeedTextView.SetBorder(true).
		SetTitle("Wallet Seed").
		SetTitleColor(helpers.SeedLabelColor)
	walletPage.AddItem(walletSeedTextView, 6, 0, true)

	storeSeedWarningTextView := primitives.WordWrappedTextView(walletcore.StoreSeedWarningText)
	storeSeedWarningTextView.SetBorder(true).
		SetTitle("IMPORTANT NOTICE").
		SetTitleColor(helpers.WarnColor)
	walletPage.AddItem(storeSeedWarningTextView, 6, 0, false)

	// add single line space before button
	walletPage.AddItem(nil, 1, 0, false)

	// text view to display error messages below the create wallet button
	var errorMessageTextView *tview.TextView

	createWalletButton := tview.NewButton("I've stored the seed. Create Wallet").SetSelectedFunc(func() {
		x, y := walletSeedTextView.GetScrollOffset()
		errorMessageTextView.SetText(fmt.Sprintf("%d, %d", x, y))
	})
	walletPage.AddItem(createWalletButton, 1, 0, false)

	// add single line space before error message text view
	walletPage.AddItem(nil, 1, 0, false)

	walletSeedTextView.ScrollToEnd()
	x, y := walletSeedTextView.GetScrollOffset()
	errorMessageTextView = primitives.WordWrappedTextView(fmt.Sprintf("%d, %d", x, y))
	walletPage.AddItem(errorMessageTextView, 3, 0, false)

	//passphraseField := tview.NewInputField().
	//	SetLabel("Wallet Passphrase:  ").
	//	SetMaskCharacter('*').
	//	SetFieldWidth(20)
	//walletPage.AddItem(passphraseField, 1, 0, true)
	//
	//// add single line space between views
	//walletPage.AddItem(nil, 1, 0, false)
	//
	//confirmPassphraseField := tview.NewInputField().
	//	SetLabel("Confirm Passphrase: ").
	//	SetMaskCharacter('*').
	//	SetFieldWidth(20)
	//walletPage.AddItem(confirmPassphraseField, 1, 0, false)
	//
	//// add single line space between views
	//walletPage.AddItem(nil, 1, 0, false)

	//createWalletFunction := func() {
		//passphrase := passphraseField.GetText()
		//if len(passphrase) == 0 {
		//	errorMessageTextView.SetText("Passphrase cannot empty")
		//	return
		//}
		//
		//confirmPassphrase := confirmPassphraseField.GetText()
		//if passphrase != confirmPassphrase {
		//	errorMessageTextView.SetText("Passphrase does not match")
		//	return
		//}
		//
		//if !storeSeedCheckbox.IsChecked() {
		//	errorMessageTextView.SetText("Please store seed in a safe location and check the box")
		//	return
		//}
		//
		//err = walletMiddleware.CreateWallet(passphrase, seed)
		//if err != nil {
		//	errorMessageTextView.SetText(err.Error())
		//}
		//
		//// wallet created, go to sync page and begin sync
		//tviewApp.SetRoot(SyncPage(tviewApp, walletMiddleware), true)
	//}


	//createWalletButton := tview.NewButton("Create Wallet").SetSelectedFunc(createWalletFunction)
	//
	//// add button, single line space and then error message text view
	//walletPage.AddItem(createWalletButton, 1, 0, false)
	//walletPage.AddItem(nil, 1, 0, false)
	//walletPage.AddItem(errorMessageTextView, 1, 0, false)

	allViews := []*tview.Box{
		walletSeedTextView.Box,
		storeSeedWarningTextView.Box,
		createWalletButton.Box,
		errorMessageTextView.Box,
	}

	// exit page on escape key press, switch focus to next view on tab press and focus previous view on backtab press
	makeInputCaptureListener := func(viewIndex int) func(event *tcell.EventKey) (nextEvent *tcell.EventKey) {
		return func(event *tcell.EventKey) (nextEvent *tcell.EventKey) {
			switch event.Key() {
			case tcell.KeyEscape:
				tviewApp.Stop()

			case tcell.KeyTAB:
				nextViewIndex := viewIndex + 1
				if nextViewIndex >= len(allViews) {
					nextViewIndex = 0 // last view reached, go back to top
				}

				nextView := allViews[nextViewIndex]
				tviewApp.SetFocus(nextView)

			case tcell.KeyBacktab:
				previousViewIndex := viewIndex - 1
				if previousViewIndex < 0 {
					previousViewIndex = len(allViews) - 1 // first view reached, go to bottom
				}

				previousView := allViews[previousViewIndex]
				tviewApp.SetFocus(previousView)

			default:
				nextEvent = event
			}
			return
		}
	}

	for i, view := range allViews {
		view.SetInputCapture(makeInputCaptureListener(i))
	}

	walletPage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			tviewApp.Stop()
			return nil
		}
		return event
	})
	tviewApp.SetFocus(walletPage)

	return walletPage
}
