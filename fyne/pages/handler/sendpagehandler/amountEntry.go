package sendpagehandler

import (
	"fmt"
	"image/color"
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
	amountLabel := canvas.NewText(values.Amount, color.RGBA{61, 88, 115, 255})
	amountLabel.TextStyle.Bold = true

	// amount entry accepts only floats
	amountEntryExpression, err := regexp.Compile(values.AmountRegExp)
	if err != nil {
		log.Println(err)
	}

	sendPage.amountEntry = widget.NewEntry()
	sendPage.amountEntry.SetPlaceHolder(values.ZeroAmount)

	sendPage.amountEntryErrorLabel = widgets.NewTextWithSize("", color.RGBA{237, 109, 71, 255}, values.DefaultErrTextSize)
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

		if numbers := strings.Split(value, "."); len(numbers) == 2 {
			if len(numbers[1]) > 8 {
				sendPage.showErrorLabel(values.AmountDecimalPlaceErr)
				return
			}
		}

		amountInFloat, err := strconv.ParseFloat(value, 64)
		if err != nil && value != "" {
			sendPage.showErrorLabel(values.ParseFloatErr)
			return
		}

		if amountInFloat == 0.0 || !sendPage.destinationAddressErrorLabel.Hidden {
			setLabelText(values.NilAmount, sendPage.transactionFeeLabel, sendPage.totalCostLabel, sendPage.balanceAfterSendLabel)
			sendPage.transactionSizeLabel.SetText(values.ZeroByte)

			sendPage.nextButton.Disable()
			widgets.Refresher(sendPage.transactionFeeLabel, sendPage.totalCostLabel, sendPage.balanceAfterSendLabel, sendPage.transactionSizeLabel)
			return
		}

		transactionAuthor, amountInAccount := sendPage.initTxAuthorAndGetAmountInWalletAccount(amountInFloat, "")
		if transactionAuthor == nil {
			return
		}

		feeAndSize, err := transactionAuthor.EstimateFeeAndSize()
		if err != nil {
			if err.Error() == dcrlibwallet.ErrInsufficientBalance {
				sendPage.amountEntryErrorLabel.Text = values.InsufficientBalanceErr
				sendPage.amountEntryErrorLabel.Show()
			} else {
				sendPage.showErrorLabel(values.TransactionFeeSizeErr)
			}
			sendPage.SendPageContents.Refresh()

			setLabelText(values.NilAmount, sendPage.transactionFeeLabel, sendPage.totalCostLabel, sendPage.balanceAfterSendLabel)
			sendPage.transactionSizeLabel.SetText(values.ZeroByte)
			widgets.Refresher(sendPage.transactionFeeLabel, sendPage.totalCostLabel, sendPage.balanceAfterSendLabel, sendPage.transactionSizeLabel)

			sendPage.nextButton.Disable()

			return
		}

		if !sendPage.amountEntryErrorLabel.Hidden {
			sendPage.amountEntryErrorLabel.Hide()
		}

		totalCostInAtom := feeAndSize.Fee.AtomValue + dcrlibwallet.AmountAtom(amountInFloat)
		balanceAfterSendInAtom := amountInAccount - totalCostInAtom

		sendPage.errorLabel.Container.Hide()

		if sendPage.destinationAddressEntry.Text != "" && sendPage.destinationAddressErrorLabel.Hidden || sendPage.destinationAddressEntry.Hidden {
			sendPage.nextButton.Enable()
		} else {
			sendPage.nextButton.Disable()
		}

		sendPage.transactionFeeLabel.SetText(fmt.Sprintf("%s %s", strconv.FormatFloat(feeAndSize.Fee.DcrValue, 'f', -1, 64), values.DCR))
		sendPage.totalCostLabel.SetText(fmt.Sprintf("%s %s", strconv.FormatFloat(dcrlibwallet.AmountCoin(totalCostInAtom), 'f', -1, 64), values.DCR))
		sendPage.transactionSizeLabel.SetText(fmt.Sprintf("%d %s", feeAndSize.EstimatedSignedSize, values.Bytes))
		sendPage.balanceAfterSendLabel.SetText(fmt.Sprintf("%s %s", strconv.FormatFloat(dcrlibwallet.AmountCoin(balanceAfterSendInAtom), 'f', -1, 64), values.DCR))
	}

	maxButton := sendPage.maxButton()

	amountEntryComponents := widget.NewVBox(
		widget.NewHBox(amountLabel, layout.NewSpacer(), sendPage.SpendableLabel, widgets.NewHSpacer(20)),
		widgets.NewVSpacer(10),

		fyne.NewContainerWithLayout(
			layouts.NewPasswordLayout(
				fyne.NewSize(widget.NewLabel(values.MaxAmountAllowedInDCR).MinSize().Width+maxButton.MinSize().Width, sendPage.amountEntry.MinSize().Height)),
			sendPage.amountEntry, maxButton.Container),
		sendPage.amountEntryErrorLabel)

	sendPage.SendPageContents.Append(amountEntryComponents)
}

func (sendPage *SendPageObjects) maxButton() *widgets.Button {
	maxButton := widgets.NewButton(color.RGBA{61, 88, 115, 255}, values.Max, func() {
		transactionAuthor, _ := sendPage.initTxAuthorAndGetAmountInWalletAccount(0, "")
		if transactionAuthor == nil {
			return
		}

		maxAmount, err := transactionAuthor.EstimateMaxSendAmount()
		if err != nil {
			if err.Error() == dcrlibwallet.ErrInsufficientBalance {
				sendPage.amountEntryErrorLabel.Text = values.NoFunds
				if !sendPage.MultiWallet.IsSynced() {
					sendPage.amountEntryErrorLabel.Text = values.NoFundsOrNotConnected
				}

				sendPage.amountEntryErrorLabel.Show()
				sendPage.SendPageContents.Refresh()
			}
			return
		}

		sendPage.amountEntryErrorLabel.Hide()
		sendPage.SendPageContents.Refresh()

		sendPage.amountEntry.SetText(strconv.FormatFloat(maxAmount.DcrValue, 'f', -1, 64))
	})

	maxButton.SetTextSize(9)
	maxButton.SetMinSize(maxButton.MinSize().Add(fyne.NewSize(8, 8)))

	return maxButton
}

func (sendPage *SendPageObjects) initTxAuthorAndGetAmountInWalletAccount(amount float64, address string) (*dcrlibwallet.TxAuthor, int64) {
	accNo, err := sendPage.Sending.SelectedWallet.AccountNumber(sendPage.Sending.SelectedAccountLabel.Text)
	if err != nil {
		sendPage.showErrorLabel(values.AccountNumberErr)
		return nil, 0
	}

	transactionAuthor := sendPage.Sending.SelectedWallet.NewUnsignedTx(int32(accNo), dcrlibwallet.DefaultRequiredConfirmations)

	if address == "" {
		address, err = sendPage.Sending.SelectedWallet.CurrentAddress(int32(accNo))
		if err != nil {
			sendPage.showErrorLabel(values.GettingAddressToSelfSendErr)
			return nil, 0
		}
	}

	accountBalance, err := sendPage.Sending.SelectedWallet.GetAccountBalance(int32(accNo), dcrlibwallet.DefaultRequiredConfirmations)
	if err != nil {
		sendPage.showErrorLabel(values.GettingAccountBalanceErr)
		return nil, 0
	}

	amountInAccount := accountBalance.Spendable

	transactionAuthor.AddSendDestination(address, 0, true)
	maxAmnt, err := transactionAuthor.EstimateMaxSendAmount()
	if err != nil {
		sendPage.showErrorLabel(values.MaxAmntErr)
		return nil, 0
	}

	if amount == maxAmnt.DcrValue || amount == 0 {
		transactionAuthor.UpdateSendDestination(0, address, dcrlibwallet.AmountAtom(amount), true)
	} else {
		transactionAuthor.UpdateSendDestination(0, address, dcrlibwallet.AmountAtom(amount), false)
	}

	return transactionAuthor, amountInAccount
}
