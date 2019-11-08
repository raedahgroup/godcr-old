package pages

import (
	"fmt"
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
	receivingAccountDropdownContent *widget.Box
	sendingAccountDropdownContent   *widget.Box

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

	// // make receivingAccountTab a clickable box thereby showing the popup
	receivingAccountClickableBox := createAccountDropdown(icons[assets.ReceiveAccountIcon], icons[assets.CollapseIcon], accounts, sendPage.receivingAccountDropdownContent, sendPage.receivingSelectedAccountLabel, sendPage.receivingSelectedAccountBalanceLabel)
	sendingAccountClickableBox := createAccountDropdown(icons[assets.ReceiveAccountIcon], icons[assets.CollapseIcon], accounts, sendPage.sendingAccountDropdownContent, sendPage.sendingSelectedAccountLabel, sendPage.sendingSelectedAccountBalanceLabel)

	receivingAccountGroup := widget.NewGroup("To", receivingAccountClickableBox)

	// sendingAccountsDropdown := widgets.NewClickableBox(receivingAccountTab, func() {

	// })

	// var accountDropdown *widgets.ClickableBox
	// accountDropdown = widgets.NewClickableBox(accountTab, func() {
	// 	accountSelectionPopup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(
	// 		accountDropdown).Add(fyne.NewPos(0, accountDropdown.Size().Height)))
	// 	accountSelectionPopup.Show()
	// })

	submit := widget.NewButton("Submut", func() {
		fmt.Println(sendPage.receivingSelectedAccountLabel.Text, sendPage.receivingSelectedAccountBalanceLabel.Text)
	})

	return widget.NewHBox(widgets.NewHSpacer(10), widget.NewVBox(baseWidgets, widget.NewVBox(receivingAccountGroup, sendingAccountClickableBox, submit)))
}

func createAccountDropdown(receiveAccountIcon, collapseIcon fyne.Resource, accounts *dcrlibwallet.Accounts, dropdownContent *widget.Box, selectedAccountLabel *widget.Label, selectedAccountBalanceLabel *widget.Label) (accountClickableBox *widgets.ClickableBox) {
	selectedAccountLabel = widget.NewLabel(accounts.Acc[0].Name)
	selectedAccountBalanceLabel = widget.NewLabel(dcrutil.Amount(accounts.Acc[0].TotalBalance).String())
	dropdownContent = widget.NewVBox()

	receivingAccountBox := widget.NewHBox(
		widgets.NewHSpacer(15),
		widget.NewIcon(receiveAccountIcon),
		widgets.NewHSpacer(20),
		selectedAccountLabel,
		widgets.NewHSpacer(30),
		selectedAccountBalanceLabel,
		widgets.NewHSpacer(8),
		widget.NewIcon(collapseIcon),
	)

	receivingAccountSelectionPopup := widget.NewPopUp(dropdownContent, fyne.CurrentApp().Driver().AllWindows()[0].Canvas())
	getAccountInBox(dropdownContent, selectedAccountLabel, selectedAccountBalanceLabel,
		accounts, receiveAccountIcon, receivingAccountSelectionPopup)
	receivingAccountSelectionPopup.Hide()

	accountClickableBox = widgets.NewClickableBox(receivingAccountBox, func() {
		receivingAccountSelectionPopup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(
			accountClickableBox).Add(fyne.NewPos(0, accountClickableBox.Size().Height)))
		receivingAccountSelectionPopup.Show()
	})

	return
}

func getAccountInBox(dropdownContent *widget.Box, selectedAccountLabel, selectedAccountBalanceLabel *widget.Label, accounts *dcrlibwallet.Accounts, receiveIcon fyne.Resource, popup *widget.PopUp) {
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

func updateAccountDropdownContent(accountBox *widget.Box, account *dcrlibwallet.Accounts) {
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
			fmt.Println("worksss")
			content.Box.Children[6] = accountBalanceBox
			widget.Refresh(content.Box)
			widget.Refresh(content)
		}
	}
}
