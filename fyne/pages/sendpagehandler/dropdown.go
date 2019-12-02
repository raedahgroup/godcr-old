package sendpagehandler

import (
	"image/color"
	"log"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/fyne/layouts"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func CreateAccountSelector(initFunction func(), accountLabel string, receiveAccountIcon, collapseIcon fyne.Resource,
	multiWallet *dcrlibwallet.MultiWallet, walletIDs []int, sendingSelectedWalletID *int,
	accountBoxes []*widget.Box, selectedAccountLabel *widget.Label,
	selectedAccountBalanceLabel *widget.Label, selectedWalletLabel *canvas.Text, contents *widget.Box) (accountClickableBox *widgets.ClickableBox) {

	dropdownContent := widget.NewVBox()

	selectAccountBox := widget.NewHBox(
		widgets.NewHSpacer(15),
		widget.NewVBox(widgets.NewVSpacer(10), widget.NewIcon(receiveAccountIcon)),
		widgets.NewHSpacer(20),
		fyne.NewContainerWithLayout(layouts.NewVBox(12), selectedAccountLabel, selectedWalletLabel),
		widgets.NewHSpacer(30),
		widget.NewVBox(widgets.NewVSpacer(4), selectedAccountBalanceLabel),
		widgets.NewHSpacer(8),
		widget.NewVBox(widgets.NewVSpacer(6), widget.NewIcon(collapseIcon)),
	)

	var accountSelectionPopup *widget.PopUp
	accountSelectionPopupHeader := widget.NewVBox(
		widgets.NewVSpacer(5),
		widget.NewHBox(
			widgets.NewHSpacer(16),
			widgets.NewImageButton(theme.CancelIcon(), nil, func() { accountSelectionPopup.Hide() }),
			widgets.NewHSpacer(16),
			widget.NewLabelWithStyle(accountLabel, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			layout.NewSpacer(),
		),
		widgets.NewVSpacer(5),
		canvas.NewLine(color.Black),
	)

	popupContent := widget.NewVBox(accountSelectionPopupHeader)
	accountSelectionPopup = widget.NewPopUp(popupContent, fyne.CurrentApp().Driver().AllWindows()[0].Canvas())
	accountSelectionPopup.Hide()

	// we cant access the children of group widget, proposed hack is to
	// create a vertical box array where all accounts would be placed,
	// then when we want to hide checkmarks we call all children of accountbox and hide checkmark icon except selected
	for walletIndex, walletID := range walletIDs {
		getAllWalletAccountsInBox(initFunction, dropdownContent, selectedAccountLabel, selectedAccountBalanceLabel, selectedWalletLabel,
			multiWallet.WalletWithID(walletID), walletIndex, walletID, sendingSelectedWalletID, accountBoxes, receiveAccountIcon, accountSelectionPopup)
	}

	dropdownContentWithScroller := fyne.NewContainerWithLayout(
		layout.NewFixedGridLayout(fyne.NewSize(dropdownContent.MinSize().Width+5, fyne.Min(dropdownContent.MinSize().Height, 100))),
		widget.NewScrollContainer(dropdownContent))
	popupContent.Append(dropdownContentWithScroller)

	accountClickableBox = widgets.NewClickableBox(selectAccountBox, func() {
		accountSelectionPopup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(
			accountClickableBox).Add(fyne.NewPos(0, accountClickableBox.Size().Height)))

		accountSelectionPopup.Show()
		accountSelectionPopup.Resize(dropdownContentWithScroller.Size().Add(fyne.NewSize(10, accountSelectionPopupHeader.MinSize().Height)))
		contents.Refresh()
	})

	return
}

func getAllWalletAccountsInBox(initFunction func(), dropdownContent *widget.Box, selectedAccountLabel,
	selectedAccountBalanceLabel *widget.Label, selectedWalletLabel *canvas.Text, wallet *dcrlibwallet.Wallet, walletIndex, walletID int,
	sendingSelectedWalletID *int, accountsBoxes []*widget.Box, receiveIcon fyne.Resource, popup *widget.PopUp) {

	accounts, err := wallet.GetAccountsRaw(dcrlibwallet.DefaultRequiredConfirmations)
	if err != nil {
		return
	}

	var groupedWalletsAccounts = widget.NewGroup(wallet.Name)
	// we cant access children of a group so a box is used
	accountsBox := widget.NewVBox()

	for index, account := range accounts.Acc {
		if account.Name == "imported" {
			continue
		}

		spendableLabel := canvas.NewText("Spendable", color.Black)
		spendableLabel.TextSize = 10

		accountName := account.Name
		accountNameLabel := widget.NewLabel(accountName)
		accountNameLabel.Alignment = fyne.TextAlignLeading
		accountNameBox := widget.NewVBox(
			accountNameLabel,
			widget.NewHBox(widgets.NewHSpacer(1), spendableLabel),
		)

		spendableAmountLabel := canvas.NewText(dcrutil.Amount(account.Balance.Spendable).String(), color.Black)
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
		if index != 0 || walletID != *sendingSelectedWalletID {
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

		accountsBox.Append(widgets.NewClickableBox(accountsView, func() {
			*sendingSelectedWalletID = walletID
			for _, boxes := range accountsBoxes {
				for _, objectsChild := range boxes.Children {
					if box, ok := objectsChild.(*widgets.ClickableBox); !ok {
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

					canvas.Refresh(objectsChild)
				}
			}

			checkmarkIcon.Show()
			if spacing, ok := accountsView.Children[9].(*fyne.Container); !ok {
				log.Println("could not reach spacing layout widget")
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
			selectedWalletLabel.Text = wallet.Name
			canvas.Refresh(selectedWalletLabel)

			if initFunction != nil {
				initFunction()
			}
			popup.Hide()
		}))
	}

	accountsBoxes[walletIndex] = accountsBox
	groupedWalletsAccounts.Append(accountsBoxes[walletIndex])
	dropdownContent.Append(groupedWalletsAccounts)
}
