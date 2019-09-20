package utils

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/decred/dcrd/dcrutil"
	"github.com/decred/dcrd/wire"
	"github.com/decred/dcrwallet/wallet/txrules"
	"github.com/raedahgroup/dcrlibwallet/txhelper"
	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/app/walletcore"
)

func DecimalPortion(n float64) string {
	decimalPlaces := fmt.Sprintf("%f", n-math.Floor(n))          // produces 0.xxxx0000
	decimalPlaces = strings.Replace(decimalPlaces, "0.", "", -1) // remove 0.
	decimalPlaces = strings.TrimRight(decimalPlaces, "0")        // remove trailing 0s
	return decimalPlaces
}

func SplitAmountIntoParts(amount float64) []string {
	balanceParts := make([]string, 3)

	wholeNumber := int(math.Floor(amount))
	balanceParts[0] = strconv.Itoa(wholeNumber)

	decimalPortion := DecimalPortion(amount)
	if len(decimalPortion) == 0 {
		balanceParts[2] = " DCR"
	} else if len(decimalPortion) <= 2 {
		balanceParts[1] = fmt.Sprintf(".%s DCR", decimalPortion)
	} else {
		balanceParts[1] = fmt.Sprintf(".%s", decimalPortion[0:2])
		balanceParts[2] = fmt.Sprintf("%s DCR", decimalPortion[2:])
	}

	return balanceParts
}

func MaxDecimalPlaces(amounts []int64) (maxDecimalPlaces int) {
	for _, amount := range amounts {
		decimalPortion := DecimalPortion(dcrutil.Amount(amount).ToCoin())
		nDecimalPlaces := len(decimalPortion)
		if nDecimalPlaces > maxDecimalPlaces {
			maxDecimalPlaces = nDecimalPlaces
		}
	}
	return
}

func FormatAmountDisplay(amount int64, maxDecimalPlaces int) string {
	dcrAmount := dcrutil.Amount(amount).ToCoin()
	wholeNumber := int(math.Floor(dcrAmount))
	decimalPortion := DecimalPortion(dcrAmount)

	if len(decimalPortion) == 0 {
		return fmt.Sprintf("%2d%-*s DCR", wholeNumber, maxDecimalPlaces+1, decimalPortion)
	} else {
		return fmt.Sprintf("%2d.%-*s DCR", wholeNumber, maxDecimalPlaces, decimalPortion)
	}
}

type Account struct {
	IsSetAsHidden         bool
	IsSetAsDefaultAccount bool
	Account               *walletcore.Account
}

func FetchAccounts(requiredConfirmations int32, settings *config.Settings, wallet walletcore.Wallet) ([]Account, error) {
	walletAccounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		return nil, err
	}

	accounts := make([]Account, len(walletAccounts))
	for index, accountItem := range walletAccounts {
		var isSetAsHidden, isSetAsDefaultAccount bool
		for _, hiddenAccount := range settings.HiddenAccounts {
			if uint32(hiddenAccount) == accountItem.Number {
				isSetAsHidden = true
				break
			}
		}

		if settings.DefaultAccount == accountItem.Number {
			isSetAsDefaultAccount = true
		}

		accounts[index] = Account{
			IsSetAsHidden:         isSetAsHidden,
			IsSetAsDefaultAccount: isSetAsDefaultAccount,
			Account:               accountItem,
		}
	}

	return accounts, nil
}

func EstimateFee(numberOfInputs int, destinations []txhelper.TransactionDestination) (dcrutil.Amount, error) {
	maxSignedSize, err := EstimateSerializeSize(numberOfInputs, destinations)
	if err != nil {
		return 0, err
	}
	relayFeePerKb := txrules.DefaultRelayFeePerKb
	maxRequiredFee := txrules.FeeForSerializeSize(relayFeePerKb, maxSignedSize)

	return maxRequiredFee, err
}

func EstimateSerializeSize(numberOfInputs int, destinations []txhelper.TransactionDestination) (int, error) {
	outputs, err := makeTxOutputs(destinations)
	if err != nil {
		return 0, err
	}

	var changeAddresses []string
	for _, destination := range destinations {
		changeAddresses = append(changeAddresses, destination.Address)
	}

	totalChangeScriptSize, err := calculateChangeScriptSize(changeAddresses)
	if err != nil {
		return 0, err
	}

	scriptSizes := make([]int, numberOfInputs)
	for i := 0; i < numberOfInputs; i++ {
		scriptSizes[i] = txhelper.RedeemP2PKHSigScriptSize
	}
	maxSignedSize := txhelper.EstimateSerializeSize(scriptSizes, outputs, totalChangeScriptSize)

	return maxSignedSize, nil
}

func makeTxOutputs(destinations []txhelper.TransactionDestination) (outputs []*wire.TxOut, err error) {
	for _, destination := range destinations {
		var output *wire.TxOut
		output, err = txhelper.MakeTxOutput(destination)
		if err != nil {
			return
		}

		outputs = append(outputs, output)
	}
	return
}

func calculateChangeScriptSize(changeAddresses []string) (int, error) {
	var totalChangeScriptSize int
	for _, changeAddress := range changeAddresses {
		changeSource, err := txhelper.MakeTxChangeSource(changeAddress)
		if err != nil {
			return 0, err
		}
		totalChangeScriptSize += changeSource.ScriptSize()
	}
	return totalChangeScriptSize, nil
}
