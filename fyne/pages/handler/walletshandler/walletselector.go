package walletshandler

import (
	"fmt"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (walletPage *WalletPageObject) accountSelector() error {
	// all wallet selectors are housed here,
	// if a wallet is deleted, we are to redeclare the children to omit the deleted wallet
	walletPage.walletSelectorBox = widget.NewVBox()

	walletPage.WalletIDToIndex = make(map[int]int)

	for index, walletID := range walletPage.OpenedWallets {
		walletPage.WalletIDToIndex[index] = walletID

		err := walletPage.getAccountsInWallet(index, walletID)
		if err != nil {
			return err
		}
	}

	scrollableSelectorBox := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(
		fyne.NewSize(walletPage.walletSelectorBox.MinSize().Width+20, walletPage.Window.Content().MinSize().Height)),
		widget.NewScrollContainer(widget.NewHBox(walletPage.walletSelectorBox, widgets.NewHSpacer(values.Padding))))

	walletPage.WalletPageContents.Append(scrollableSelectorBox)
	return nil
}

func (walletPage *WalletPageObject) getAccountsInWallet(index, selectedWalletID int) error {
	selectedWallet := walletPage.MultiWallet.WalletWithID(selectedWalletID)
	accts, err := selectedWallet.GetAccountsRaw(dcrlibwallet.DefaultRequiredConfirmations)
	if err != nil {
		return err
	}

	var totalBalance int64
	for _, acc := range accts.Acc {
		totalBalance += acc.TotalBalance
	}
	balanceInString := strconv.FormatFloat(dcrlibwallet.AmountCoin(totalBalance), 'f', -1, 64)

	walletPage.WalletTotalAmountLabel[index].Text = fmt.Sprintf(values.AmountInDCR, balanceInString)

	notBackedUpLabel := canvas.NewText("Not backed up", values.ErrorColor)

	// add extra padding to account selector on hiding "Not backed up" due to extra VBox padding
	extraPadding1 := widgets.NewVSpacer((notBackedUpLabel.MinSize().Height / 2) - 2)
	extraPadding2 := widgets.NewVSpacer((notBackedUpLabel.MinSize().Height / 2) - 3)

	if selectedWallet.Seed == "" {
		notBackedUpLabel.Hide()
	} else {
		extraPadding2.Hide()
		extraPadding1.Hide()
	}

	accountLabel := widgets.NewVBox(
		layout.NewSpacer(),
		widgets.CenterObject(canvas.NewText(selectedWallet.Name, values.DefaultTextColor), true),
		widgets.CenterObject(notBackedUpLabel, true),
		layout.NewSpacer())

	walletPage.walletExpandCollapseIcon[index] = widget.NewIcon(walletPage.icons[assets.Expand])

	walletIcon := widget.NewIcon(walletPage.icons[assets.WalletIcon])

	accountBox := widgets.NewHBox(
		widgets.NewHSpacer(values.SpacerSize12),
		widgets.CenterObject(walletPage.walletExpandCollapseIcon[index], true),
		widgets.NewHSpacer(values.SpacerSize4),
		widgets.CenterObject(walletIcon, true),
		widgets.NewHSpacer(values.SpacerSize12),
		widgets.CenterObject(accountLabel, true),
		widgets.NewHSpacer(values.SpacerSize50),
		layout.NewSpacer(),
		widgets.CenterObject(walletPage.WalletTotalAmountLabel[index], true),
		widgets.NewHSpacer(4),
		widgets.NewImageButton(walletPage.icons[assets.MoreIcon], nil, func() {

		}),
		widgets.NewHSpacer(values.SpacerSize12),
	)

	accountBoxSpacer := accountBox.MinSize().Width - values.SpacerSize44

	walletPage.walletAccountsBox[index], err = walletPage.accountDropdown(accountBoxSpacer, selectedWallet)
	if err != nil {
		return err
	}
	walletPage.walletAccountsBox[index].Hide()

	accountSelector := widgets.NewClickableWidget(accountBox, func() {
		// hide other multiwallet accounts boxes
		for i, propertieBox := range walletPage.walletAccountsBox {
			if i == index {
				continue
			}

			if !propertieBox.Hidden {
				propertieBox.Hide()
				walletPage.walletExpandCollapseIcon[i].SetResource(walletPage.icons[assets.Expand])
			}
		}

		if walletPage.walletAccountsBox[index].Hidden {
			walletPage.walletExpandCollapseIcon[index].SetResource(walletPage.icons[assets.CollapseIcon])
			walletPage.walletAccountsBox[index].Show()
		} else {
			walletPage.walletExpandCollapseIcon[index].SetResource(walletPage.icons[assets.Expand])
			walletPage.walletAccountsBox[index].Hide()
		}

	})

	textBox := widgets.NewVBox(
		extraPadding1,
		widgets.NewVSpacer(values.SpacerSize12),
		accountSelector,
		extraPadding2,
		widgets.NewVSpacer(values.SpacerSize4),
		walletPage.walletAccountsBox[index],
		widgets.NewVSpacer(values.SpacerSize4))

	walletPage.walletSelectorBox.Append(widget.NewVBox(
		textBox,
		widgets.NewVSpacer(values.SpacerSize4), // add spacing between wallet account selector
	))

	return nil
}

func (walletPage *WalletPageObject) accountDropdown(walletBoxSize int, wallet *dcrlibwallet.Wallet) (*widgets.Box, error) {
	accounts, err := wallet.GetAccountsRaw(dcrlibwallet.DefaultRequiredConfirmations)
	if err != nil {
		return nil, err
	}

	accountObjects := widgets.NewVBox(canvas.NewLine(values.StrippedLineColor))

	for _, account := range accounts.Acc {
		accountObjects.Append(walletPage.walletAccountBox(walletBoxSize, account))
	}

	addAccount := widgets.NewClickableBox(widgets.NewHBox(
		widgets.CenterObject(widget.NewIcon(theme.ContentAddIcon()), true),
		widgets.NewHSpacer(values.SpacerSize12),
		widgets.CenterObject(widget.NewLabel("Add new account"), true)),
		func() {
			fmt.Println("Add new account")
		})

	accountObjects.Append(widgets.NewVSpacer(values.SpacerSize12))
	accountObjects.Append(addAccount)
	accountObjects.Append(widgets.NewVSpacer(values.SpacerSize12))

	accountBoxSpacer := widget.NewIcon(walletPage.icons[assets.Expand]).MinSize().Width + values.SpacerSize24
	accountBoxes := widgets.NewHBox(widgets.NewHSpacer(accountBoxSpacer), accountObjects)

	return accountBoxes, nil
}

func (walletPage *WalletPageObject) walletAccountBox(walletBoxSize int, account *dcrlibwallet.Account) *widgets.Box {
	walletIcon := walletPage.icons[assets.WalletAccount]
	if account.Name == "imported" {
		walletIcon = walletPage.icons[assets.ImportedAccount]
	}
	accountName := canvas.NewText(account.Name, values.DefaultTextColor)

	accountNameWithSpendableLabel := widgets.NewVBox(
		layout.NewSpacer(),
		accountName,
		widgets.NewHSpacer(values.SpacerSize4),
		canvas.NewText("Spendable", values.SpendableLabelColor),
		layout.NewSpacer())

	totalBalanceInString := strconv.FormatFloat(dcrlibwallet.AmountCoin(account.TotalBalance), 'f', -1, 64)
	spendableBalanceInString := strconv.FormatFloat(dcrlibwallet.AmountCoin(account.Balance.Spendable), 'f', -1, 64)

	accountBalAndSpendableBal := widgets.NewVBox(
		layout.NewSpacer(),
		widgets.NewTextAndAlign(fmt.Sprintf(values.AmountInDCR, totalBalanceInString), values.DefaultTextColor, fyne.TextAlignTrailing),
		widgets.NewHSpacer(values.SpacerSize4),
		widgets.NewTextAndAlign(fmt.Sprintf(values.AmountInDCR, spendableBalanceInString), values.SpendableLabelColor, fyne.TextAlignTrailing),
		layout.NewSpacer())

	accountHBox := widgets.NewHBox(widgets.CenterObject(widget.NewIcon(walletIcon), true), widgets.NewHSpacer(values.SpacerSize14),
		accountNameWithSpendableLabel,
		layout.NewSpacer(),
		widgets.NewHSpacer(values.NilSpacer),
		accountBalAndSpendableBal, widgets.NewHSpacer(values.SpacerSize12))

	spacerSize := walletBoxSize - accountHBox.MinSize().Width - 4
	accountHBox.Children[4] = widgets.NewHSpacer(spacerSize)

	iconWidthSize := widget.NewIcon(walletPage.icons[assets.WalletAccount]).MinSize().Width + values.SpacerSize14

	accountBoxWithLiner := widgets.NewVBox(
		widgets.NewVSpacer(values.SpacerSize14),
		accountHBox,
		widgets.NewVSpacer(values.SpacerSize8))

	clickableAccountBox := widgets.NewClickableBox(accountBoxWithLiner, func() {
		walletPage.accountDetailsPopUp(walletIcon, account, accountName)
	})

	canvasLine := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(clickableAccountBox.MinSize().Width-iconWidthSize-values.SpacerSize18, 1)),
		canvas.NewLine(values.StrippedLineColor))

	return widgets.NewVBox(
		clickableAccountBox,
		widgets.NewHBox(widgets.NewHSpacer(iconWidthSize), canvasLine))
}
