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
	"github.com/raedahgroup/godcr/fyne/layouts"
	"github.com/raedahgroup/godcr/fyne/pages/handler/constantvalues"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type AccountSelectorStruct struct {
	OnAccountChange         func()
	SendingSelectedWalletID *int
	WalletIDs               []int

	AccountBoxes []*widget.Box

	SelectedAccountLabel        *widget.Label
	SelectedAccountBalanceLabel *widget.Label
	SelectedWalletLabel         *canvas.Text

	PageContents *widget.Box

	SelectedWallet *dcrlibwallet.Wallet
	MultiWallet    *dcrlibwallet.MultiWallet

	Window fyne.Window
}

func (accountSelector *AccountSelectorStruct) CreateAccountSelector(accountLabel string) (*widgets.ClickableBox, error) {
	icons, err := assets.GetIcons(assets.ReceiveAccountIcon, assets.CollapseIcon)
	if err != nil {
		return nil, errors.New(constantvalues.AccountSelectorIconErr)
	}

	accountSelector.SelectedWallet = accountSelector.MultiWallet.WalletWithID(accountSelector.WalletIDs[0])
	accountSelector.SelectedWalletLabel = canvas.NewText(accountSelector.SelectedWallet.Name, color.RGBA{137, 151, 165, 255})

	dropdownContent := widget.NewVBox()

	selectAccountBox := widget.NewHBox(
		widgets.NewHSpacer(15),
		widget.NewVBox(widgets.NewVSpacer(10), widget.NewIcon(icons[assets.ReceiveAccountIcon])),
		widgets.NewHSpacer(20),
		fyne.NewContainerWithLayout(layouts.NewVBox(12), accountSelector.SelectedAccountLabel, accountSelector.SelectedWalletLabel),
		widgets.NewHSpacer(30),
		widget.NewVBox(widgets.NewVSpacer(4), accountSelector.SelectedAccountBalanceLabel),
		widgets.NewHSpacer(8),
		widget.NewVBox(widgets.NewVSpacer(6), widget.NewIcon(icons[assets.CollapseIcon])),
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
		accountSelector.PageContents.Refresh()
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
		if account.Name == constantvalues.Imported {
			continue
		}

		spendableLabel := canvas.NewText(constantvalues.Spendable, color.Black)
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
		if index != 0 || walletID != *accountSelector.SendingSelectedWalletID {
			checkmarkIcon.Hide()
			spacing = widgets.NewHSpacer(35)
		} else {
			spacing = widgets.NewHSpacer(15)
		}

		accountsView := widget.NewHBox(
			widgets.NewHSpacer(15),
			widget.NewIcon(receiveAccountIcon),
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
						accountSelector.SelectedAccountBalanceLabel.SetText(balanceLabel.Text)
					}
				}
			}

			accountSelector.SelectedAccountLabel.SetText(accountName)
			accountSelector.SelectedWalletLabel.Text = wallet.Name
			accountSelector.SelectedWalletLabel.Refresh()

			if accountSelector.OnAccountChange != nil {
				accountSelector.OnAccountChange()
			}
			popup.Hide()
		}))
	}

	accountSelector.AccountBoxes[walletIndex] = accountsBox
	groupedWalletsAccounts.Append(accountSelector.AccountBoxes[walletIndex])
	dropdownContent.Append(groupedWalletsAccounts)
}
