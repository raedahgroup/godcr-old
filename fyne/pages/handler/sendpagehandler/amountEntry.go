package sendpagehandler

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/fyne/layouts"
	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (sendPage *SendPageObjects) initAmountEntryComponents() {
	amountLabel := canvas.NewText(values.Amount, values.DarkerBlueGrayTextColor)
	// amount entry accepts only floats
	amountEntryExpression, err := regexp.Compile(values.AmountRegExp)
	if err != nil {
		log.Println(err)
	}

	sendPage.amountEntry = widget.NewEntry()
	sendPage.amountEntry.SetPlaceHolder(fmt.Sprintf("%s %s", values.ZeroAmount, values.DCR))

	sendPage.amountEntryErrorLabel = widgets.NewTextWithSize("", values.ErrorColor, values.DefaultErrTextSize)
	sendPage.amountEntryErrorLabel.Hide()

	sendPage.amountEntry.OnChanged = func(value string) {
		if len(value) > 0 && !amountEntryExpression.MatchString(value) {
			if len(value) == 1 {
				sendPage.amountEntry.SetText("")
			} else {
				//fix issue with crash on paste here
				value = value[:sendPage.amountEntry.CursorColumn-1] + value[sendPage.amountEntry.CursorColumn:]
				//todo: using setText, cursor column count doesnt increase or reduce. Create an issue on this
				sendPage.amountEntry.CursorColumn--
				sendPage.amountEntry.SetText(value)
			}

			return
		}

		if sendPage.amountEntry.Focused() {
			sendPage.sendMax = false
		}

		sendPage.initTxDetails(value)
	}

	sendPage.initMaxButton()

	amountEntryComponents := widget.NewVBox(
		widget.NewHBox(amountLabel, layout.NewSpacer(), sendPage.SpendableLabel, widgets.NewHSpacer(values.SpacerSize20)),
		widgets.NewVSpacer(values.SpacerSize10),

		fyne.NewContainerWithLayout(
			layouts.NewPasswordLayout(
				fyne.NewSize(widget.NewLabel(values.MaxAmountAllowedInDCR).MinSize().Width+sendPage.maxButton.MinSize().Width, sendPage.amountEntry.MinSize().Height)),
			sendPage.amountEntry, sendPage.maxButton.Container),
		sendPage.amountEntryErrorLabel)

	sendPage.SendPageContents.Append(amountEntryComponents)
}

func (sendPage *SendPageObjects) initMaxButton() {
	sendPage.maxButton = widgets.NewButton(values.MaxButtonBackgroundColor, values.Max, func() {
		sendPage.sendMax = true
		transactionAuthor, _ := sendPage.initTxAuthorAndGetAmountInWalletAccount(0, "")
		if transactionAuthor == nil {
			return
		}

		maxAmount, err := transactionAuthor.EstimateMaxSendAmount()
		if err != nil && err.Error() == dcrlibwallet.ErrInsufficientBalance {
			sendPage.amountEntryErrorLabel.Text = values.NoFunds
			if !sendPage.MultiWallet.IsSynced() {
				sendPage.amountEntryErrorLabel.Text = values.NoFundsOrNotConnected
			}

			sendPage.SendPageContents.Refresh()
			sendPage.amountEntryErrorLabel.Show()
			sendPage.SendPageContents.Refresh()

			return
		}
		sendPage.SendPageContents.Refresh()
		sendPage.amountEntryErrorLabel.Hide()
		sendPage.SendPageContents.Refresh()

		sendPage.amountEntryErrorLabel.Text = ""
		sendPage.SendPageContents.Refresh()

		// entry widget doesn't perform the onchanged function if values are same
		sendPage.amountEntry.Text = ""
		sendPage.amountEntry.SetText(strconv.FormatFloat(maxAmount.DcrValue, 'f', -1, 64))
	})

	sendPage.maxButton.SetTextSize(values.TextSize10)
	sendPage.maxButton.SetTextStyle(fyne.TextStyle{Bold: true})
	sendPage.maxButton.SetMinSize(sendPage.maxButton.MinSize().Add(fyne.NewSize(8, 8)))

}

func (sendPage *SendPageObjects) initTxAuthorAndGetAmountInWalletAccount(amount float64, address string) (*dcrlibwallet.TxAuthor, int64) {
	transactionAuthor := sendPage.Sending.SelectedWallet.NewUnsignedTx(int32(*sendPage.Sending.SendingSelectedAccountID), dcrlibwallet.DefaultRequiredConfirmations)
	var err error
	if address == "" {
		address, err = sendPage.Sending.SelectedWallet.CurrentAddress(int32(*sendPage.Sending.SendingSelectedAccountID))
		if err != nil {
			sendPage.showErrorLabel(values.GettingAddressToSelfSendErr)
			return nil, 0
		}
	}

	accountBalance, err := sendPage.Sending.SelectedWallet.GetAccountBalance(int32(*sendPage.Sending.SendingSelectedAccountID), dcrlibwallet.DefaultRequiredConfirmations)
	if err != nil {
		sendPage.showErrorLabel(values.GettingAccountBalanceErr)
		return nil, 0
	}

	amountInAccount := accountBalance.Spendable

	transactionAuthor.AddSendDestination(address, dcrlibwallet.AmountAtom(amount), sendPage.sendMax)

	return transactionAuthor, amountInAccount
}

func (sendPage *SendPageObjects) initTxDetails(amountInString string) {
	if numbers := strings.Split(amountInString, "."); len(numbers) == 2 && len(numbers[1]) > 8 {
		sendPage.showErrorLabel(values.AmountDecimalPlaceErr)
		return
	}

	amountInFloat, err := strconv.ParseFloat(amountInString, 64)
	if err != nil && amountInString != "" {
		sendPage.showErrorLabel(values.ParseFloatErr)
		return
	}

	if amountInFloat == 0.0 || !sendPage.destinationAddressErrorLabel.Hidden {
		setLabelText(values.NilAmount, sendPage.transactionFeeLabel, sendPage.totalCostLabel, sendPage.balanceAfterSendLabel)
		sendPage.transactionSizeLabel.Text = values.ZeroByte
		sendPage.SendPageContents.Refresh()

		setLabelColor(values.NilAmountColor, sendPage.transactionFeeLabel, sendPage.totalCostLabel, sendPage.balanceAfterSendLabel)
		sendPage.SendPageContents.Refresh()

		sendPage.nextButton.Disable()
		sendPage.SendPageContents.Refresh()

		return
	}

	transactionAuthor, amountInAccount := sendPage.initTxAuthorAndGetAmountInWalletAccount(amountInFloat, "")
	if transactionAuthor == nil {
		return
	}

	feeAndSize, err := transactionAuthor.EstimateFeeAndSize()
	if err != nil {
		sendPage.SendPageContents.Refresh()
		if err.Error() == dcrlibwallet.ErrInsufficientBalance {
			sendPage.amountEntryErrorLabel.Text = values.InsufficientBalanceErr
			sendPage.amountEntryErrorLabel.Show()
		} else {
			sendPage.showErrorLabel(values.TransactionFeeSizeErr)
		}
		sendPage.SendPageContents.Refresh()

		setLabelText(values.NilAmount, sendPage.transactionFeeLabel, sendPage.totalCostLabel, sendPage.balanceAfterSendLabel)
		sendPage.transactionSizeLabel.Text = values.ZeroByte
		setLabelColor(values.NilAmountColor, sendPage.transactionFeeLabel, sendPage.totalCostLabel, sendPage.balanceAfterSendLabel)

		sendPage.SendPageContents.Refresh()
		sendPage.nextButton.Disable()
		sendPage.SendPageContents.Refresh()
		return
	}

	totalCostInAtom := feeAndSize.Fee.AtomValue + dcrlibwallet.AmountAtom(amountInFloat)
	balanceAfterSendInAtom := amountInAccount - totalCostInAtom
	sendPage.SendPageContents.Refresh()
	if !sendPage.amountEntryErrorLabel.Hidden {
		sendPage.amountEntryErrorLabel.Hide()
	}
	sendPage.SendPageContents.Refresh()

	sendPage.SendPageContents.Refresh()
	sendPage.transactionFeeLabel.Text = strconv.FormatFloat(feeAndSize.Fee.DcrValue, 'f', -1, 64)
	sendPage.totalCostLabel.Text = strconv.FormatFloat(dcrlibwallet.AmountCoin(totalCostInAtom), 'f', -1, 64)
	sendPage.transactionSizeLabel.Text = fmt.Sprintf("%d %s", feeAndSize.EstimatedSignedSize, values.Bytes)
	sendPage.balanceAfterSendLabel.Text = strconv.FormatFloat(dcrlibwallet.AmountCoin(balanceAfterSendInAtom), 'f', -1, 64)

	setLabelColor(values.DefaultTextColor, sendPage.transactionFeeLabel, sendPage.totalCostLabel, sendPage.balanceAfterSendLabel)
	sendPage.SendPageContents.Refresh()

	if sendPage.destinationAddressEntry.Text != "" && sendPage.destinationAddressErrorLabel.Hidden || sendPage.destinationAddressEntry.Hidden {
		sendPage.nextButton.Enable()
	} else {
		sendPage.nextButton.Disable()
	}

	sendPage.SendPageContents.Refresh()
}
