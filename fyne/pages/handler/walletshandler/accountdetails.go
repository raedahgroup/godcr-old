package walletshandler

import (
	"fmt"
	"strconv"

	"github.com/raedahgroup/godcr/fyne/pages/handler/multipagecomponents"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

func (walletPage *WalletPageObject) accountDetailsPopUp(walletIcon *fyne.StaticResource, account *dcrlibwallet.Account) {
	var popUp *widget.PopUp

	editAccountButton := widgets.NewImageButton(walletPage.icons[assets.Edit], nil, func() {

	})

	exitButton := widgets.NewImageButton(theme.CancelIcon(), nil, func() {
		popUp.Hide()
	})

	baseWidget := widget.NewHBox(exitButton, widgets.NewHSpacer(values.SpacerSize12),
		widgets.NewTextWithSize(account.Name, values.DefaultTextColor, 20),
		widgets.NewHSpacer(values.SpacerSize202), editAccountButton)

	totalAmountInString := strconv.FormatFloat(dcrlibwallet.AmountCoin(account.TotalBalance), 'f', -1, 64)
	spendableAmountInString := strconv.FormatFloat(dcrlibwallet.AmountCoin(account.Balance.Spendable), 'f', -1, 64)

	// consist of totalbalance label to Imature stake gen
	balanceDetailsBox := widget.NewVBox(
		canvas.NewText(values.TotalBalance, values.SpendableLabelColor),
		widgets.NewVSpacer(values.SpacerSize12),
		widgets.NewTextWithSize(spendableAmountInString, values.DefaultTextColor, 20),
		canvas.NewText(values.Spendable, values.TransactionInfoColor),
		widgets.NewVSpacer(values.SpacerSize12),
	)
	addStakingBalance(balanceDetailsBox, account.Balance)

	accountDetailsSpacer := widget.NewIcon(walletIcon).MinSize().Width + values.SpacerSize14

	accountProperties := accountPropertiesBox(account, baseWidget.MinSize().Width)
	accountProperties.Hide()

	showPropertiesText := widgets.NewTextWithStyle(values.ShowProperties, values.Blue, fyne.TextStyle{}, fyne.TextAlignCenter, 16)
	clickablePropertiesText := widgets.NewClickableWidget(widget.NewVBox(showPropertiesText), func() {
		walletPage.WalletPageContents.Refresh()
		if accountProperties.Hidden {
			accountProperties.Show()
			showPropertiesText.Text = values.HideProperties
		} else {
			accountProperties.Hide()
			showPropertiesText.Text = values.ShowProperties
		}
		walletPage.WalletPageContents.Refresh()
	})

	accountBalanceBox := widget.NewVBox(
		baseWidget,

		widget.NewHBox(centerObject(widget.NewIcon(walletIcon)), widgets.NewHSpacer(values.SpacerSize14), multipagecomponents.AmountFormatBox(totalAmountInString, 32, 30)),

		widget.NewHBox(widgets.NewHSpacer(accountDetailsSpacer), balanceDetailsBox),

		canvas.NewLine(values.StrippedLineColor),

		widgets.NewVSpacer(values.SpacerSize12),
		clickablePropertiesText,
		widgets.NewVSpacer(values.SpacerSize12),
	)

	popupContent := widget.NewHBox(widgets.NewHSpacer(values.SpacerSize14), accountBalanceBox, widgets.NewHSpacer(values.SpacerSize14))

	popUp = widget.NewModalPopUp(popupContent, walletPage.Window.Canvas())
}

func addStakingBalance(box *widget.Box, balance *dcrlibwallet.Balance) {
	if balance.ImmatureReward == 0 && balance.ImmatureStakeGeneration == 0 && balance.LockedByTickets == 0 && balance.VotingAuthority == 0 {
		return
	}

	immatureRewardBalance := strconv.FormatFloat(dcrlibwallet.AmountCoin(balance.ImmatureReward), 'f', -1, 64)
	lockedByTickets := strconv.FormatFloat(dcrlibwallet.AmountCoin(balance.LockedByTickets), 'f', -1, 64)
	votingAuthority := strconv.FormatFloat(dcrlibwallet.AmountCoin(balance.VotingAuthority), 'f', -1, 64)
	immatureStakeGen := strconv.FormatFloat(dcrlibwallet.AmountCoin(balance.ImmatureStakeGeneration), 'f', -1, 64)

	box.Append(multipagecomponents.AmountFormatBox(immatureRewardBalance, values.SpacerSize20, values.SpacerSize14))
	box.Append(canvas.NewText(values.ImmatureRewards, values.TransactionInfoColor))
	box.Append(widgets.NewHSpacer(values.SpacerSize12))

	box.Append(multipagecomponents.AmountFormatBox(lockedByTickets, values.SpacerSize20, values.SpacerSize14))
	box.Append(canvas.NewText(values.LockedByTickets, values.TransactionInfoColor))
	box.Append(widgets.NewHSpacer(values.SpacerSize12))

	box.Append(multipagecomponents.AmountFormatBox(votingAuthority, values.SpacerSize20, values.SpacerSize14))
	box.Append(canvas.NewText(values.VotingAuthority, values.TransactionInfoColor))
	box.Append(widgets.NewHSpacer(values.SpacerSize12))

	box.Append(multipagecomponents.AmountFormatBox(immatureStakeGen, values.SpacerSize20, values.SpacerSize14))
	box.Append(canvas.NewText(values.ImmatureStakeGen, values.TransactionInfoColor))
	box.Append(widgets.NewHSpacer(values.SpacerSize12))
}

func accountPropertiesBox(account *dcrlibwallet.Account, popupMinSizeWidth int) *widget.Box {
	propertiesInfo := func(ID string, value string) *widget.Box {
		return widget.NewHBox(widgets.NewHSpacer(values.SpacerSize12), canvas.NewText(ID, values.TransactionInfoColor),
			layout.NewSpacer(),
			canvas.NewText(value, values.DefaultTextColor), widgets.NewHSpacer(values.SpacerSize12))
	}

	propertiesVBox := widget.NewVBox(
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(popupMinSizeWidth, 1)), canvas.NewLine(values.StrippedLineColor)),

		widget.NewHBox(widgets.NewHSpacer(values.SpacerSize12), canvas.NewText(values.Properties, values.TransactionInfoColor)),
		widgets.NewVSpacer(values.SpacerSize12),

		propertiesInfo(values.Properties, fmt.Sprintf("%d", account.Number)),
		widgets.NewVSpacer(values.SpacerSize14),

		propertiesInfo(values.HDPath, fmt.Sprintf(values.HDPathFormat, account.Number)),
		widgets.NewVSpacer(values.SpacerSize14),

		propertiesInfo(values.Keys, fmt.Sprintf("%d external, %d internal, %d imported",
			account.ExternalKeyCount, account.InternalKeyCount, account.ImportedKeyCount)),
		widgets.NewVSpacer(values.SpacerSize14),
	)

	return propertiesVBox
}

func centerObject(object fyne.CanvasObject) *widget.Box {
	return widget.NewVBox(layout.NewSpacer(), object, layout.NewSpacer())
}
