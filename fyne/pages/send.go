package pages

import (
	"image/color"

	"github.com/raedahgroup/dcrlibwallet"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type sendPageDynamicData struct {
	// houses all clickable box
	receivingAccount *widget.Box
	sendingAccount   *widget.Box

	errorLabel                           *widget.Label
	sendingSelectedAccountLabel          *widget.Label
	sendingSelectedAccountBalanceLabel   *widget.Label
	receivingSelectedAccountLabel        *widget.Label
	receivingSelectedAccountBalanceLabel *widget.Label
}

var sendPage sendPageDynamicData

func sendPageContent(dcrlw *dcrlibwallet.LibWallet) fyne.CanvasObject {
	// acct, _ := dcrlw.AccountNumber("default")
	// txauthor := dcrlw.NewUnsignedTx(int32(acct), 0)
	// txauthor.AddSendDestination("TsfDLrRkk9ciUuwfp2b8PawwnukYD7yAjGd", dcrlibwallet.AmountAtom(10), false)
	// fmt.Println(txauthor.EstimateMaxSendAmount())
	// hash, err := txauthor.Broadcast([]byte("admin"))
	// fmt.Println(hash, err)
	// fmt.Println(chainhash.NewHash(hash))

	icons, err := assets.GetIcons(assets.InfoIcon, assets.MoreIcon, assets.ReceiveAccountIcon, assets.CollapseIcon)
	if err != nil {
		return widget.NewLabelWithStyle(err.Error(), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	}

	// define base widget consisting of label, more icon and info button
	sendLabel := widget.NewLabelWithStyle("Send DCR", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true, Italic: true})
	clickabelInfoIcon := widgets.NewImageButton(icons[assets.InfoIcon], nil, func() {

	})
	clickabelMoreIcon := widgets.NewImageButton(icons[assets.MoreIcon], nil, func() {

	})

	baseWidgets := widget.NewHBox(sendLabel, layout.NewSpacer(), clickabelInfoIcon, clickabelMoreIcon)

	accounts, err := dcrlw.GetAccountsRaw(dcrlibwallet.DefaultRequiredConfirmations)
	if err != nil {
		return widget.NewLabelWithStyle("could not retrieve account, "+err.Error(), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	}

	sendPage.receivingSelectedAccountLabel = widget.NewLabel(accounts.Acc[0].Name)
	sendPage.receivingSelectedAccountBalanceLabel = widget.NewLabel(dcrutil.Amount(accounts.Acc[0].TotalBalance).String())
	sendPage.receivingAccount = widget.NewVBox()

	receivingAccountBox := widget.NewHBox(
		widgets.NewHSpacer(15),
		widget.NewIcon(icons[assets.ReceiveAccountIcon]),
		widgets.NewHSpacer(20),
		sendPage.receivingSelectedAccountLabel,
		widgets.NewHSpacer(30),
		sendPage.receivingSelectedAccountBalanceLabel,
		widgets.NewHSpacer(8),
		widget.NewIcon(icons[assets.CollapseIcon]),
	)

	receivingAccountSelectionPopup := widget.NewPopUp(sendPage.receivingAccount, fyne.CurrentApp().Driver().AllWindows()[0].Canvas())
	getAccountInBox(sendPage.receivingAccount, sendPage.receivingSelectedAccountLabel, sendPage.receivingSelectedAccountBalanceLabel,
		accounts, icons[assets.ReceiveAccountIcon], receivingAccountSelectionPopup)
	receivingAccountSelectionPopup.Hide()

	var receivingAccountClickableBox *widgets.ClickableBox
	receivingAccountClickableBox = widgets.NewClickableBox(receivingAccountBox, func() {
		receivingAccountSelectionPopup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(
			receivingAccountClickableBox).Add(fyne.NewPos(0, receivingAccountClickableBox.Size().Height)))
		receivingAccountSelectionPopup.Show()
	})

	// // make receivingAccountTab a clickable box thereby showing the popup

	receivingAccountGroup := widget.NewGroup("To", receivingAccountClickableBox)

	// sendingAccountsDropdown := widgets.NewClickableBox(receivingAccountTab, func() {

	// })

	// var accountDropdown *widgets.ClickableBox
	// accountDropdown = widgets.NewClickableBox(accountTab, func() {
	// 	accountSelectionPopup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(
	// 		accountDropdown).Add(fyne.NewPos(0, accountDropdown.Size().Height)))
	// 	accountSelectionPopup.Show()
	// })

	return widget.NewHBox(widgets.NewHSpacer(10), widget.NewVBox(baseWidgets, receivingAccountGroup))
}

func getAccountInBox(dropdownContent *widget.Box, selectedAccountLabel *widget.Label, selectedAccountBalanceLabel *widget.Label, accounts *dcrlibwallet.Accounts, receiveIcon fyne.Resource, popup *widget.PopUp) {
	for index, account := range accounts.Acc {
		if account.Name == "imported" {
			continue
		}

		spendableLabel := canvas.NewText("Spendable", color.White)
		spendableLabel.TextSize = 10

		accountName := account.Name
		accountNameLabel := widget.NewLabel(accountName)
		accountNameLabel.Alignment = fyne.TextAlignLeading
		accountNameBox := widget.NewVBox(
			accountNameLabel,
			widget.NewHBox(widgets.NewHSpacer(1), spendableLabel),
		)

		spendableAmountLabel := canvas.NewText(dcrutil.Amount(account.Balance.Spendable).String(), color.White)
		spendableAmountLabel.TextSize = 10
		spendableAmountLabel.Alignment = fyne.TextAlignTrailing

		amount := dcrutil.Amount(account.Balance.Total).String()
		accountBalance := amount
		accountBalanceLabel := widget.NewLabel(accountBalance)
		accountBalanceLabel.Alignment = fyne.TextAlignTrailing

		accountBalanceBox := widget.NewVBox(
			accountBalanceLabel,
			spendableAmountLabel,
		)

		checkmarkIcon := widget.NewIcon(theme.ConfirmIcon())
		var spacing fyne.CanvasObject
		if index != 0 {
			checkmarkIcon.Hide()
			spacing = widgets.NewHSpacer(35)
		} else {
			spacing = widgets.NewHSpacer(15)
		}

		accountsView := widget.NewHBox(
			widgets.NewHSpacer(15),
			widget.NewIcon(receiveIcon),
			widgets.NewHSpacer(20),
			accountNameBox,
			layout.NewSpacer(),
			widgets.NewHSpacer(30),
			accountBalanceBox,
			widgets.NewHSpacer(30),
			checkmarkIcon,
			spacing,
		)

		dropdownContent.Append(widgets.NewClickableBox(accountsView, func() {
			// hide checkmark icon of other accounts
			for _, children := range dropdownContent.Children {
				if box, ok := children.(*widgets.ClickableBox); !ok {
					continue
				} else {
					if len(box.Children) != 10 {
						continue
					}

					if icon, ok := box.Children[8].(*widget.Icon); !ok {
						continue
					} else {
						icon.Hide()
					}
					if spacing, ok := box.Children[9].(*fyne.Container); !ok {
						continue
					} else {
						spacing.Layout = layout.NewFixedGridLayout(fyne.NewSize(35, 0))
						canvas.Refresh(spacing)
					}
				}
				canvas.Refresh(children)
			}

			checkmarkIcon.Show()
			if spacing, ok := accountsView.Children[9].(*fyne.Container); !ok {
			} else {
				spacing.Layout = layout.NewFixedGridLayout(fyne.NewSize(15, 0))
				canvas.Refresh(spacing)
			}

			if accountbalanceBox, ok := accountsView.Children[6].(*widget.Box); ok {
				if len(accountbalanceBox.Children) == 2 {
					if balanceLabel, ok := accountbalanceBox.Children[0].(*widget.Label); ok {
						selectedAccountBalanceLabel.SetText(balanceLabel.Text)
					}
				}
			}
			selectedAccountLabel.SetText(accountName)
			popup.Hide()
		}))
	}
}

func updateAccountBoxContent(accountBox *widget.Box, account *dcrlibwallet.Accounts) {
	for index, boxContent := range accountBox.Children {
		spendableAmountLabel := canvas.NewText(dcrutil.Amount(account.Acc[index].Balance.Spendable).String(), color.White)
		spendableAmountLabel.TextSize = 10
		spendableAmountLabel.Alignment = fyne.TextAlignTrailing

		accountBalance := dcrutil.Amount(account.Acc[index].Balance.Total).String()
		accountBalanceLabel := widget.NewLabel(accountBalance)
		accountBalanceLabel.Alignment = fyne.TextAlignTrailing

		accountBalanceBox := widget.NewVBox(
			accountBalanceLabel,
			spendableAmountLabel,
		)

		accountBalance = dcrutil.Amount(account.Acc[index].Balance.Total).String()

		if content, ok := boxContent.(*widgets.ClickableBox); ok {
			content.Box.Children[6] = accountBalanceBox
			widget.Refresh(content.Box)
			widget.Refresh(content)
		}
	}
}
