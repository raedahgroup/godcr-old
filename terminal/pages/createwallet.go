package pages

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

func CreateWalletPage(tviewApp *tview.Application, walletMiddleware app.WalletMiddleware) tview.Primitive {
	createWalletPage := tview.NewFlex().SetDirection(tview.FlexRow)
	createWalletPage.SetBorderPadding(1, 1, 2, 2).SetBackgroundColor(tcell.ColorBlack)

	// page title and hint
	pageTitle := primitives.NewCenterAlignedTextView("First Time? Create Wallet")
	createWalletPage.AddItem(pageTitle, 1, 0, false)

	// attempt to get seed and display any error to user
	seed, err := walletMiddleware.GenerateNewWalletSeed()
	if err != nil {
		errTextView := primitives.NewCenterAlignedTextView(err.Error())
		createWalletPage.AddItem(errTextView, 0, 1, false)

		createWalletPage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyEsc {
				tviewApp.Stop()
				return nil
			}
			return event
		})

		tviewApp.SetFocus(createWalletPage)
		return createWalletPage
	}

	createWalletForm := primitives.NewForm(false)
	createWalletPage.AddItem(createWalletForm, 0, 1, true)

	passphraseField := tview.NewInputField().
		SetLabel("Wallet Passphrase:  ").
		SetMaskCharacter('*').
		SetFieldWidth(20)
	createWalletForm.AddFormItem(passphraseField)

	confirmPassphraseField := tview.NewInputField().
		SetLabel("Confirm Passphrase: ").
		SetMaskCharacter('*').
		SetFieldWidth(20)
	createWalletForm.AddFormItem(confirmPassphraseField)

	walletSeedTextView := primitives.WordWrappedTextView(seed)
	walletSeedTextView.SetBorder(true).
		SetTitle("Wallet Seed").
		SetTitleColor(helpers.DecredLightBlueColor)
	createWalletForm.AddFormItem(primitives.NewTextViewFormItem(walletSeedTextView, 20, 1, true))

	storeSeedWarningTextView := primitives.WordWrappedTextView(walletcore.StoreSeedWarningText)
	storeSeedWarningTextView.SetBorder(true).
		SetTitle("IMPORTANT NOTICE").
		SetTitleColor(helpers.DecredOrangeColor)
	createWalletForm.AddFormItem(primitives.NewTextViewFormItem(storeSeedWarningTextView, 20, 1, true))

	storeSeedCheckbox := tview.NewCheckbox().SetLabel("I've stored the seed securely")
	createWalletForm.AddFormItem(storeSeedCheckbox)

	var isShowingMessage bool
	clearMessages := func() {
		if isShowingMessage {
			createWalletForm.RemoveFormItem(createWalletForm.GetFormItemsCount() - 1)
			isShowingMessage = false
			tviewApp.ForceDraw()
		}
	}

	showMessage := func(message string, isError bool) {
		var messageColor tcell.Color
		if isError {
			messageColor = helpers.DecredOrangeColor
			message = fmt.Sprintf("Error: %s", message)
		} else {
			messageColor = helpers.DecredGreenColor
			message = fmt.Sprintf("Success: %s", message)
		}

		messageTextView := primitives.NewCenterAlignedTextView(message)
		messageTextView.SetTextColor(messageColor)

		messageTextViewAsFormItem := primitives.NewTextViewFormItem(messageTextView, 20, 1, true)
		createWalletForm.AddFormItem(messageTextViewAsFormItem)

		isShowingMessage = true
	}

	var isCreatingWallet bool
	createWalletForm.AddButton("Create Wallet", func() {
		clearMessages()

		passphrase := passphraseField.GetText()
		if len(passphrase) == 0 {
			showMessage("Passphrase cannot empty", true)
			return
		}

		confirmPassphrase := confirmPassphraseField.GetText()
		if passphrase != confirmPassphrase {
			showMessage("Passphrase does not match", true)
			return
		}

		if !storeSeedCheckbox.IsChecked() {
			showMessage("Please store seed in a safe location and check the box", true)
			return
		}

		// create wallet in subroutine to prevent blocking the UI
		isCreatingWallet = true
		createWalletForm.GetButton(0).SetLabel("Creating...")
		go func() {
			err = walletMiddleware.CreateWallet(passphrase, seed)
			if err != nil {
				showMessage(err.Error(), true)
				isCreatingWallet = false
				return
			}

			// wallet created, display success message
			tviewApp.QueueUpdateDraw(func() {
				showMessage(fmt.Sprintf(`%s wallet created successfully`, strings.Title(walletMiddleware.NetType())), false)
				createWalletForm.GetButton(0).SetLabel("Done!")
			})

			// wait briefly then go to sync page and begin sync
			time.Sleep(1 * time.Second)

			tviewApp.QueueUpdateDraw(func() {
				LaunchSyncPage(tviewApp, walletMiddleware)
			})
		}()
	})

	createWalletForm.SetCancelFunc(func() {
		if !isCreatingWallet {
			tviewApp.Stop()
		}
	})

	tviewApp.SetFocus(createWalletPage)

	hintText := primitives.WordWrappedTextView("(Use TAB and Shift+TAB to move between fields and ESC to cancel)")
	hintText.SetTextColor(tcell.ColorGray)
	createWalletPage.AddItem(hintText, 2, 0, false)

	return createWalletPage
}
