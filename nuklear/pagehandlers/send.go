package pagehandlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aarzilli/nucular"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet/txhelper"
	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/styles"
	"github.com/raedahgroup/godcr/nuklear/widgets"
)

const (
	addressFieldWidth = 300
	amountFieldWidth  = 150
	sectionSpacing    = 20
)

type SendHandler struct {
	wallet               walletcore.Wallet
	refreshWindowDisplay func()

	spendUnconfirmed      bool
	accountSelectorWidget *widgets.AccountSelector

	selectCustomInputs  bool
	isFetchingUTXOS     bool
	utxosFetchError     error
	utxos               []*utxoSelection
	utxosSelectionError string

	sendDestinations []*sendDestination

	isSubmitting bool
	sendErr      error
	successHash  string
}

type utxoSelection struct {
	selected bool
	utxo     *walletcore.UnspentOutput
}

type sendDestination struct {
	address    *nucular.TextEditor
	addressErr string
	amount     *nucular.TextEditor
	amountErr  string
}

func (handler *SendHandler) BeforeRender(wallet walletcore.Wallet, settings *config.Settings, refreshWindowDisplay func()) bool {
	handler.wallet = wallet
	handler.refreshWindowDisplay = refreshWindowDisplay

	handler.spendUnconfirmed = false // todo should use the value in settings
	handler.accountSelectorWidget = widgets.AccountSelectorWidget("From:", handler.spendUnconfirmed, true, wallet)

	handler.selectCustomInputs = false
	handler.isFetchingUTXOS = false
	handler.utxosFetchError = nil
	handler.utxos = nil
	handler.utxosSelectionError = ""

	handler.sendDestinations = nil
	handler.addSendDestination(false)

	handler.isSubmitting = false
	handler.sendErr = nil
	handler.successHash = ""

	return true
}

func (handler *SendHandler) Render(window *nucular.Window) {
	widgets.PageContentWindowDefaultPadding("Send", window, func(contentWindow *widgets.Window) {
		handler.accountSelectorWidget.Render(contentWindow)
		contentWindow.AddCheckbox("Spend Unconfirmed", &handler.spendUnconfirmed, func() {
			// reload account balance and refresh display
			handler.accountSelectorWidget = widgets.AccountSelectorWidget("From:", handler.spendUnconfirmed,
				true, handler.wallet)
			handler.accountSelectorWidget.Render(contentWindow)
			handler.refreshWindowDisplay()

			// use updated spend unconfirmed value to reload utxos list (if use custom inputs is checked)
			if handler.selectCustomInputs {
				go handler.fetchCustomInputs()
			}
		})

		/* CUSTOM INPUTS SECTION */
		contentWindow.AddHorizontalSpace(sectionSpacing) // add space before drawing the custom inputs section
		contentWindow.AddCheckbox("Select Custom Inputs", &handler.selectCustomInputs, handler.fetchCustomInputsCheck)
		if handler.isFetchingUTXOS {
			contentWindow.DisplayIsLoadingMessage()
		} else if handler.utxosFetchError != nil {
			contentWindow.DisplayErrorMessage("Unable to load inputs", handler.utxosFetchError)
		} else if handler.utxos != nil {
			utxosTable := widgets.NewTable()

			// add table header using nav font
			utxosTable.AddRowWithFont(styles.NavFont,
				widgets.NewLabelTableCell("", widgets.LeftCenterAlign),
				widgets.NewLabelTableCell("Address", widgets.LeftCenterAlign),
				widgets.NewLabelTableCell("Amount", widgets.LeftCenterAlign),
				widgets.NewLabelTableCell("Time", widgets.LeftCenterAlign),
				widgets.NewLabelTableCell("Confirmations", widgets.LeftCenterAlign),
			)

			for _, utxo := range handler.utxos {
				receiveTime := time.Unix(utxo.utxo.ReceiveTime, 0).Format(time.RFC1123)
				confirmations := strconv.Itoa(int(utxo.utxo.Confirmations))

				utxosTable.AddRow(
					widgets.NewCheckboxTableCell("", &utxo.selected, handler.calculateInputsPercentage),
					widgets.NewLabelTableCell(utxo.utxo.Address, widgets.LeftCenterAlign),
					widgets.NewLabelTableCell(utxo.utxo.Amount.String(), widgets.LeftCenterAlign),
					widgets.NewLabelTableCell(receiveTime, widgets.LeftCenterAlign),
					widgets.NewLabelTableCell(confirmations, widgets.LeftCenterAlign),
				)
			}
			utxosTable.Render(contentWindow)

			// error selecting utxo?
			if handler.utxosSelectionError != "" {
				contentWindow.DisplayMessage(handler.utxosSelectionError, styles.DecredOrangeColor)
			}
		}

		/* SEND DESTINATIONS SECTION */
		contentWindow.AddHorizontalSpace(sectionSpacing) // add space before drawing the send destinations section
		columnWidths := []int{addressFieldWidth, amountFieldWidth}
		// add headers
		contentWindow.AddLabelsWithWidths(columnWidths,
			widgets.NewLabelTableCell("Destination Address", widgets.LeftCenterAlign),
			widgets.NewLabelTableCell("Amount (DCR)", widgets.LeftCenterAlign),
		)
		// add destination fields
		for _, destination := range handler.sendDestinations {
			contentWindow.AddEditorsWithWidths(columnWidths, destination.address, destination.amount)

			// add errors if exist
			if destination.addressErr != "" || destination.amountErr != "" {
				errorLabels := make([]*widgets.LabelTableCell, 2)
				if destination.addressErr != "" {
					errorLabels[0] = widgets.NewColoredLabelTableCell(destination.addressErr, widgets.LeftCenterAlign,
						styles.DecredOrangeColor)
				}
				if destination.amountErr != "" {
					errorLabels[1] = widgets.NewColoredLabelTableCell(destination.amountErr, widgets.LeftCenterAlign,
						styles.DecredOrangeColor)
				}
				contentWindow.AddLabelsWithWidths(columnWidths, errorLabels...)
			}
		}
		// add/remove destination buttons
		contentWindow.Row(widgets.ButtonHeight).Static(
			contentWindow.ButtonWidth("Add another address"),
			contentWindow.ButtonWidth("Remove last address"),
		)
		contentWindow.AddButtonToCurrentRow("Add another address", func() {
			handler.addSendDestination(true)
		})
		if len(handler.sendDestinations) > 1 {
			contentWindow.AddButtonToCurrentRow("Remove last address", handler.removeLastSendInputPair)
		}

		contentWindow.AddHorizontalSpace(sectionSpacing) // add space before drawing submit button
		submitButtonText := "Submit"
		if handler.isSubmitting {
			submitButtonText = "Submitting"
		}
		contentWindow.AddButton(submitButtonText, func() {
			if !handler.isSubmitting && handler.validateForm() {
				handler.getPassphraseAndSubmit(contentWindow)
			}
		})

		// show result of last send op if exists
		if handler.sendErr != nil {
			contentWindow.DisplayErrorMessage("Send tx error", handler.sendErr)
		} else if handler.successHash != "" {
			successMessage := "The transaction was published successfully. Hash: " + handler.successHash
			contentWindow.AddWrappedLabelWithColor(successMessage, widgets.LeftCenterAlign, styles.DecredGreenColor)
		}
	})
}

// addSendDestination adds a address and amount input field pair on user click of
// the 'add another address' button. This function is called at least once in the lifetime of
// the send page
func (handler *SendHandler) addSendDestination(refreshWindowDisplay bool) {
	if handler.sendDestinations == nil {
		handler.sendDestinations = []*sendDestination{}
	}

	amountItemEditor := &nucular.TextEditor{}
	amountItemEditor.Flags = nucular.EditClipboard | nucular.EditSimple

	destinationAddressItemEditor := &nucular.TextEditor{}
	destinationAddressItemEditor.Flags = nucular.EditClipboard | nucular.EditSimple

	item := &sendDestination{
		amount:  amountItemEditor,
		address: destinationAddressItemEditor,
	}
	handler.sendDestinations = append(handler.sendDestinations, item)

	if refreshWindowDisplay {
		handler.refreshWindowDisplay()
	}
}

// removeLastSendInputPair removes the last address and amount input pair
// when the 'remove last address' button is clicked. The button is hidden when only one
// input pair exists on form
func (handler *SendHandler) removeLastSendInputPair() {
	handler.sendDestinations[len(handler.sendDestinations)-1] = nil
	handler.sendDestinations = handler.sendDestinations[:len(handler.sendDestinations)-1]
	handler.refreshWindowDisplay()
}

// fetchCustomInputsCheck is called everytime the 'select custom inputs' checkbox is checked or unchecked.
// for everytime the checkbox is checked, the fetchCustomInputs function is called
func (handler *SendHandler) fetchCustomInputsCheck() {
	if handler.selectCustomInputs {
		go handler.fetchCustomInputs()
		return
	}

	handler.isFetchingUTXOS = false
	handler.utxos = nil
	handler.refreshWindowDisplay()
}

// fetchCustomInputs fetches custom inputs to add to send transaction. It adds unconfirmed inputs if the
// 'spend unconfirmed' checkbox is checked
func (handler *SendHandler) fetchCustomInputs() {
	handler.isFetchingUTXOS = true
	handler.utxos = nil
	handler.refreshWindowDisplay()

	defer func() {
		handler.isFetchingUTXOS = false
		handler.refreshWindowDisplay()
	}()

	var requiredConfirmations int32 = walletcore.DefaultRequiredConfirmations
	if handler.spendUnconfirmed {
		requiredConfirmations = 0
	}

	accountNumber := handler.accountSelectorWidget.GetSelectedAccountNumber()
	utxos, err := handler.wallet.UnspentOutputs(accountNumber, 0, requiredConfirmations)
	if err != nil {
		handler.utxosFetchError = err
		return
	}

	handler.utxos = make([]*utxoSelection, len(utxos))
	for index, utxo := range utxos {
		utxoItem := &utxoSelection{
			selected: false,
			utxo:     utxo,
		}
		handler.utxos[index] = utxoItem
	}
}

// todo this should prolly be implemented as is done in the web interface
func (handler *SendHandler) calculateInputsPercentage() {

}

func (handler *SendHandler) validateForm() bool {
	isClean := true

	// clear errors before continuing
	handler.sendErr = nil
	handler.successHash = ""
	for _, destination := range handler.sendDestinations {
		destination.amountErr, destination.addressErr = "", ""
	}
	handler.utxosSelectionError = ""
	handler.refreshWindowDisplay()

	totalSendAmount := 0.0
	for _, destination := range handler.sendDestinations {
		address := string(destination.address.Buffer)
		if address == "" {
			destination.addressErr = "This address field is required"
		} else {
			isValid, err := handler.wallet.ValidateAddress(address)
			if err != nil {
				destination.addressErr = fmt.Sprintf("Error checking destination address: %s", err.Error())
			} else if !isValid {
				destination.addressErr = "Invalid address"
			}
		}

		amountStr := string(destination.amount.Buffer)
		if amountStr == "" {
			destination.amountErr = "This amount field is required"
		} else {
			amountFloat, err := strconv.ParseFloat(amountStr, 64)
			if err != nil {
				destination.amountErr = "This is not a valid number"
			} else if amountFloat < 1 {
				destination.amountErr = "Send amount must be greater than 0 DCR"
			} else {
				totalSendAmount += amountFloat
			}
		}

		if destination.addressErr != "" || destination.amountErr != "" {
			isClean = false
		}
	}

	// check if total selected utxo amount is not less than total send amount
	if handler.selectCustomInputs {
		totalSelectedUtxoAmount := 0.0
		for _, utxo := range handler.utxos {
			if utxo.selected {
				totalSelectedUtxoAmount += utxo.utxo.Amount.ToCoin()
			}
		}

		if totalSendAmount > totalSelectedUtxoAmount {
			handler.utxosSelectionError = fmt.Sprintf("Total send amount (%f DCR) is higher than the total input amount (%f DCR)",
				totalSendAmount, totalSelectedUtxoAmount)
			isClean = false
		}

	}

	handler.refreshWindowDisplay()
	return isClean
}

func (handler *SendHandler) getPassphraseAndSubmit(window *widgets.Window) {
	// clear success and/or error message
	handler.sendErr = nil
	handler.successHash = ""

	passphraseChan := make(chan string)
	widgets.NewPassphraseWidget().Get(window.Window, passphraseChan)

	go func() {
		passphrase := <-passphraseChan
		if passphrase != "" {
			handler.submit(passphrase, window)
		}
	}()
}

func (handler *SendHandler) submit(passphrase string, window *widgets.Window) {
	if handler.isSubmitting {
		return
	}

	handler.isSubmitting = true
	handler.refreshWindowDisplay()

	defer func() {
		handler.isSubmitting = false
		handler.refreshWindowDisplay()
	}()

	sendDestinations := make([]txhelper.TransactionDestination, len(handler.sendDestinations))
	for index := range handler.sendDestinations {
		amount, err := strconv.ParseFloat(string(handler.sendDestinations[index].amount.Buffer), 64)
		if err != nil {
			handler.sendErr = err
			return
		}

		sendDestinations[index] = txhelper.TransactionDestination{
			Address: string(handler.sendDestinations[index].address.Buffer),
			Amount:  amount,
		}
	}

	accountNumber := handler.accountSelectorWidget.GetSelectedAccountNumber()
	var requiredConfirmations int32 = walletcore.DefaultRequiredConfirmations
	if handler.spendUnconfirmed {
		requiredConfirmations = 0
	}

	if handler.selectCustomInputs {
		utxos, totalInputAmount := handler.getUTXOSAndSelectedAmount()

		changeAddress, err := handler.wallet.GenerateNewAddress(accountNumber)
		if err != nil {
			handler.sendErr = err
			return
		}

		changeAmount, err := txhelper.EstimateChange(len(handler.utxos), int64(totalInputAmount), sendDestinations, []string{changeAddress})
		if err != nil {
			handler.sendErr = err
			return
		}

		changeDestinations := []txhelper.TransactionDestination{{
			Amount:  dcrutil.Amount(changeAmount).ToCoin(),
			Address: changeAddress,
		}}

		handler.successHash, handler.sendErr = handler.wallet.SendFromUTXOs(accountNumber, requiredConfirmations, utxos,
			sendDestinations, changeDestinations, passphrase)
	} else {
		handler.successHash, handler.sendErr = handler.wallet.SendFromAccount(accountNumber, requiredConfirmations,
			sendDestinations, passphrase)
	}

	if handler.successHash != "" {
		handler.resetForm(window)
	} else {
		handler.refreshWindowDisplay()
	}
}

func (handler *SendHandler) getUTXOSAndSelectedAmount() (utxos []string, totalInputAmount dcrutil.Amount) {
	for _, utxo := range handler.utxos {
		if utxo.selected {
			totalInputAmount += utxo.utxo.Amount
			utxos = append(utxos, utxo.utxo.OutputKey)
		}
	}
	return
}

func (handler *SendHandler) resetForm(window *widgets.Window) {
	handler.spendUnconfirmed = false // todo should use the value in settings
	handler.accountSelectorWidget = widgets.AccountSelectorWidget("From:", handler.spendUnconfirmed,
		true, handler.wallet)
	handler.accountSelectorWidget.Render(window)

	handler.selectCustomInputs = false
	handler.isFetchingUTXOS = false
	handler.utxosFetchError = nil
	handler.utxos = nil
	handler.utxosSelectionError = ""

	handler.sendDestinations = nil
	handler.addSendDestination(false)

	handler.isSubmitting = false
	handler.sendErr = nil
	handler.successHash = ""

	handler.refreshWindowDisplay()
}
