package handlers

import (
	"strconv"
	"time"

	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/handlers/widgets"
	"github.com/raedahgroup/godcr/nuklear/helpers"
)

type utxoSelection struct {
	selected bool
	utxo     *walletcore.UnspentOutput
}

type SendHandler struct {
	err            error
	fetchUTXOError error
	isRendering    bool

	utxos           []*utxoSelection
	isFetchingUTXOS bool

	accountNumbers   []uint32
	accountOverviews []string

	destinationAddressInputs []nucular.TextEditor
	sendAmountInputs         []nucular.TextEditor

	selectedAccountIndex int
	spendUnconfirmed     bool
	selectCustomInputs   bool
}

func (handler *SendHandler) BeforeRender() {
	handler.isRendering = false
	handler.spendUnconfirmed = false
	handler.isFetchingUTXOS = false
}

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
		handler.accountOverviews[index] = account.Name + string(account.Balance.Spendable)
		handler.accountNumbers[index] = account.Number
	}
	handler.selectedAccountIndex = 0
}

func (handler *SendHandler) addSendInputPair(updateWindow bool, window nucular.MasterWindow) {
	if handler.sendAmountInputs == nil {
		handler.sendAmountInputs = []nucular.TextEditor{}
		handler.destinationAddressInputs = []nucular.TextEditor{}
	}

	handler.sendAmountInputs = append(handler.sendAmountInputs, nucular.TextEditor{})
	handler.destinationAddressInputs = append(handler.destinationAddressInputs, nucular.TextEditor{})

	if updateWindow {
		window.Changed()
	}
}

func (handler *SendHandler) removeLastSendInputPair(window nucular.MasterWindow) {
	handler.sendAmountInputs = handler.sendAmountInputs[:len(handler.sendAmountInputs)-1]
	handler.destinationAddressInputs = handler.sendAmountInputs[:len(handler.destinationAddressInputs)-1]

	window.Changed()
}

func (handler *SendHandler) fetchCustomInputsCheck(wallet walletcore.Wallet, window *nucular.Window) {
	if handler.selectCustomInputs {
		handler.fetchCustomInputs(wallet, window)
		return
	}

	handler.isFetchingUTXOS = false
	handler.utxos = nil
	window.Master().Changed()
}

func (handler *SendHandler) fetchCustomInputs(wallet walletcore.Wallet, window *nucular.Window) {
	handler.isFetchingUTXOS = true
	window.Master().Changed()

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
	window.Master().Changed()
}

func (handler *SendHandler) Render(window *nucular.Window, wallet walletcore.Wallet) {
	if !handler.isRendering {
		handler.isRendering = true
		handler.fetchAccounts(wallet)
		handler.addSendInputPair(false, window.Master())
	}

	// draw page
	if pageWindow := helpers.NewWindow("Send Page", window, 0); pageWindow != nil {
		pageWindow.DrawHeader("Send")

		// content window
		if contentWindow := pageWindow.ContentWindow("Send Form"); contentWindow != nil {
			if handler.err != nil {
				contentWindow.SetErrorMessage(handler.err.Error())
			} else {
				helpers.SetFont(window, helpers.PageContentFont)

				contentWindow.Row(10).Dynamic(2)
				contentWindow.Label("Source Account", "LC")

				contentWindow.Row(25).Dynamic(2)
				handler.selectedAccountIndex = contentWindow.ComboSimple(handler.accountOverviews, handler.selectedAccountIndex, 25)

				contentWindow.Row(15).Dynamic(2)
				if contentWindow.CheckboxText("Spend Unconfirmed", &handler.spendUnconfirmed) {
					handler.fetchCustomInputsCheck(wallet, contentWindow.Window)
				}

				for i := 0; i < len(handler.destinationAddressInputs); i++ {
					contentWindow.Row(10).Dynamic(2)
					contentWindow.Label("Destination Address", "LC")
					contentWindow.Label("Amount (DCR)", "LC")

					contentWindow.Row(25).Dynamic(2)
					handler.destinationAddressInputs[i].Edit(contentWindow.Window)
					handler.sendAmountInputs[i].Edit(contentWindow.Window)
				}

				contentWindow.Row(25).Dynamic(2)
				if contentWindow.ButtonText("Add anohter address") {
					handler.addSendInputPair(true, window.Master())
				}

				if len(handler.sendAmountInputs) > 1 {
					if contentWindow.ButtonText("Remove last address") {
						handler.removeLastSendInputPair(window.Master())
					}
				}

				contentWindow.Row(15).Dynamic(2)
				if contentWindow.CheckboxText("Select custom inputs", &handler.selectCustomInputs) {
					handler.fetchCustomInputsCheck(wallet, contentWindow.Window)
				}

				if handler.isFetchingUTXOS {
					widgets.ShowIsFetching(contentWindow)
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
			}
			contentWindow.End()
		}
		pageWindow.End()
	}
}
