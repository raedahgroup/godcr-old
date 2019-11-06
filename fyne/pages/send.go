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
	fromAccountSelect *widget.Box
	toAccountSelect   *widget.Box
	errorLabel        *widget.Label
}

var sendPage sendPageDynamicData

func sendPageContent(dcrlw *dcrlibwallet.LibWallet) fyne.CanvasObject {
	icons, err := assets.GetIcons(assets.InfoIcon, assets.MoreIcon)
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

	// fromAccountSelectionBox := widget.NewGroup("From")
	// ToAccountSelectionBox := widget.NewGroup("To")

	return widget.NewHBox(widgets.NewHSpacer(10), widget.NewVBox(baseWidgets, widget.NewLabel("Hello")))
}

func getAccountInBox(accountListWidget *widget.Box, accounts *dcrlibwallet.Accounts, receiveIcon fyne.Resource, popup *widget.PopUp) {
	for index, account := range accounts.Acc {
		if account.Name == "imported" {
			continue
		}

		spendableLabel := canvas.NewText("Spendable", color.White)
		spendableLabel.TextSize = 10
		spendableLabel.Alignment = fyne.TextAlignLeading

		spendableAmountLabel := canvas.NewText(dcrutil.Amount(account.Balance.Spendable).String(), color.White)
		spendableAmountLabel.TextSize = 10
		spendableAmountLabel.Alignment = fyne.TextAlignTrailing

		accountName := account.Name
		accountNameLabel := widget.NewLabel(accountName)
		accountNameLabel.Alignment = fyne.TextAlignLeading
		accountNameBox := widget.NewVBox(
			accountNameLabel,
			spendableLabel,
		)

		accountBalance := dcrutil.Amount(account.Balance.Total).String()
		accountBalanceLabel := widget.NewLabel(accountBalance)
		accountBalanceLabel.Alignment = fyne.TextAlignTrailing
		accountBalanceBox := widget.NewVBox(
			accountBalanceLabel,
			spendableAmountLabel,
		)

		checkmarkIcon := widget.NewIcon(theme.ConfirmIcon())
		if index != 0 {
			checkmarkIcon.Hide()
		}

		accountsView := widget.NewHBox(
			widgets.NewHSpacer(15),
			widget.NewIcon(receiveIcon),
			widgets.NewHSpacer(20),
			accountNameBox,
			widgets.NewHSpacer(20),
			accountBalanceBox,
			widgets.NewHSpacer(30),
			checkmarkIcon,
			widgets.NewHSpacer(15),
		)

		accountListWidget.Append(widgets.NewClickableBox(accountsView, func() {
			// hide checkmark icon of other accounts
			for _, children := range accountListWidget.Children {
				if box, ok := children.(*widgets.ClickableBox); !ok {
					continue
				} else {
					if len(box.Children) != 9 {
						continue
					}

					if icon, ok := box.Children[7].(*widget.Icon); !ok {
						continue
					} else {
						icon.Hide()
					}
				}
			}

			checkmarkIcon.Show()
			// what happens after selecting
			popup.Hide()
		}))
	}
}
