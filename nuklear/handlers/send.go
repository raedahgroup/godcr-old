package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aarzilli/nucular"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet/txhelper"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/handlers/widgets"
	"github.com/raedahgroup/godcr/nuklear/helpers"
)

type utxoSelection struct {
	selected bool
	utxo     *walletcore.UnspentOutput
}

type SendDetailInputPair struct {
	addressErr         string
	amountErr          string
	destinationAddress nucular.TextEditor
	amount             nucular.TextEditor
}

type SendHandler struct {
	err            error
	fetchUTXOError error
	isRendering    bool

	utxos           []*utxoSelection
	isFetchingUTXOS bool

	accountNumbers       []uint32
	accountOverviews     []string
	selectedAccountIndex int

	sendDetailInputPairs []*SendDetailInputPair

	spendUnconfirmed   bool
	selectCustomInputs bool

	isSubmitting bool

	successHash string
}

func (handler *SendHandler) BeforeRender() {
	handler.err = nil
	handler.fetchUTXOError = nil
	handler.utxos = nil
	handler.isFetchingUTXOS = false
	handler.accountNumbers = nil
	handler.accountOverviews = nil
	handler.sendDetailInputPairs = nil
	handler.spendUnconfirmed = false
	handler.selectCustomInputs = false
	handler.isSubmitting = false
	handler.isRendering = false
	handler.successHash = ""
}

func (handler *SendHandler) Render(window *nucular.Window, wallet walletcore.Wallet) {
	if !handler.isRendering {
		handler.isRendering = true
		handler.fetchAccounts(wallet)
		handler.addSendInputPair(window.Master(), true)
	}

	masterWindow := window.Master()

	// draw page
	if pageWindow := helpers.NewWindow("Send Page", window, 0); pageWindow != nil {
		pageWindow.DrawHeader("Send")

		// content window
		if contentWindow := pageWindow.ContentWindow("Send Form"); contentWindow != nil {
			helpers.SetFont(window, helpers.PageContentFont)

			if handler.err != nil {
				contentWindow.Row(10).Dynamic(1)
				contentWindow.LabelColored(handler.err.Error(), "LC", helpers.DangerColor)
			}

			if handler.successHash != "" {
				contentWindow.Row(10).Dynamic(1)
				contentWindow.LabelColored("The transaction was published successfully. Hash: "+handler.successHash, "LC", helpers.SuccessColor)
			}

			contentWindow.Row(10).Dynamic(2)
			contentWindow.Label("Source Account", "LC")

			contentWindow.Row(25).Dynamic(1)
			handler.selectedAccountIndex = contentWindow.ComboSimple(handler.accountOverviews, handler.selectedAccountIndex, 25)

			contentWindow.Row(15).Dynamic(2)
			if contentWindow.CheckboxText("Spend Unconfirmed", &handler.spendUnconfirmed) {
				handler.fetchCustomInputsCheck(wallet, masterWindow)
			}

			contentWindow.Spacing(2)

			for _, destInputPair := range handler.sendDetailInputPairs {
				contentWindow.Row(10).Dynamic(2)
				contentWindow.Label("Destination Address", "LC")
				contentWindow.Label("Amount (DCR)", "LC")

				contentWindow.Row(25).Dynamic(2)
				destInputPair.destinationAddress.Edit(contentWindow.Window)
				destInputPair.amount.Edit(contentWindow.Window)

				if destInputPair.addressErr != "" || destInputPair.amountErr != "" {
					contentWindow.Row(10).Dynamic(2)
					contentWindow.LabelColored(destInputPair.addressErr, "LC", helpers.DangerColor)
					contentWindow.LabelColored(destInputPair.amountErr, "LC", helpers.DangerColor)
				}
			}

			contentWindow.Row(25).Dynamic(2)
			if contentWindow.ButtonText("Add another address") {
				handler.addSendInputPair(window.Master(), true)
			}

			if len(handler.sendDetailInputPairs) > 1 && contentWindow.ButtonText("Remove last address") {
				handler.removeLastSendInputPair(window.Master())
			}

			contentWindow.Spacing(2)

			contentWindow.Row(15).Dynamic(2)
			if contentWindow.CheckboxText("Select custom inputs", &handler.selectCustomInputs) {
				if handler.validate(window, wallet) {
					handler.fetchCustomInputsCheck(wallet, masterWindow)
				} else {
					handler.selectCustomInputs = false
				}
			}

			if handler.isFetchingUTXOS {
				widgets.ShowIsFetching(contentWindow)
			} else if handler.fetchUTXOError != nil {
				contentWindow.Row(10).Dynamic(1)
				contentWindow.LabelColored(handler.fetchUTXOError.Error(), "LC", helpers.DangerColor)
			} else if handler.utxos != nil {
				contentWindow.Row(20).Ratio(0.1, 0.3, 0.2, 0.2, 0.2)
				contentWindow.Label("", "LC")
				contentWindow.Label("Address", "LC")
				contentWindow.Label("Amount", "LC")
				contentWindow.Label("Time", "LC")
				contentWindow.Label("Confirmations", "LC")

				for _, utxo := range handler.utxos {
					amountStr := utxo.utxo.Amount.String()
					receiveTime := time.Unix(utxo.utxo.ReceiveTime, 0).Format(time.RFC1123)
					confirmations := strconv.Itoa(int(utxo.utxo.Confirmations))

					contentWindow.Row(20).Ratio(0.04, 0.36, 0.2, 0.2, 0.2)
					contentWindow.CheckboxText("", &utxo.selected)
					contentWindow.Label(utxo.utxo.Address, "LC")
					contentWindow.Label(amountStr, "LC")
					contentWindow.Label(receiveTime, "LC")
					contentWindow.Label(confirmations, "LC")
				}
			}

			contentWindow.Spacing(2)

			submitButtonText := "Submit"
			if handler.isSubmitting {
				submitButtonText = "Submitting"
			}

			contentWindow.Row(25).Dynamic(2)
			if contentWindow.ButtonText(submitButtonText) {
				if handler.validate(window, wallet) {
					handler.getPassphraseAndSubmit(window, wallet)
				}
			}
			contentWindow.End()
		}
		pageWindow.End()
	}
}

// fetch accounts to display in select source account dropdown
func (handler *SendHandler) fetchAccounts(wallet walletcore.Wallet) {
	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		handler.err = err
		return
	}

	numAccounts := len(accounts)
	handler.accountNumbers = make([]uint32, numAccounts)
	handler.accountOverviews = make([]string, numAccounts)

	for index, account := range accounts {
		handler.accountOverviews[index] = fmt.Sprintf("%s - Total %s (Spendable %s)", account.Name, account.Balance.Total.String(), account.Balance.Spendable.String())
		handler.accountNumbers[index] = account.Number
	}
	handler.selectedAccountIndex = 0
}

// addSendInputPair adds a destinationAddress and amount input field pair on user click of
// the 'add another address' button. This function is called at least once in the lifetime of
// the send page
func (handler *SendHandler) addSendInputPair(window nucular.MasterWindow, updateWindow bool) {
	if handler.sendDetailInputPairs == nil {
		handler.sendDetailInputPairs = []*SendDetailInputPair{}
	}

	amountItemEditor := nucular.TextEditor{}
	amountItemEditor.Flags = nucular.EditClipboard | nucular.EditSimple

	destinationAddressItemEditor := nucular.TextEditor{}
	destinationAddressItemEditor.Flags = nucular.EditClipboard | nucular.EditSimple

	item := &SendDetailInputPair{
		amount:             amountItemEditor,
		destinationAddress: destinationAddressItemEditor,
	}
	handler.sendDetailInputPairs = append(handler.sendDetailInputPairs, item)

	if updateWindow {
		window.Changed()
	}
}

// removeLastSendInputPair removes the last destinationAddress and amount input pair
// when the 'remove last address' button is clicked. The button is hidden when only one
// input pair exists on form
func (handler *SendHandler) removeLastSendInputPair(window nucular.MasterWindow) {
	handler.sendDetailInputPairs[len(handler.sendDetailInputPairs)-1] = nil
	handler.sendDetailInputPairs = handler.sendDetailInputPairs[:len(handler.sendDetailInputPairs)-1]
	window.Changed()
}

// fetchCustomInputsCheck is called everytime the 'select custom inputs' checkbox is checked or unchecked.
// for everytime the checkbox is checked, the fetchCustomInputs function is called
func (handler *SendHandler) fetchCustomInputsCheck(wallet walletcore.Wallet, masterWindow nucular.MasterWindow) {
	if handler.selectCustomInputs {
		handler.fetchCustomInputs(wallet, masterWindow)
		return
	}

	handler.isFetchingUTXOS = false
	handler.utxos = nil
	masterWindow.Changed()
}

// fetchCustomInputs fetches custom inputs to add to send transaction. It adds unconfirmed inputs if the
// 'spend unconfirmed' checkbox is checked
func (handler *SendHandler) fetchCustomInputs(wallet walletcore.Wallet, masterWindow nucular.MasterWindow) {
	handler.isFetchingUTXOS = true
	masterWindow.Changed()

	defer func() {
		handler.isFetchingUTXOS = false
		masterWindow.Changed()
	}()

	var requiredConfirmations int32 = walletcore.DefaultRequiredConfirmations
	if handler.spendUnconfirmed {
		requiredConfirmations = 0
	}

	accountNumber := handler.accountNumbers[handler.selectedAccountIndex]
	utxos, err := wallet.UnspentOutputs(accountNumber, 0, requiredConfirmations)
	if err != nil {
		handler.fetchUTXOError = err
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

func (handler *SendHandler) validate(window *nucular.Window, wallet walletcore.Wallet) bool {
	isClean := true

	for _, input := range handler.sendDetailInputPairs {
		amountStr := string(input.amount.Buffer)
		if amountStr == "" {
			input.amountErr = "This amount field is required"
		} else {
			amountInt, err := strconv.ParseFloat(amountStr, 64)
			if err != nil {
				input.amountErr = "This is not a valid number"
			} else if amountInt < 1 {
				input.amountErr = "Send amount must be greater than 0DCR"
			} else {
				input.amountErr = ""
			}
		}
		if input.amountErr != "" {
			isClean = false
		}

		address := string(input.destinationAddress.Buffer)
		if address == "" {
			input.addressErr = "This address field is required"
		} else {
			isValid, err := wallet.ValidateAddress(address)
			if err != nil {
				input.addressErr = fmt.Sprintf("error validating address: %s", err.Error())
			} else if !isValid {
				input.addressErr = "Invalid address"
			} else {
				input.addressErr = ""
			}
		}

		if input.addressErr != "" {
			isClean = false
		}
	}

	window.Master().Changed()
	return isClean
}

func (handler *SendHandler) getPassphraseAndSubmit(window *nucular.Window, wallet walletcore.Wallet) {
	passphraseChan := make(chan string)
	widgets.NewPassphraseWidget().Get(window, passphraseChan)

	go func() {
		passphrase := <-passphraseChan
		if passphrase != "" {
			handler.submit(passphrase, window, wallet)
		}
	}()
}

func (handler *SendHandler) submit(passphrase string, window *nucular.Window, wallet walletcore.Wallet) {
	if handler.isSubmitting {
		return
	}

	handler.isSubmitting = true
	window.Master().Changed()

	defer func() {
		handler.isSubmitting = false
		window.Master().Changed()
	}()

	sendDestinations := make([]txhelper.TransactionDestination, len(handler.sendDetailInputPairs))
	for index := range handler.sendDetailInputPairs {
		amount, err := strconv.ParseFloat(string(handler.sendDetailInputPairs[index].amount.Buffer), 64)
		if err != nil {
			handler.err = err
			return
		}

		sendDestinations[index] = txhelper.TransactionDestination{
			Address: string(handler.sendDetailInputPairs[index].destinationAddress.Buffer),
			Amount:  amount,
		}
	}

	accountNumber := handler.accountNumbers[handler.selectedAccountIndex]
	var requiredConfirmations int32 = walletcore.DefaultRequiredConfirmations
	if handler.spendUnconfirmed {
		requiredConfirmations = 0
	}

	if handler.selectCustomInputs {
		utxos, totalInputAmount := handler.getUTXOSAndSelectedAmount()

		changeAddress, err := wallet.GenerateNewAddress(accountNumber)
		if err != nil {
			handler.err = err
			return
		}

		changeAmount, err := txhelper.EstimateChange(len(handler.utxos), int64(totalInputAmount), sendDestinations, []string{changeAddress})
		if err != nil {
			handler.err = err
			return
		}

		changeDestinations := []txhelper.TransactionDestination{{
			Amount:  dcrutil.Amount(changeAmount).ToCoin(),
			Address: changeAddress,
		}}

		handler.successHash, handler.err = wallet.SendFromUTXOs(accountNumber, requiredConfirmations, utxos, sendDestinations, changeDestinations, passphrase)
	} else {
		handler.successHash, handler.err = wallet.SendFromAccount(accountNumber, requiredConfirmations, sendDestinations, passphrase)
	}

	if handler.successHash != "" {
		handler.resetForm()
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

func (handler *SendHandler) resetForm() {
	if len(handler.accountNumbers) > 1 {
		handler.selectedAccountIndex = 0
	}

	for _, input := range handler.sendDetailInputPairs {
		input.amount.Buffer = []rune("")
		input.destinationAddress.Buffer = []rune("")
	}

	handler.spendUnconfirmed = false
	handler.selectCustomInputs = false
}
