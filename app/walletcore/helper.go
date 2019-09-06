package walletcore

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/decred/dcrd/dcrutil"
	"github.com/decred/dcrwallet/wallet"
	"github.com/raedahgroup/dcrlibwallet/txhelper"
	"github.com/raedahgroup/dcrlibwallet/txindex"
	"github.com/raedahgroup/dcrlibwallet/utils"
)

const (
	// standard decred min confirmations is 2, this should be used as default for wallet operations
	DefaultRequiredConfirmations = 2

	// default number of transactions to return per call to Wallet.TransactionHistory()
	TransactionHistoryCountPerPage = 25

	ReceivingDecredHint = "Each time you request payment, a new address is generated to protect your privacy."

	TestnetHDPath       = "m / 44' / 1' / "
	LegacyTestnetHDPath = "m / 44' / 11' / "
	MainnetHDPath       = "m / 44' / 42' / "
	LegacyMainnetHDPath = "m / 44' / 20' / "

	TransactionFilterSent     = "Sent"
	TransactionFilterReceived = "Received"
	TransactionFilterYourself = "Yourself"
	TransactionFilterStaking  = "Staking"
	TransactionFilterCoinbase = "Coinbase"
)

var TransactionFilters = []string{
	TransactionFilterSent,
	TransactionFilterReceived,
	TransactionFilterYourself,
	TransactionFilterStaking,
	TransactionFilterCoinbase,
}

func BuildTransactionFilter(filters ...string) *txindex.ReadFilter {
	var (
		txFilter       = txindex.Filter()
		stakingTxTypes = []string{
			txhelper.FormatTransactionType(wallet.TransactionTypeTicketPurchase),
			txhelper.FormatTransactionType(wallet.TransactionTypeVote),
			txhelper.FormatTransactionType(wallet.TransactionTypeRevocation),
		}
	)
	for _, filter := range filters {
		switch filter {
		case TransactionFilterSent:
			txFilter = txFilter.AndForDirections(txhelper.TransactionDirectionSent)
			break
		case TransactionFilterReceived:
			txFilter = txFilter.AndForDirections(txhelper.TransactionDirectionReceived)
			break
		case TransactionFilterYourself:
			txFilter = txFilter.AndForDirections(txhelper.TransactionDirectionYourself).AndNotWithTxTypes(stakingTxTypes...)

			break
		case TransactionFilterCoinbase:
			txFilter = txFilter.AndWithTxTypes(txhelper.FormatTransactionType(wallet.TransactionTypeCoinbase))
			break
		case TransactionFilterStaking:
			txFilter = txFilter.OrWithTxTypes(stakingTxTypes...)
			break
		}
	}
	return txFilter
}

// NormalizeBalance adds 0s the right of balance to make it x.xxxxxxxx DCR
func NormalizeBalance(balance float64) string {
	return fmt.Sprintf("%010.8f DCR", balance)
}

func WalletBalance(accounts []*Account) string {
	var totalBalance, spendableBalance dcrutil.Amount
	for _, account := range accounts {
		totalBalance += account.Balance.Total
		spendableBalance += account.Balance.Spendable
	}

	if totalBalance != spendableBalance {
		return fmt.Sprintf("Total %s (Spendable %s)", totalBalance.String(), spendableBalance.String())
	} else {
		return totalBalance.String()
	}
}

// GetChangeDestinationsWithRandomAmounts generates change destination(s) based on the number of change addresses the user wants.
func GetChangeDestinationsWithRandomAmounts(wallet Wallet, nChangeOutputs int, amountInAtom int64, sourceAccount uint32,
	nUtxoSelection int, sendDestinations []txhelper.TransactionDestination) (changeOutputDestinations []txhelper.TransactionDestination, err error) {

	var changeAddresses []string
	for i := 0; i < nChangeOutputs; i++ {
		address, err := wallet.GenerateNewAddress(sourceAccount)
		if err != nil {
			return nil, fmt.Errorf("error generating address: %s", err.Error())
		}
		changeAddresses = append(changeAddresses, address)
	}

	changeAmount, err := txhelper.EstimateChange(nUtxoSelection, amountInAtom, sendDestinations, changeAddresses)
	if err != nil {
		return nil, fmt.Errorf("error in getting change amount: %s", err.Error())
	}
	if changeAmount <= 0 {
		return
	}

	var portionRations []float64
	var rationSum float64
	for i := 0; i < nChangeOutputs; i++ {
		portion := rand.Float64()
		portionRations = append(portionRations, portion)
		rationSum += portion
	}

	for i, portion := range portionRations {
		portionPercentage := portion / rationSum
		amount := portionPercentage * float64(changeAmount)

		changeOutput := txhelper.TransactionDestination{
			Address: changeAddresses[i],
			Amount:  dcrutil.Amount(amount).ToCoin(),
		}
		changeOutputDestinations = append(changeOutputDestinations, changeOutput)
	}
	return
}

func BuildTxDestinations(destinationAddresses, destinationAccounts,
	destinationAmounts, sendMaxAmountValues []string, generateAddressFn func(accountNumber uint32) (string, error)) (
	destinations []txhelper.TransactionDestination, totalAmount dcrutil.Amount, err error) {

	destinationAccountAddresses := make([]string, len(destinationAccounts))
	for i, accountNumberStr := range destinationAccounts {
		account, err := strconv.ParseInt(accountNumberStr, 10, 32)
		if err != nil {
			return nil, 0, fmt.Errorf("Invalid account number")
		}
		address, err := generateAddressFn(uint32(account))
		destinationAccountAddresses[i] = address
	}

	addressLength := len(destinationAddresses)
	if addressLength == 0 {
		addressLength = len(destinationAccountAddresses)
	}

	for i := 0; i < addressLength; i++ {
		var address string
		if len(destinationAddresses) > 0 {
			address = destinationAddresses[i]
		} else {
			address = destinationAccountAddresses[i]
		}
		destination := txhelper.TransactionDestination{
			Address: address,
			// only set SendMax to true if `sendMaxAmountValues` is not nil and this particular sendMaxAmountValue is "true"
			SendMax: sendMaxAmountValues != nil && sendMaxAmountValues[i] == "true",
		}

		if !destination.SendMax {
			var dcrSendAmount dcrutil.Amount
			destination.Amount, err = strconv.ParseFloat(destinationAmounts[i], 64)
			if err == nil {
				dcrSendAmount, err = dcrutil.NewAmount(destination.Amount)
			}

			if err != nil {
				err = fmt.Errorf("invalid destination amount: %s", destinationAmounts[i])
				return
			}
			totalAmount += dcrSendAmount
		}

		if destination.Amount == 0 && !destination.SendMax {
			err = fmt.Errorf("invalid request, cannot send 0 amount to %s", destination.Address)
			return
		}

		destinations = append(destinations, destination)
	}

	return
}

func SumUtxosInAccount(wallet Wallet, accountNumber uint32, requiredConfirmations int32) (utxos []string, total dcrutil.Amount, err error) {
	allUtxos, err := wallet.UnspentOutputs(accountNumber, 0, requiredConfirmations)
	if err != nil {
		return nil, 0, err
	}

	var totalInputAmountDcr float64
	for _, utxo := range allUtxos {
		utxos = append(utxos, utxo.OutputKey)
		totalInputAmountDcr += utxo.Amount.ToCoin()
	}

	total, err = dcrutil.NewAmount(totalInputAmountDcr)
	return
}

func TxDetails(tx *txhelper.Transaction, confirmations int32) *Transaction {
	return &Transaction{
		Transaction:   tx,
		ShortTime:     utils.ExtractDateOrTime(tx.Timestamp),
		LongTime:      utils.FormatUTCTime(tx.Timestamp),
		Confirmations: confirmations,
		Status:        txhelper.TxStatus(confirmations),
	}
}
