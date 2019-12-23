package walletshandler

import (
	"fmt"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (walletPage *WalletPageObject) accountSelector() error {
	icons, err := assets.GetIcons(assets.Expand, assets.AccountsIcon, assets.ImportedAccount, assets.MoreIcon)
	if err != nil {
		return err
	}

	// all wallet selectors are housed here,
	// if a wallet is deleted, we are to redeclare the children to omit the deleted wallet
	walletPage.walletSelectorBox = widget.NewVBox()

	walletPage.walletIDToIndex = make(map[int]int)

	for index, walletID := range walletPage.OpenedWallets {
		walletPage.walletIDToIndex[index] = walletID
		walletPage.getAccountsInWallet(icons, index, walletID)
	}

	walletPage.WalletPageContents.Append(walletPage.walletSelectorBox)
	return nil
}

func (walletPage *WalletPageObject) getAccountsInWallet(icons map[string]*fyne.StaticResource, index, selectedWalletID int) {
	selectedWallet := walletPage.MultiWallet.WalletWithID(selectedWalletID)
	accts, err := selectedWallet.GetAccountsRaw(dcrlibwallet.DefaultRequiredConfirmations)
	if err != nil {
		return
	}

	var totalBalance int64
	for _, acc := range accts.Acc {
		totalBalance += acc.TotalBalance
	}
	balanceInString := strconv.FormatFloat(dcrlibwallet.AmountCoin(totalBalance), 'f', 8, 64)

	walletPage.WalletTotalAmountLabel[index].Text = fmt.Sprintf(values.AmountInDCR, balanceInString)

	var accountLabel fyne.CanvasObject

	notBackedUpLabel := canvas.NewText("Not backed up", values.ErrorColor)
	extraPadding1 := widgets.NewVSpacer((notBackedUpLabel.MinSize().Height / 2) - 1)
	extraPadding2 := widgets.NewVSpacer((notBackedUpLabel.MinSize().Height / 2) - 1)

	if selectedWallet.Seed == "" {
		notBackedUpLabel.Hide()
		accountLabel = widgets.NewVBox(layout.NewSpacer(), canvas.NewText(selectedWallet.Name, values.DefaultTextColor), layout.NewSpacer())
	} else {
		extraPadding2.Hide()
		extraPadding1.Hide()
	}

	accountLabel = widgets.NewVBox(layout.NewSpacer(), canvas.NewText(selectedWallet.Name, values.DefaultTextColor), notBackedUpLabel, layout.NewSpacer())

	accountBox := widgets.NewHBox(
		widgets.NewHSpacer(12),
		widget.NewIcon(icons[assets.Expand]), widgets.NewHSpacer(4),
		widget.NewIcon(icons[assets.AccountsIcon]), widgets.NewHSpacer(12),
		accountLabel, widgets.NewHSpacer(50),
		layout.NewSpacer(),
		walletPage.WalletTotalAmountLabel[index], widgets.NewHSpacer(4),
		widgets.NewImageButton(icons[assets.MoreIcon], nil, func() {

		}), widgets.NewHSpacer(12))

	fmt.Println(accountBox.MinSize())
	toShow := widgets.NewVBox(
		widget.NewLabel("Hello"),
		widget.NewLabel("Hello"),
		widget.NewLabel("Hello"),
	)
	toShow.Hide()

	accountSelector := widgets.NewClickableWidget(accountBox, func() {
		fmt.Println("Hello")
		if toShow.Hidden {
			toShow.Show()
		} else {
			toShow.Hide()
		}
		walletPage.WalletPageContents.Refresh()
	})

	textBox := widgets.NewVBox(
		extraPadding1,
		widgets.NewVSpacer(12),
		accountSelector,
		toShow,
		widgets.NewVSpacer(12),
		extraPadding2,
	)
	walletPage.walletSelectorBox.Append(widget.NewVBox(
		textBox,
		widgets.NewVSpacer(4),
	))
}
