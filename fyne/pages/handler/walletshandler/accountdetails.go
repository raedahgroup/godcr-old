package walletshandler

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/pages/handler/multipagecomponents"
	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"

	"github.com/raedahgroup/dcrlibwallet"
)

func (walletPage *WalletPageObject) accountDetailsPopUp(walletIcon *fyne.StaticResource, account *dcrlibwallet.Account, accountNameInMainPage *canvas.Text) {
	var popUp *widget.PopUp
	parentPopUpAccountName := widgets.NewTextWithSize(account.Name, values.DefaultTextColor, 20)
	successLabel := widgets.NewBorderedText("", fyne.NewSize(10, 10), values.Green)

	editAccountButton := widgets.NewImageButton(walletPage.icons[assets.Edit], nil, func() {
		walletPage.renameAccountPopUp(parentPopUpAccountName, accountNameInMainPage,
			account.WalletID, int(account.Number), successLabel, popUp)
	})

	exitButton := widgets.NewImageButton(theme.CancelIcon(), nil, func() {
		popUp.Hide()
	})

	baseWidget := widget.NewHBox(exitButton, widgets.NewHSpacer(values.SpacerSize12),
		parentPopUpAccountName,
		widgets.NewHSpacer(values.SpacerSize202), editAccountButton)

	totalAmountInString := strconv.FormatFloat(dcrlibwallet.AmountCoin(account.TotalBalance), 'f', -1, 64)
	spendableAmountInString := strconv.FormatFloat(dcrlibwallet.AmountCoin(account.Balance.Spendable), 'f', -1, 64)

	// consist of totalbalance label to Imature stake gen
	balanceDetailsBox := widget.NewVBox(
		canvas.NewText(values.TotalBalance, values.SpendableLabelColor),
		widgets.NewVSpacer(values.SpacerSize12),
		multipagecomponents.AmountFormatBox(spendableAmountInString, values.TextSize22, values.TextSize14),
		canvas.NewText(values.Spendable, values.TransactionInfoColor),
		widgets.NewVSpacer(values.SpacerSize12),
	)
	walletPage.addStakingBalance(balanceDetailsBox, account.Balance)

	accountDetailsSpacer := widget.NewIcon(walletIcon).MinSize().Width + values.SpacerSize20

	accountProperties := accountPropertiesBox(account, baseWidget.MinSize().Width)
	accountProperties.Hide()

	showPropertiesText := widgets.NewTextWithStyle(values.ShowProperties, values.Blue, fyne.TextStyle{}, fyne.TextAlignCenter, 16)
	clickablePropertiesText := widgets.NewClickableWidget(widget.NewVBox(showPropertiesText), func() {
		if accountProperties.Hidden {
			accountProperties.Show()
			showPropertiesText.Text = values.HideProperties
		} else {
			accountProperties.Hide()
			showPropertiesText.Text = values.ShowProperties
		}

		showPropertiesText.Refresh()
		walletPage.WalletPageContents.Refresh()
	})

	accountBalanceBox := widget.NewVBox(
		widgets.NewVSpacer(values.SpacerSize14),
		baseWidget,
		widgets.NewVSpacer(values.SpacerSize12),
		widget.NewHBox(widgets.CenterObject(widget.NewIcon(walletIcon), false), widgets.NewHSpacer(values.SpacerSize14),
			multipagecomponents.AmountFormatBox(totalAmountInString, 32, 20)),
		widget.NewHBox(widgets.NewHSpacer(accountDetailsSpacer), balanceDetailsBox),

		accountProperties,
		canvas.NewLine(values.StrippedLineColor),

		widgets.NewVSpacer(values.SpacerSize12),
		clickablePropertiesText,
		widgets.NewVSpacer(values.SpacerSize12),
	)

	popupContent := widget.NewHBox(widgets.NewHSpacer(values.SpacerSize14), accountBalanceBox, widgets.NewHSpacer(values.SpacerSize14))

	popUp = widget.NewModalPopUp(popupContent, walletPage.Window.Canvas())
}

func (walletPage *WalletPageObject) renameAccountPopUp(accountNameInParentPopup, accountNameInPage *canvas.Text,
	walletID, accountID int, successLabel *widgets.BorderedText, accountPopUp *widget.PopUp) {

	errorLabel := canvas.NewText("", values.ErrorColor)
	errorLabel.Hide()

	var popup *widget.PopUp

	showAccountPage := func() {
		popup.Hide()
		accountPopUp.Show()
	}

	renameAccountLabel := canvas.NewText(values.RenameAccount, values.DefaultTextColor)
	renameAccountLabel.TextStyle.Bold = true

	accountTextBox := widget.NewEntry()
	accountTextBox.PlaceHolder = "Account name"

	cancelButton := widgets.NewClickableWidget(widget.NewVBox(widgets.NewTextWithStyle(values.Cancel, values.Blue, fyne.TextStyle{Bold: true}, fyne.TextAlignLeading, 16)),
		func() {
			showAccountPage()
		})

	renameButton := widgets.NewButton(values.Blue, "Rename", func() {
		wallet := walletPage.MultiWallet.WalletWithID(walletID)
		if wallet == nil {
			errorLabel.Text = "could not initialize wallet for account"
			errorLabel.Show()
			return
		}

		err := wallet.RenameAccount(int32(accountID), accountTextBox.Text)
		if err != nil {
			errorLabel.Text = err.Error()
			errorLabel.Show()
			return
		}

		accountNameInParentPopup.Text = accountTextBox.Text
		accountNameInPage.Text = accountTextBox.Text

		showAccountPage()
		successLabel.SetText("Account renamed")
		successLabel.Container.Show()
		time.AfterFunc(time.Second*5, func() {
			successLabel.Container.Hide()
		})

		walletPage.WalletPageContents.Refresh()
	})
	renameButton.SetMinSize(renameButton.MinSize().Add(fyne.NewSize(32, 24)))
	renameButton.SetTextStyle(fyne.TextStyle{Bold: true})

	accountTextBox.OnChanged = func(value string) {
		if value == "" {
			renameButton.Disable()
		} else {
			renameButton.Enable()
		}
	}

	popUpContent := widget.NewVBox(
		widgets.NewVSpacer(values.Padding),
		renameAccountLabel,
		widgets.NewVSpacer(values.SpacerSize20),
		accountTextBox,
		errorLabel,
		widgets.NewVSpacer(values.SpacerSize20),
		widget.NewHBox(widgets.NewHSpacer(values.SpacerSize140), layout.NewSpacer(), widgets.CenterObject(cancelButton, false),
			widgets.NewHSpacer(20), widgets.CenterObject(renameButton.Container, false)),
		widgets.NewVSpacer(values.Padding))

	contentWithBorder := widget.NewHBox(widgets.NewHSpacer(values.Padding), popUpContent, widgets.NewHSpacer(values.Padding))

	popup = widget.NewModalPopUp(contentWithBorder, walletPage.Window.Canvas())
}

func accountPropertiesBox(account *dcrlibwallet.Account, popupMinSizeWidth int) *widget.Box {
	propertiesInfo := func(ID string, value string) *widget.Box {
		return widget.NewHBox(widgets.NewHSpacer(values.SpacerSize12), canvas.NewText(ID, values.TransactionInfoColor),
			layout.NewSpacer(),
			canvas.NewText(value, values.DefaultTextColor), widgets.NewHSpacer(values.SpacerSize12))
	}

	propertiesVBox := widget.NewVBox(
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(popupMinSizeWidth, 1)), canvas.NewLine(values.StrippedLineColor)),
		widgets.NewVSpacer(values.SpacerSize12),
		widget.NewHBox(widgets.NewHSpacer(values.SpacerSize12), canvas.NewText(values.Properties, values.TransactionInfoColor)),
		widgets.NewVSpacer(values.SpacerSize12),

		propertiesInfo(values.AccountNumber, fmt.Sprintf("%d", account.Number)),
		widgets.NewVSpacer(values.SpacerSize14),

		propertiesInfo(values.HDPath, fmt.Sprintf(values.HDPathFormat, account.Number)),
		widgets.NewVSpacer(values.SpacerSize14),

		propertiesInfo(values.Keys, fmt.Sprintf(values.KeyValues,
			account.ExternalKeyCount, account.InternalKeyCount, account.ImportedKeyCount)),
		widgets.NewVSpacer(values.SpacerSize14),
	)

	return propertiesVBox
}

func (walletPage *WalletPageObject) addStakingBalance(box *widget.Box, balance *dcrlibwallet.Balance) {
	if balance.ImmatureReward == 0 && balance.ImmatureStakeGeneration == 0 && balance.LockedByTickets == 0 && balance.VotingAuthority == 0 {
		return
	}
	walletPage.Window.Resize(
		walletPage.TabMenu.MinSize().Add(fyne.NewSize(0, 250)))

	immatureRewardBalance := strconv.FormatFloat(dcrlibwallet.AmountCoin(balance.ImmatureReward), 'f', -1, 64)
	lockedByTickets := strconv.FormatFloat(dcrlibwallet.AmountCoin(balance.LockedByTickets), 'f', -1, 64)
	votingAuthority := strconv.FormatFloat(dcrlibwallet.AmountCoin(balance.VotingAuthority), 'f', -1, 64)
	immatureStakeGen := strconv.FormatFloat(dcrlibwallet.AmountCoin(balance.ImmatureStakeGeneration), 'f', -1, 64)

	box.Append(multipagecomponents.AmountFormatBox(immatureRewardBalance, values.TextSize22, values.TextSize14))
	box.Append(canvas.NewText(values.ImmatureRewards, values.TransactionInfoColor))
	box.Append(widgets.NewHSpacer(values.SpacerSize12))

	box.Append(multipagecomponents.AmountFormatBox(lockedByTickets, values.TextSize22, values.TextSize14))
	box.Append(canvas.NewText(values.LockedByTickets, values.TransactionInfoColor))
	box.Append(widgets.NewHSpacer(values.SpacerSize12))

	box.Append(multipagecomponents.AmountFormatBox(votingAuthority, values.TextSize22, values.TextSize14))
	box.Append(canvas.NewText(values.VotingAuthority, values.TransactionInfoColor))
	box.Append(widgets.NewHSpacer(values.SpacerSize12))

	box.Append(multipagecomponents.AmountFormatBox(immatureStakeGen, values.TextSize22, values.TextSize14))
	box.Append(canvas.NewText(values.ImmatureStakeGen, values.TransactionInfoColor))
	box.Append(widgets.NewHSpacer(values.SpacerSize12))
}
