package routes

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet/txhelper"
	"github.com/raedahgroup/godcr/app/walletcore"
)

type sendPagePayload struct {
	utxos                 []string
	sourceAccount         uint32
	passphrase            string
	requiredConfirmations int32
	useCustom             bool
	sendDestinations      []txhelper.TransactionDestination
	totalSendAmount		  int64
	changeDestinations    []txhelper.TransactionDestination
	totalInputAmount      int64
}

func retrieveSendPagePayload(req *http.Request) (payload *sendPagePayload, err error) {
	payload = new(sendPagePayload)

	req.ParseForm()
	payload.utxos = req.Form["utxo"]
	selectedAccount := req.FormValue("source-account")
	payload.passphrase = req.FormValue("wallet-passphrase")
	spendUnconfirmed := req.FormValue("spend-unconfirmed")
	payload.useCustom = req.FormValue("use-custom") != ""

	destinationAddresses := req.Form["destination-address"]
	destinationAmounts := req.Form["destination-amount"]

	sendDestinations, err := walletcore.BuildTxDestinations(destinationAddresses, destinationAmounts)
	if err != nil {
		return nil, errors.New("invalid source account selected")
	}
	payload.sendDestinations = sendDestinations

	account, err := strconv.ParseUint(selectedAccount, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("error in retreiving selected account: %s", err.Error())
	}
	payload.sourceAccount = uint32(account)

	payload.requiredConfirmations = walletcore.DefaultRequiredConfirmations
	if spendUnconfirmed != "" {
		payload.requiredConfirmations = 0
	}

	totalSelectedInputAmountDcr := req.FormValue("totalSelectedInputAmountDcr")
	totalInputAmountDcr, err := strconv.ParseFloat(totalSelectedInputAmountDcr, 64)
	if err != nil {
		return nil, errors.New("cannot read total send amount value")
	}

	totalInputAmount, err := dcrutil.NewAmount(totalInputAmountDcr)
	if err != nil {
		err = errors.New("cannot read total send amount value")
		return
	}

	payload.totalInputAmount = int64(totalInputAmount)

	changeOutputAddresses := req.Form["change-output-address"]
	changeOutputAmounts := req.Form["change-output-amount"]

	changeDestinations, err := walletcore.BuildTxDestinations(changeOutputAddresses, changeOutputAmounts)
	if err != nil {
		return nil, fmt.Errorf("error in building change destinations: %s", err.Error())
	}
	payload.changeDestinations = changeDestinations

	return
}
