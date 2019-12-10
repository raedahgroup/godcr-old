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
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (sendPage *SendPageObjects) initAmountEntryComponents() {

	amountLabel := canvas.NewText("Amount", color.RGBA{61, 88, 115, 255})
	amountLabel.TextStyle.Bold = true

	// amount entry accepts only floats
	amountEntryExpression, err := regexp.Compile("^\\d*\\.?\\d*$")
	if err != nil {
		log.Println(err)
	}

	sendPage.amountEntry = widget.NewEntry()
	sendPage.amountEntry.SetPlaceHolder(ZeroAmount)

	sendPage.amountEntryErrorLabel = canvas.NewText("", color.RGBA{237, 109, 71, 255})
	sendPage.amountEntryErrorLabel.TextSize = 12
	sendPage.amountEntryErrorLabel.Hide()

	var container *fyne.Container

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
				sendPage.showErrorLabel("Amount has more than 8 decimal places.")
				return
			}
		}

		amountInFloat, err := strconv.ParseFloat(value, 64)
		if err != nil && value != "" {
			sendPage.showErrorLabel("Could not parse float")
			return
		}

		if amountInFloat == 0.0 {
			setLabelText(NilAmount, sendPage.transactionFeeLabel, sendPage.totalCostLabel, sendPage.balanceAfterSendLabel)
			sendPage.transactionSizeLabel.SetText(ZeroByte)

			sendPage.nextButton.Disable()
			widgets.Refresher(sendPage.transactionFeeLabel, sendPage.totalCostLabel, sendPage.balanceAfterSendLabel, sendPage.transactionSizeLabel)
			container.Refresh()
			sendPage.SendPageContents.Refresh()
			return
		}

		transactionAuthor, amountInAccount := sendPage.initTxAuthorAndGetAmountInWalletAccount(amountInFloat, "")
		if transactionAuthor == nil {
			sendPage.showErrorLabel("unable to initialize TxAuthor")

			return
		}

		feeAndSize, err := transactionAuthor.EstimateFeeAndSize()
		if err != nil {
			if err.Error() == dcrlibwallet.ErrInsufficientBalance {
				sendPage.amountEntryErrorLabel.Text = errInsufficientBalance
				sendPage.amountEntryErrorLabel.Show()
			} else {
				sendPage.showErrorLabel("Could not retrieve transaction fee and size")
				log.Println(fmt.Sprintf("could not retrieve transaction fee and size %s", err.Error()))
			}

			setLabelText(NilAmount, sendPage.transactionFeeLabel, sendPage.totalCostLabel, sendPage.balanceAfterSendLabel)
			sendPage.transactionSizeLabel.SetText(ZeroByte)

			sendPage.nextButton.Disable()
			widgets.Refresher(sendPage.transactionFeeLabel, sendPage.totalCostLabel, sendPage.balanceAfterSendLabel, sendPage.transactionSizeLabel)
			container.Refresh()
			sendPage.SendPageContents.Refresh()

			return
		}

		if !sendPage.amountEntryErrorLabel.Hidden {
			sendPage.amountEntryErrorLabel.Hide()
		}

		sendPage.transactionFeeLabel.SetText(fmt.Sprintf("%f DCR", feeAndSize.Fee.DcrValue))
		sendPage.totalCostLabel.SetText(fmt.Sprintf("%f DCR", feeAndSize.Fee.DcrValue+amountInFloat))
		sendPage.balanceAfterSendLabel.SetText(fmt.Sprintf("%f DCR", amountInAccount-(feeAndSize.Fee.DcrValue+amountInFloat)))
		sendPage.transactionSizeLabel.SetText(fmt.Sprintf("%d bytes", feeAndSize.EstimatedSignedSize))

		if sendPage.destinationAddressEntry.Text != "" && sendPage.destinationAddressErrorLabel.Hidden || sendPage.destinationAddressEntry.Hidden {
			sendPage.nextButton.Enable()
		} else {
			sendPage.nextButton.Disable()
		}

		sendPage.errorLabel.Container.Hide()
		widgets.Refresher(sendPage.transactionFeeLabel, sendPage.totalCostLabel, sendPage.balanceAfterSendLabel, sendPage.transactionSizeLabel)
		container.Refresh()
		sendPage.SendPageContents.Refresh()
	}

	maxButton := sendPage.maxButton()

	amountEntryComponents := widget.NewVBox(
		widget.NewHBox(amountLabel, layout.NewSpacer(), sendPage.SpendableLabel),
		widgets.NewVSpacer(10),

		fyne.NewContainerWithLayout(
			layouts.NewPasswordLayout(
				fyne.NewSize(widget.NewLabel(maxAmountInDCR).MinSize().Width+maxButton.MinSize().Width, sendPage.amountEntry.MinSize().Height)),
			sendPage.amountEntry, maxButton.Container),
		sendPage.amountEntryErrorLabel)

	container = fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(
		fyne.Max(amountEntryComponents.MinSize().Width, 312), amountEntryComponents.MinSize().Height)),
		amountEntryComponents)

	sendPage.SendPageContents.Append(container)
}

func (sendPage *SendPageObjects) maxButton() *widgets.Button {
	maxButton := widgets.NewButton(color.RGBA{61, 88, 115, 255}, "MAX", func() {
		transactionAuthor, _ := sendPage.initTxAuthorAndGetAmountInWalletAccount(0, "")
		maxAmount, err := transactionAuthor.EstimateMaxSendAmount()
		if err != nil {
			if err.Error() == dcrlibwallet.ErrInsufficientBalance {
				sendPage.amountEntryErrorLabel.Text = noFunds
				if !sendPage.MultiWallet.IsSynced() {
					sendPage.amountEntryErrorLabel.Text = noFundsOrNotConnected
				}

				sendPage.amountEntryErrorLabel.Show()
				sendPage.SendPageContents.Refresh()
				return
			}
		}

		sendPage.amountEntryErrorLabel.Hide()
		sendPage.SendPageContents.Refresh()

		sendPage.amountEntry.SetText(fmt.Sprintf("%f", maxAmount.DcrValue-0.000012))
	})

	maxButton.SetTextSize(9)
	maxButton.SetMinSize(maxButton.MinSize().Add(fyne.NewSize(8, 8)))

	return maxButton
}

func (sendPage *SendPageObjects) initTxAuthorAndGetAmountInWalletAccount(amount float64, address string) (*dcrlibwallet.TxAuthor, float64) {
	accNo, err := sendPage.Sending.selectedWallet.AccountNumber(sendPage.Sending.SelectedAccountLabel.Text)
	if err != nil {
		sendPage.showErrorLabel("could not retrieve account number")
		return nil, 0
	}

	transactionAuthor := sendPage.Sending.selectedWallet.NewUnsignedTx(int32(accNo), dcrlibwallet.DefaultRequiredConfirmations)

	if address == "" {
		address, err = sendPage.Sending.selectedWallet.CurrentAddress(int32(accNo))
		if err != nil {
			sendPage.showErrorLabel("could not get temporary account for send destination")
			return nil, 0
		}
	}

	accountBalance, err := sendPage.Sending.selectedWallet.GetAccountBalance(int32(accNo), dcrlibwallet.DefaultRequiredConfirmations)
	if err != nil {
		sendPage.showErrorLabel("could not get account balace for send destination")
		return nil, 0
	}

	amountInAccount := dcrlibwallet.AmountCoin(accountBalance.Spendable)

	var sendMax bool
	if amount == 0 {
		sendMax = true
	}
	transactionAuthor.AddSendDestination(address, dcrlibwallet.AmountAtom(amount), sendMax)

	return transactionAuthor, amountInAccount
}
