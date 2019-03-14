package terminal

import (
	"context"

	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/rivo/tview"
)

func pageLoader(ctx context.Context, tviewApp *tview.Application, walletMiddleware app.WalletMiddleware) tview.Primitive {

	walletExists, err := openWalletIfExist(ctx, walletMiddleware)
	if err != nil {
		return helpers.CenterAlignedTextView(err.Error())
	}

	pages := tview.NewPages()

	if walletExists {
		pages.AddPage("main", terminalLayout(ctx, tviewApp, walletMiddleware), true, true)
	} else {

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

		form := tview.NewForm()
		var password string
		form.AddPasswordField("Passphrase", "", 20, '*', func(text string) {
			password = text
		})

		var confPassword string
		form.AddPasswordField("Confirm Passphrase", "", 20, '*', func(text string) {
			confPassword = text
		})

		var hasStoredSeed bool
		form.AddCheckbox("I have stored wallet seed", false, func(checked bool) {
			hasStoredSeed = checked
		})

		form.AddButton("Create", func() {
			if len(password) == 0 {
				outputMessage("Passphrase cannot empty")
				return
			}
			if password != confPassword {
				outputMessage("passphrase does not match")
				return
			}
			if hasStoredSeed == false {
				outputMessage("please Store seed in a szfe location and check the box")
				return
			}

			err := CreateWallet(tviewApp, seed, password, walletMiddleware)
			if err != nil {
				outputMessage(err.Error())
				return
			}

			pages.AddAndSwitchToPage("main", terminalLayout(ctx, tviewApp, walletMiddleware), true)
		})

		form.AddButton("Quit", func() {
			tviewApp.Stop()
		})

		layout.AddItem(form, 10, 1, true)

		seedView := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(seed).SetRegions(true).SetWordWrap(true)
		seedView.SetBorder(true).SetTitle("Wallet Seed").SetTitleColor(helpers.PrimaryColor)
		layout.AddItem(seedView, 7, 1, false)

		seedInfo := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(`Keep the seed in a safe place as you will NOT be able to restore your wallet without it. Please keep in mind that anyone who has access to the seed can also restore your wallet thereby giving them access to all your funds, so it is imperative that you keep it in a secure location.`).
		SetRegions(true).SetWordWrap(true)
		seedInfo.SetBorder(true).SetTitle("IMPORTANT NOTICE:").SetTitleColor(helpers.WarnColor)
		layout.AddItem(seedInfo, 7, 1, false)

		layout.SetFullScreen(true).SetBorderPadding(3, 1, 6, 4)

		layout.SetFullScreen(true)

		pages.AddPage("create", layout, true, true)
	}

	return pages
}