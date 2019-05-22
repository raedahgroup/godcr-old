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
	totalSendAmount       dcrutil.Amount
	changeDestinations    []txhelper.TransactionDestination
	totalInputAmount      dcrutil.Amount
}

// retrieveSendPagePayload parses the req for the send parameters submitted;
// the order of form fields on the front end is followed:
// source account - spend unconfirmed - custom inputs - send destinations - custom change outputs
func retrieveSendPagePayload(req *http.Request, addressFunc func(accountNumber uint32) (string, error)) (payload *sendPagePayload, err error) {
	payload = new(sendPagePayload)

	err = req.ParseForm()
	if err != nil {
		return nil, fmt.Errorf("error in parsing request: %s", err.Error())
	}

	selectedAccount := req.FormValue("source-account")
	account, err := strconv.ParseUint(selectedAccount, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("error in retreiving selected account: %s", err.Error())
	}
	payload.sourceAccount = uint32(account)

	spendUnconfirmed := req.FormValue("spend-unconfirmed")
	payload.requiredConfirmations = walletcore.DefaultRequiredConfirmations
	if spendUnconfirmed != "" {
		payload.requiredConfirmations = 0
	}

	// parse custom inputs form data
	payload.useCustom = req.FormValue("use-custom") != ""
	if payload.useCustom {
		// set selected utxos
		payload.utxos = req.Form["utxo"]

		// set total selected inputs amount
		var totalInputAmountDcr float64
		totalSelectedInputAmountDcr := req.FormValue("totalSelectedInputAmountDcr")
		if totalInputAmountDcr, err = strconv.ParseFloat(totalSelectedInputAmountDcr, 64); err == nil {
			payload.totalInputAmount, err = dcrutil.NewAmount(totalInputAmountDcr)
		}
		if err != nil {
			return nil, errors.New("cannot read total send amount value")
		}
	}

	// parse send destinations form data
	destinationAddresses := req.Form["destination-address"]
	destinationAccountNumbers := req.Form["destination-account-number"]
	destinationAmounts := req.Form["destination-amount"]
	sendMaxAmountChecks := req.Form["send-max-amount"]

	// sendMaxAmountChecks is a slice of send max amount true/false values
	// there are 2 `send-max-amount` input elements per destination - an hidden input field and a checkbox input field
	// the hidden input field always produces a "false" value, while the second input field produces "true" (if and only if checked)
	// no value is submitted for the second input element (the checkbox) if the checkbox is not checked
	// the implication is that for checkboxes that are checked, there'd be 2 values ("false" and "true")
	// while unchecked checkboxes will return only one value ("false")
	// if any value in the sendMaxAmountChecks slice is true, then the previous "false" value should be ignored
	// as both values refer to the same destination
	// at the end of the day, the number of send max amount check values should be equal to the number of destination address values
	actualSendMaxAmountValues := make([]string, 0, len(destinationAmounts))
	for _, sendMaxAmountCheckValue := range sendMaxAmountChecks {
		if sendMaxAmountCheckValue == "true" {
			// replace previous value in `actualSendMaxAmountValues` slice to true and ignore the previously set value
			previousValueIndex := len(actualSendMaxAmountValues) - 1
			actualSendMaxAmountValues[previousValueIndex] = sendMaxAmountCheckValue
		} else {
			actualSendMaxAmountValues = append(actualSendMaxAmountValues, sendMaxAmountCheckValue)
		}
	}

	payload.sendDestinations, payload.totalSendAmount, err = walletcore.BuildTxDestinations(destinationAddresses,
		destinationAccountNumbers, destinationAmounts, actualSendMaxAmountValues, addressFunc)
	if err != nil {
		return nil, fmt.Errorf("error in parsing send destinations: %s", err.Error())
	}

	// parse custom change outputs form data
	changeOutputAddresses := req.Form["change-output-address"]
	changeOutputAmounts := req.Form["change-output-amount"]
	payload.changeDestinations, _, err = walletcore.BuildTxDestinations(changeOutputAddresses, nil, changeOutputAmounts, nil, addressFunc)
	if err != nil {
		return nil, fmt.Errorf("error in parsing change destinations: %s", err.Error())
	}

	payload.passphrase = req.FormValue("wallet-passphrase")

	return
}
