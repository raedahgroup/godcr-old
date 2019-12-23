package multipagecomponents

import (
	"errors"
	"image/color"
	"log"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type AccountSelectorStruct struct {
	OnAccountChange          func()
	SendingSelectedWalletID  *int
	SendingSelectedAccountID *int
	WalletIDs                []int

	AccountBoxes []*widget.Box

	SelectedAccountLabel        *canvas.Text
	SelectedAccountBalanceLabel *canvas.Text
	SelectedWalletLabel         *canvas.Text

	PageContents *widget.Box

	SelectedWallet *dcrlibwallet.Wallet
	MultiWallet    *dcrlibwallet.MultiWallet

	Window fyne.Window
}

func (accountSelector *AccountSelectorStruct) CreateAccountSelector(accountLabel string) (*widgets.ClickableBox, error) {
	icons, err := assets.GetIcons(assets.ReceiveAccountIcon, assets.CollapseIcon)
	if err != nil {
		return nil, errors.New(values.AccountSelectorIconErr)
	}

	accountSelector.SelectedWallet = accountSelector.MultiWallet.WalletWithID(accountSelector.WalletIDs[0])
	accountSelector.SelectedWalletLabel = canvas.NewText(accountSelector.SelectedWallet.Name, values.WalletLabelColor)

	dropdownContent := widget.NewVBox()

	selectAccountBox := widgets.NewHBox(
		widgets.NewHSpacer(values.SpacerSize16),
		widgets.NewVBox(layout.NewSpacer(), widget.NewIcon(icons[assets.ReceiveAccountIcon]), layout.NewSpacer()),
		widgets.NewHSpacer(values.SpacerSize20),
		widgets.NewVBox(widgets.NewHBox(layout.NewSpacer(), accountSelector.SelectedAccountLabel, layout.NewSpacer()),
			widgets.NewHBox(layout.NewSpacer(), widgets.NewHSpacer(values.NilSpacer), accountSelector.SelectedWalletLabel, layout.NewSpacer())),
		widgets.NewHSpacer(values.SpacerSize30),
		layout.NewSpacer(),
		widgets.NewVBox(layout.NewSpacer(), accountSelector.SelectedAccountBalanceLabel, layout.NewSpacer()),
		widgets.NewHSpacer(values.SpacerSize8),
		widgets.NewVBox(layout.NewSpacer(), widget.NewIcon(icons[assets.CollapseIcon]), layout.NewSpacer()),
		widgets.NewHSpacer(values.SpacerSize16),
	)

	var accountSelectionPopup *widget.PopUp
	accountSelectionPopupHeader := widget.NewVBox(
		widgets.NewVSpacer(values.SpacerSize4),
		widget.NewHBox(
			widgets.NewHSpacer(values.SpacerSize16),
			widgets.NewImageButton(theme.CancelIcon(), nil, func() { accountSelectionPopup.Hide() }),
			widgets.NewHSpacer(values.SpacerSize16),
			widget.NewLabelWithStyle(accountLabel, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			layout.NewSpacer(),
		),
		widgets.NewVSpacer(values.SpacerSize4),
		canvas.NewLine(color.Black),
	)

	popupContent := widget.NewVBox(accountSelectionPopupHeader)
	accountSelectionPopup = widget.NewPopUp(popupContent, accountSelector.Window.Canvas())
	accountSelectionPopup.Hide()

	// we cant access the children of group widget, proposed hack is to
	// create a vertical box array where all accounts would be placed,
	// then when we want to hide checkmarks we call all children of accountbox and hide checkmark icon except selected
	for walletIndex, walletID := range accountSelector.WalletIDs {
		accountSelector.getAllWalletAccountsInBox(icons[assets.ReceiveAccountIcon], dropdownContent, accountSelector.MultiWallet.WalletWithID(walletID),
			walletIndex, walletID, accountSelectionPopup)
	}

	dropdownContentWithScroller := fyne.NewContainerWithLayout(
		layout.NewFixedGridLayout(fyne.NewSize(dropdownContent.MinSize().Width+5, fyne.Min(dropdownContent.MinSize().Height, 100))),
		widget.NewScrollContainer(dropdownContent))
	popupContent.Append(dropdownContentWithScroller)

	var accountClickableBox *widgets.ClickableBox
	accountClickableBox = widgets.NewClickableBox(selectAccountBox, func() {
		accountSelectionPopup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(
			accountClickableBox).Add(fyne.NewPos(0, accountClickableBox.Size().Height)))

		accountSelectionPopup.Show()
		accountSelectionPopup.Resize(dropdownContentWithScroller.Size().Add(fyne.NewSize(10, accountSelectionPopupHeader.MinSize().Height)))
	})

	return accountClickableBox, err
}

func (accountSelector *AccountSelectorStruct) getAllWalletAccountsInBox(receiveAccountIcon fyne.Resource, dropdownContent *widget.Box,
	wallet *dcrlibwallet.Wallet, walletIndex, walletID int, popup *widget.PopUp) {

	accounts, err := wallet.GetAccountsRaw(dcrlibwallet.DefaultRequiredConfirmations)
	if err != nil {
		return
	}

	var groupedWalletsAccounts = widget.NewGroup(wallet.Name)
	// we cant access children of a group so a box is used
	accountsBox := widget.NewVBox()

	for index, account := range accounts.Acc {
		if account.Name == values.Imported {
			continue
		}

		spendableLabel := canvas.NewText(values.Spendable, values.TransactionInfoColor)
		spendableLabel.TextSize = 10

		accountName := account.Name
		accountNameLabel := widget.NewLabel(accountName)
		accountNameLabel.Alignment = fyne.TextAlignLeading
		accountNameBox := widget.NewVBox(
			accountNameLabel,
			widget.NewHBox(widgets.NewHSpacer(values.NilSpacer), spendableLabel),
		)

		spendableAmountLabel := canvas.NewText(dcrutil.Amount(account.Balance.Spendable).String(), values.DefaultTextColor)
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
		if index != 0 || walletID != *accountSelector.SendingSelectedWalletID {
			checkmarkIcon.Hide()
			spacing = widgets.NewHSpacer(values.SpacerSize36)
		} else {
			spacing = widgets.NewHSpacer(values.SpacerSize16)
		}

		accountsView := widget.NewHBox(
			widgets.NewHSpacer(values.SpacerSize16),
			widget.NewIcon(receiveAccountIcon),
			widgets.NewHSpacer(values.SpacerSize20),
			accountNameBox,
			layout.NewSpacer(),
			widgets.NewHSpacer(values.SpacerSize30),
			accountBalanceBox,
			widgets.NewHSpacer(values.SpacerSize30),
			checkmarkIcon,
			spacing,
		)

		accountNumber := account.Number
		accountsBox.Append(widgets.NewClickableWidget(accountsView, func() {
			if *accountSelector.SendingSelectedAccountID == int(accountNumber) && *accountSelector.SendingSelectedWalletID == walletID {
				if accountSelector.OnAccountChange != nil {
					accountSelector.OnAccountChange()
					popup.Hide()
				}
				return
			}

			*accountSelector.SendingSelectedAccountID = int(accountNumber)
			*accountSelector.SendingSelectedWalletID = walletID
			accountSelector.SelectedWallet = accountSelector.MultiWallet.WalletWithID(walletID)

			for _, boxes := range accountSelector.AccountBoxes {
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
						accountSelector.SelectedAccountBalanceLabel.Text = balanceLabel.Text
					}
				}
			}

			accountSelector.SelectedAccountLabel.Text = accountName
			accountSelector.SelectedWalletLabel.Text = wallet.Name

			popup.Hide()

			accountSelector.SelectedWalletLabel.Refresh()
			accountSelector.SelectedAccountLabel.Refresh()
			accountSelector.PageContents.Refresh()

			if accountSelector.OnAccountChange != nil {
				accountSelector.OnAccountChange()
			}

		}))
	}

	accountSelector.AccountBoxes[walletIndex] = accountsBox
	groupedWalletsAccounts.Append(accountSelector.AccountBoxes[walletIndex])
	dropdownContent.Append(groupedWalletsAccounts)
}
