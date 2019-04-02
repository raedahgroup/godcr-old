package routes

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet/txhelper"
	"github.com/raedahgroup/godcr/app/walletcore"
)

type sendPagePayload struct {
	Utxos []string
	SourceAccount uint32
	Passphrase string
	RequiredConfirmations int32
	UseCustom bool
	SendDestinations []txhelper.TransactionDestination
	ChangeDestinations []txhelper.TransactionDestination
	TotalInputAmount int64

}


func retrieveSendPagePayload(req *http.Request) (payload *sendPagePayload, err error) {
	payload = new(sendPagePayload)

	req.ParseForm()
	payload.Utxos = req.Form["utxo"]
	selectedAccount := req.FormValue("source-account")
	payload.Passphrase = req.FormValue("wallet-passphrase")
	spendUnconfirmed := req.FormValue("spend-unconfirmed")
	payload.UseCustom = strings.EqualFold(req.FormValue("use-custom"), "true")

	destinationAddresses := req.Form["destination-address"]
	destinationAmounts := req.Form["destination-amount"]

	sendDestinations, err := walletcore.BuildTxDestinations(destinationAddresses, destinationAmounts)
	if err != nil {
		return nil, err
	}
	payload.SendDestinations = sendDestinations

	account, err := strconv.ParseUint(selectedAccount, 10, 32)
	if err != nil {
		return nil, err
	}
	payload.SourceAccount = uint32(account)

	payload.RequiredConfirmations = walletcore.DefaultRequiredConfirmations
	if spendUnconfirmed != "" {
		payload.RequiredConfirmations = 0
	}

	totalSelectedInputAmountDcr := req.FormValue("totalSelectedInputAmountDcr")
	totalInputAmountDcr, err := strconv.ParseFloat(totalSelectedInputAmountDcr, 64)
	if err != nil {
		return
	}

	totalInputAmount, err := dcrutil.NewAmount(totalInputAmountDcr)
	if err != nil {
		return
	}

	payload.TotalInputAmount = int64(totalInputAmount)

	changeOutputAddresses := req.Form["change-output-address"]
	changeOutputAmounts := req.Form["change-output-amount"]

	changeDestinations, err := walletcore.BuildTxDestinations(changeOutputAddresses, changeOutputAmounts)
	if err != nil {
		return
	}
	payload.ChangeDestinations = changeDestinations

	return
}
