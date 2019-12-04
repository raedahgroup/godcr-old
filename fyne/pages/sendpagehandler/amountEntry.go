package sendpagehandler

import (
	"fmt"
	"image/color"
	"log"
	"regexp"
	"strconv"
	"strings"

	"fyne.io/fyne/layout"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/layouts"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func AmountEntryComponents(errorLabel *widgets.Button, showErrorLabel func(string), temporaryAddress *string, amountInAccount *float64,
	transactionAuthor *dcrlibwallet.TxAuthor, transactionFeeLabel, totalCostLabel, balanceAfterSendLabel, transactionSizeLabel *widget.Label,
	destinationAddressEntryText *string, isDestinationAddressEntryHidden, isDestinationAddressErrorLabelHidden *bool,
	contents *widget.Box, nextButton *widgets.Button, spendableLabel *canvas.Text, multiWallet *dcrlibwallet.MultiWallet) (*fyne.Container, *widget.Entry, bool) {

	amountLabel := canvas.NewText("Amount", color.RGBA{61, 88, 115, 255})
	amountLabel.TextStyle.Bold = true

	// amount entry accepts only floats
	amountEntryExpression, err := regexp.Compile("^\\d*\\.?\\d*$")
	if err != nil {
		log.Println(err)
	}

	amountEntry := widget.NewEntry()
	amountEntry.SetPlaceHolder("0 DCR")

	amountErrorLabel := canvas.NewText("", color.RGBA{237, 109, 71, 255})
	amountErrorLabel.TextSize = 12
	amountErrorLabel.Hide()

	amountEntry.OnChanged = func(value string) {
		if len(value) > 0 && !amountEntryExpression.MatchString(value) {
			if len(value) == 1 {
				amountEntry.SetText("")
			} else {
				//fix issue with crash on paste here
				value = value[:amountEntry.CursorColumn-1] + value[amountEntry.CursorColumn:]
				//todo: using setText, cursor column count doesnt increase or reduce. Create an issue on this
				amountEntry.CursorColumn--
				amountEntry.SetText(value)
			}
			return
		}

		if numbers := strings.Split(value, "."); len(numbers) == 2 {
			if len(numbers[1]) > 8 {
				showErrorLabel("Amount has more than 8 decimal places.")
				return
			}
		}

		amountInFloat, err := strconv.ParseFloat(value, 64)
		if err != nil && value != "" {
			showErrorLabel("Could not parse float")
			return
		}

		if amountInFloat == 0.0 {
			transactionFeeLabel.SetText("- DCR")
			totalCostLabel.SetText("- DCR")
			balanceAfterSendLabel.SetText("- DCR")
			transactionSizeLabel.SetText("0 bytes")
			amountErrorLabel.Hide()

			nextButton.Disable()
			widgets.Refresher(transactionFeeLabel, totalCostLabel, balanceAfterSendLabel, transactionSizeLabel)
			// paintedtransactionInfoform.Refresh()
			//contents.Refresh()
			return
		}

		transactionAuthor.UpdateSendDestination(0, *temporaryAddress, dcrlibwallet.AmountAtom(amountInFloat), false)
		feeAndSize, err := transactionAuthor.EstimateFeeAndSize()
		if err != nil {
			if err.Error() == dcrlibwallet.ErrInsufficientBalance {
				amountErrorLabel.Text = "Insufficient balance"
				amountErrorLabel.Show()
			} else {
				showErrorLabel("Could not retrieve transaction fee and size")
				log.Println(fmt.Sprintf("could not retrieve transaction fee and size %s", err.Error()))
			}

			transactionFeeLabel.SetText(fmt.Sprintf("- DCR"))
			totalCostLabel.SetText("- DCR")
			balanceAfterSendLabel.SetText("- DCR")
			transactionSizeLabel.SetText("0 bytes")

			nextButton.Disable()
			//paintedtransactionInfoform.Refresh()
			widgets.Refresher(transactionFeeLabel, totalCostLabel, balanceAfterSendLabel, transactionSizeLabel) //, costAndBalanceAfterSendBox)
			contents.Refresh()
			return
		}

		if !amountErrorLabel.Hidden {
			amountErrorLabel.Hide()
		}

		transactionFeeLabel.SetText(fmt.Sprintf("%f DCR", feeAndSize.Fee.DcrValue))
		totalCostLabel.SetText(fmt.Sprintf("%f DCR", feeAndSize.Fee.DcrValue+amountInFloat))
		balanceAfterSendLabel.SetText(fmt.Sprintf("%f DCR", *amountInAccount-(feeAndSize.Fee.DcrValue+amountInFloat)))
		transactionSizeLabel.SetText(fmt.Sprintf("%d bytes", feeAndSize.EstimatedSignedSize))

		if *destinationAddressEntryText != "" && *isDestinationAddressErrorLabelHidden || *isDestinationAddressEntryHidden {
			nextButton.Enable()
		} else {
			nextButton.Disable()
		}

		errorLabel.Container.Hide()
		//paintedtransactionInfoform.Refresh()
		widgets.Refresher(transactionFeeLabel, totalCostLabel, balanceAfterSendLabel, transactionSizeLabel) //, costAndBalanceAfterSendBox)
		contents.Refresh()
	}

	maxButton := maxButton(temporaryAddress, amountErrorLabel, amountEntry, transactionAuthor, multiWallet, contents)

	amountEntryComponents := widget.NewVBox(
		widget.NewHBox(amountLabel, layout.NewSpacer(), spendableLabel),
		widgets.NewVSpacer(10),

		fyne.NewContainerWithLayout(
			layouts.NewPasswordLayout(
				fyne.NewSize(widget.NewLabel("12345678.12345678").MinSize().Width+maxButton.MinSize().Width, amountEntry.MinSize().Height)),
			amountEntry, maxButton.Container),

		amountErrorLabel)

	return fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(fyne.Max(amountEntryComponents.MinSize().Width, 312), amountEntryComponents.MinSize().Height)), amountEntryComponents), amountEntry, amountErrorLabel.Hidden
}

func maxButton(temporaryAddress *string, amountErrorLabel *canvas.Text, amountEntry *widget.Entry,
	transactionAuthor *dcrlibwallet.TxAuthor, multiWallet *dcrlibwallet.MultiWallet, contents *widget.Box) *widgets.Button {

	maxButton := widgets.NewButton(color.RGBA{61, 88, 115, 255}, "MAX", func() {
		transactionAuthor.UpdateSendDestination(0, *temporaryAddress, 0, true)

		maxAmount, err := transactionAuthor.EstimateMaxSendAmount()
		if err != nil {
			if err.Error() == dcrlibwallet.ErrInsufficientBalance {
				amountErrorLabel.Text = "Not enough funds"
				if !multiWallet.IsSynced() {
					amountErrorLabel.Text = "Not enough funds (or not connected)."
				}

				amountErrorLabel.Show()
				contents.Refresh()
				return
			}
		}

		amountErrorLabel.Hide()
		contents.Refresh()
		amountEntry.SetText(fmt.Sprintf("%f", maxAmount.DcrValue-0.000012))
	})

	maxButton.SetTextSize(9)
	maxButton.SetMinSize(maxButton.MinSize().Add(fyne.NewSize(8, 8)))

	return maxButton
}
