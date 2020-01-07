package walletshandler

import (
	"fmt"
	"strconv"

	"github.com/raedahgroup/godcr/fyne/pages/handler/multipagecomponents"

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

	walletLabel := canvas.NewText(selectedWallet.Name, values.DefaultTextColor)

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

	var clickableMoreDialog *widgets.ImageButton
	clickableMoreDialog = widgets.NewImageButton(walletPage.icons[assets.MoreIcon], nil, func() {
		dialogPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(clickableMoreDialog)
		walletPage.dialogMenu(walletLabel, dialogPos, selectedWalletID)
	})

	accountLabel := widgets.NewVBox(
		layout.NewSpacer(),
		widgets.CenterObject(walletLabel, true),
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
		clickableMoreDialog,
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
	var accountBoxes *widgets.Box

	var addAccount *widgets.ClickableBox

	addAccount = widgets.NewClickableBox(widgets.NewHBox(
		widgets.CenterObject(widget.NewIcon(theme.ContentAddIcon()), true),
		widgets.NewHSpacer(values.SpacerSize12),
		widgets.CenterObject(canvas.NewText("Add new account", values.DefaultTextColor), true)),
		func() {
			addAccountFunc := func(account *dcrlibwallet.Account) {
				child := accountObjects.Children
				// omit to import account clickable widget
				allAccounts := child[:len(child)-4]
				fmt.Println(len(allAccounts))
				importedAccount := child[len(child)-4]

				allAccounts = append(allAccounts, walletPage.walletAccountBox(walletBoxSize, account))
				accountObjects.Children = allAccounts
				accountObjects.Children = append(accountObjects.Children, importedAccount)

				accountObjects.Append(widgets.NewVSpacer(values.SpacerSize12))
				accountObjects.Append(addAccount)
				accountObjects.Append(widgets.NewVSpacer(values.SpacerSize12))

				accountObjects.Refresh()
				walletPage.WalletPageContents.Refresh()
			}

			walletPage.createNewAccountPopUp(wallet, addAccountFunc)
		})

	accountObjects.Append(widgets.NewVSpacer(values.SpacerSize12))
	accountObjects.Append(addAccount)
	accountObjects.Append(widgets.NewVSpacer(values.SpacerSize12))

	accountBoxSpacer := widget.NewIcon(walletPage.icons[assets.Expand]).MinSize().Width + values.SpacerSize24
	accountBoxes = widgets.NewHBox(widgets.NewHSpacer(accountBoxSpacer), accountObjects)

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

func (walletPage *WalletPageObject) createNewAccountPopUp(wallet *dcrlibwallet.Wallet, addAccount func(*dcrlibwallet.Account)) {
	createNewAccountLabel := widgets.NewTextWithStyle(values.CreateNewAccount, values.DefaultTextColor,
		fyne.TextStyle{Bold: true}, fyne.TextAlignLeading, 20)

	infoLabel := widget.NewHBox(
		widget.NewIcon(walletPage.icons[assets.Alert]),
		canvas.NewText("Accounts CANNOT be deleted once created.", values.DarkerBlueGrayTextColor),
	)

	textBox := widget.NewEntry()
	textBox.SetPlaceHolder(values.AccountNamePlaceHolder)

	var popup *widget.PopUp

	cancel := canvas.NewText(values.Cancel, values.Blue)
	cancel.TextStyle.Bold = true

	cancelButton := widgets.NewClickableWidget(widget.NewVBox(cancel), func() {
		popup.Hide()
	})

	var accountNo int32
	initOnConfirmation := func(value string) (err error) {
		accountNo, err = wallet.NextAccount(textBox.Text, []byte(value))
		return err
	}
	extraCall := func() {
		account, err := wallet.GetAccount(accountNo, dcrlibwallet.DefaultRequiredConfirmations)
		if err == nil {
			addAccount(account)
		}

		walletPage.showLabel("Account created", walletPage.successLabel)
	}
	onCancel := func() {
		popup.Show()
	}

	createAccountButton := widgets.NewButton(values.Blue, values.CreateNewAccountButtonText, func() {
		confirmPasswordPopUp := multipagecomponents.PasswordPopUpObjects{
			Window:             walletPage.Window,
			InitOnConfirmation: initOnConfirmation,
			ExtraCalls:         extraCall,
			InitOnCancel:       onCancel,
			Title:              values.ConfirmToCreateAcc,
		}

		confirmPasswordPopUp.PasswordPopUp()
	})
	createAccountButton.SetTextStyle(fyne.TextStyle{Bold: true})
	createAccountButton.SetMinSize(createAccountButton.MinSize().Add(fyne.NewSize(32, 24)))

	popUpContent := widget.NewVBox(
		widgets.NewVSpacer(values.SpacerSize20),
		createNewAccountLabel,
		widgets.NewVSpacer(values.SpacerSize20),
		infoLabel,
		widgets.NewVSpacer(values.SpacerSize20),
		textBox,
		widgets.NewVSpacer(values.SpacerSize20),
		widget.NewHBox(layout.NewSpacer(), widgets.CenterObject(cancelButton, false), widgets.NewHSpacer(values.SpacerSize20), createAccountButton.Container),
		widgets.NewVSpacer(values.SpacerSize20),
	)

	contentWithBorder := widget.NewHBox(widgets.NewHSpacer(values.SpacerSize20), popUpContent, widgets.NewHSpacer(values.SpacerSize20))

	popup = widget.NewModalPopUp(contentWithBorder, walletPage.Window.Canvas())
}
