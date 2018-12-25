package dcrwalletrpc

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/dcrutil"
	"github.com/decred/dcrwallet/netparams"
	"github.com/decred/dcrwallet/rpc/walletrpc"
	"github.com/raedahgroup/dcrcli/app/walletcore"
)

func amountToAtom(amountInDCR float64) (int64, error) {
	amountInAtom, err := dcrutil.NewAmount(amountInDCR)
	if err != nil {
		return 0, err
	}

	// type of amountInAtom is `dcrutil.Amount` which is an int64 alias
	return int64(amountInAtom), nil
}

func (c *WalletPRCClient) unspentOutputStream(account uint32, targetAmount int64) (walletrpc.WalletService_UnspentOutputsClient, error) {
	req := &walletrpc.UnspentOutputsRequest{
		Account:                  account,
		TargetAmount:             targetAmount,
		RequiredConfirmations:    0,
		IncludeImmatureCoinbases: true,
	}

	return c.walletService.UnspentOutputs(context.Background(), req)
}

func (c *WalletPRCClient) signAndPublishTransaction(serializedTx []byte, passphrase string) (string, error) {
	ctx := context.Background()

	// sign transaction
	signRequest := &walletrpc.SignTransactionRequest{
		Passphrase:            []byte(passphrase),
		SerializedTransaction: serializedTx,
	}

	signResponse, err := c.walletService.SignTransaction(ctx, signRequest)
	if err != nil {
		return "", fmt.Errorf("error signing transaction: %s", err.Error())
	}

	// publish transaction
	publishRequest := &walletrpc.PublishTransactionRequest{
		SignedTransaction: signResponse.Transaction,
	}

	publishResponse, err := c.walletService.PublishTransaction(ctx, publishRequest)
	if err != nil {
		return "", fmt.Errorf("error publishing transaction: %s", err.Error())
	}

	transactionHash, err := chainhash.NewHash(publishResponse.TransactionHash)
	if err != nil {
		return "", fmt.Errorf("error parsing successful transaction hash: %s", err.Error())
	}

	return transactionHash.String(), nil
}

func processTransactions(transactionDetails []*walletrpc.TransactionDetails) ([]*walletcore.Transaction, error) {
	transactions := make([]*walletcore.Transaction, 0, len(transactionDetails))

	for _, txDetail := range transactionDetails {
		// use any of the addresses in inputs/outputs to determine if this is a testnet tx
		var isTestnet bool
		for _, output := range txDetail.Credits {
			addr, err := dcrutil.DecodeAddress(output.Address)
			if err != nil {
				continue
			}

			isTestnet = !addr.IsForNet(netparams.MainNetParams.Params)
			break
		}

		hash, err := chainhash.NewHash(txDetail.Hash)
		if err != nil {
			return nil, err
		}

		amount, direction := transactionAmountAndDirection(txDetail)

		tx := &walletcore.Transaction{
			Hash:          hash.String(),
			Amount:        dcrutil.Amount(amount).ToCoin(),
			Fee:           dcrutil.Amount(txDetail.Fee).ToCoin(),
			Type:          txDetail.TransactionType.String(),
			Direction:     direction,
			Testnet:       isTestnet,
			Timestamp:     txDetail.Timestamp,
			FormattedTime: time.Unix(txDetail.Timestamp, 0).Format("Mon Jan 2, 2006 3:04PM"),
		}

		transactions = append(transactions, tx)
	}

	return transactions, nil
}

func transactionAmountAndDirection(txDetail *walletrpc.TransactionDetails) (int64, walletcore.TransactionDirection) {
	var outputAmounts int64
	for _, credit := range txDetail.Credits {
		outputAmounts += int64(credit.Amount)
	}

	var inputAmounts int64
	for _, debit := range txDetail.Debits {
		inputAmounts += int64(debit.PreviousAmount)
	}

	var amount int64
	var direction walletcore.TransactionDirection

	if txDetail.TransactionType == walletrpc.TransactionDetails_REGULAR {
		amountDifference := outputAmounts - inputAmounts
		if amountDifference < 0 && (float64(txDetail.Fee) == math.Abs(float64(amountDifference))) {
			// transferred internally, the only real amount spent was transaction fee
			direction = walletcore.TransactionDirectionTransferred
			amount = int64(txDetail.Fee)
		} else if amountDifference > 0 {
			// received
			direction = walletcore.TransactionDirectionReceived

			for _, credit := range txDetail.Credits {
				amount += int64(credit.Amount)
			}
		} else {
			// sent
			direction = walletcore.TransactionDirectionSent

			for _, debit := range txDetail.Debits {
				amount += int64(debit.PreviousAmount)
			}
			for _, credit := range txDetail.Credits {
				amount -= int64(credit.Amount)
			}
			amount -= int64(txDetail.Fee)
		}
	}

	return amount, direction
}
