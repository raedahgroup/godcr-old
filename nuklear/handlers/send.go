package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/decred/dcrd/dcrutil"

	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/dcrlibwallet/txhelper"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/handlers/widgets"
	"github.com/raedahgroup/godcr/nuklear/helpers"
)

type utxoSelection struct {
	selected bool
	utxo     *walletcore.UnspentOutput
}

type input struct {
	addressErrStr      string
	amountErrStr       string
	destinationAddress nucular.TextEditor
	amount             nucular.TextEditor
}

type SendHandler struct {
	err            error
	fetchUTXOError error
	isRendering    bool

	utxos           []*utxoSelection
	isFetchingUTXOS bool

	accountNumbers   []uint32
	accountOverviews []string

	inputs []*input

	selectedAccountIndex int
	spendUnconfirmed     bool
	selectCustomInputs   bool

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
	handler.inputs = nil
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
		handler.addSendInputPair(false, window.Master())
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

			for i := 0; i < len(handler.inputs); i++ {
				contentWindow.Row(10).Dynamic(2)
				contentWindow.Label("Destination Address", "LC")
				contentWindow.Label("Amount (DCR)", "LC")

				contentWindow.Row(25).Dynamic(2)
				handler.inputs[i].destinationAddress.Edit(contentWindow.Window)
				handler.inputs[i].amount.Edit(contentWindow.Window)

				if handler.inputs[i].addressErrStr != "" || handler.inputs[i].amountErrStr != "" {
					contentWindow.Row(10).Dynamic(2)
					contentWindow.LabelColored(handler.inputs[i].addressErrStr, "LC", helpers.DangerColor)
					contentWindow.LabelColored(handler.inputs[i].amountErrStr, "LC", helpers.DangerColor)
				}
			}

			contentWindow.Row(25).Dynamic(2)
			if contentWindow.ButtonText("Add anohter address") {
				handler.addSendInputPair(true, window.Master())
			}

			if len(handler.inputs) > 1 {
				if contentWindow.ButtonText("Remove last address") {
					handler.removeLastSendInputPair(window.Master())
				}
			}

			contentWindow.Row(15).Dynamic(2)
			if contentWindow.CheckboxText("Select custom inputs", &handler.selectCustomInputs) {
				handler.fetchCustomInputsCheck(wallet, masterWindow)
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

			submitButtonText := "Submit"
			if handler.isSubmitting {
				submitButtonText = "Submitting"
			}

			contentWindow.Row(25).Dynamic(2)
			if contentWindow.ButtonText(submitButtonText) {
				handler.validateAndSubmit(window, wallet)
			}
			contentWindow.End()
		}
		pageWindow.End()
	}
}

// fetch accounts for select source account field
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
func (handler *SendHandler) addSendInputPair(updateWindow bool, window nucular.MasterWindow) {
	if handler.inputs == nil {
		handler.inputs = []*input{}
	}

	item := &input{
		amount:             nucular.TextEditor{},
		destinationAddress: nucular.TextEditor{},
	}
	handler.inputs = append(handler.inputs, item)

	if updateWindow {
		window.Changed()
	}
}

// removeLastSendInputPair removes the last destinationAddress and amount input pair
// when the 'remove last address' button is clicked. The button is hidden when only one
// input pair exists on form
func (handler *SendHandler) removeLastSendInputPair(window nucular.MasterWindow) {
	handler.inputs = handler.inputs[:len(handler.inputs)-1]
	window.Changed()
}

// fetchCustomInputsCheck is called everytime the 'select custom inputs' checbox is checked or unchecked.
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
	handler.isFetchingUTXOS = false
	masterWindow.Changed()
}

func (handler *SendHandler) validateAndSubmit(window *nucular.Window, wallet walletcore.Wallet) {
	isClean := true

	for _, input := range handler.inputs {
		if string(input.amount.Buffer) == "" {
			input.amountErrStr = "This amount field is required"
			isClean = false
		} else {
			input.amountErrStr = ""
		}

		if string(input.destinationAddress.Buffer) == "" {
			input.addressErrStr = "This address field is required"
			isClean = false
		} else {
			input.addressErrStr = ""
		}
	}

	if isClean {
		passphraseChan := make(chan string)
		widgets.NewPassphraseWidget().Get(window, passphraseChan)

		go func() {
			passphrase := <-passphraseChan
			if passphrase != "" {
				handler.submit(passphrase, window, wallet)
			}
		}()
		return
	}
	window.Master().Changed()
}

func (handler *SendHandler) submit(passphrase string, window *nucular.Window, wallet walletcore.Wallet) {
	handler.isSubmitting = true
	window.Master().Changed()

	defer window.Master().Changed()

	sendDestinations := make([]txhelper.TransactionDestination, len(handler.inputs))
	for index := range handler.inputs {
		amount, err := strconv.ParseFloat(string(handler.inputs[index].amount.Buffer), 64)
		if err != nil {
			handler.err = err
			return
		}

		sendDestinations[index] = txhelper.TransactionDestination{
			Address: string(handler.inputs[index].destinationAddress.Buffer),
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
}

func (handler *SendHandler) getUTXOSAndSelectedAmount() ([]string, dcrutil.Amount) {
	var totalInputAmount dcrutil.Amount
	var utxos []string

	for i := range handler.utxos {
		if handler.utxos[i].selected {
			totalInputAmount += handler.utxos[i].utxo.Amount
		}
		utxos = append(utxos, handler.utxos[i].utxo.Address)
	}

	return utxos, totalInputAmount
}
