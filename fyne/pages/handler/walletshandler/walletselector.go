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
	icons, err := assets.GetIcons(assets.Expand, assets.CollapseIcon, assets.WalletIcon, assets.ImportedAccount, assets.MoreIcon)
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

	notBackedUpLabel := canvas.NewText("Not backed up", values.ErrorColor)
	// add extra padding to account selector on hiding "Not backed up" due to extra VBox padding
	extraPadding1 := widgets.NewVSpacer((notBackedUpLabel.MinSize().Height / 2) - 1)
	extraPadding2 := widgets.NewVSpacer((notBackedUpLabel.MinSize().Height / 2) - 1)

	if selectedWallet.Seed == "" {
		notBackedUpLabel.Hide()
	} else {
		extraPadding2.Hide()
		extraPadding1.Hide()
	}

	accountLabel := widgets.NewVBox(layout.NewSpacer(), canvas.NewText(selectedWallet.Name, values.DefaultTextColor), notBackedUpLabel, layout.NewSpacer())

	expandIcon := widget.NewIcon(icons[assets.Expand])

	walletIcon := canvas.NewImageFromResource(icons[assets.WalletIcon])
	walletIcon.SetMinSize(fyne.NewSize(24, 24))

	accountBox := widgets.NewHBox(
		widgets.NewHSpacer(12),
		widgets.NewVBox(layout.NewSpacer(), expandIcon, layout.NewSpacer()),
		widgets.NewHSpacer(4),
		widgets.NewVBox(layout.NewSpacer(), walletIcon, layout.NewSpacer()),
		widgets.NewHSpacer(12),
		accountLabel, widgets.NewHSpacer(50),
		layout.NewSpacer(),
		walletPage.WalletTotalAmountLabel[index], widgets.NewHSpacer(4),
		widgets.NewImageButton(icons[assets.MoreIcon], nil, func() {

		}), widgets.NewHSpacer(12))

	toShow := widgets.NewVBox(
		widget.NewLabel("To do"))
	toShow.Hide()

	accountSelector := widgets.NewClickableWidget(accountBox, func() {
		if toShow.Hidden {
			expandIcon.SetResource(icons[assets.CollapseIcon])
			toShow.Show()
		} else {
			expandIcon.SetResource(icons[assets.Expand])
			toShow.Hide()
		}
	})

	textBox := widgets.NewVBox(
		extraPadding1,
		widgets.NewVSpacer(12),
		accountSelector,
		extraPadding2,
		toShow,
		widgets.NewVSpacer(4),
	)
	walletPage.walletSelectorBox.Append(widget.NewVBox(
		textBox,
		widgets.NewVSpacer(4),
	))
}
