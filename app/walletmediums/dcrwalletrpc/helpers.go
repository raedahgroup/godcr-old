package dcrwalletrpc

import (
	"context"
	"fmt"
	"time"

	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/dcrutil"
	"github.com/decred/dcrwallet/rpc/walletrpc"
	"github.com/raedahgroup/dcrlibwallet/txhelper"
	"github.com/raedahgroup/godcr/app/walletcore"
)

func (c *WalletRPCClient) unspentOutputStream(account uint32, targetAmount int64, requiredConfirmations int32) (walletrpc.WalletService_UnspentOutputsClient, error) {
	req := &walletrpc.UnspentOutputsRequest{
		Account:                  account,
		TargetAmount:             targetAmount,
		RequiredConfirmations:    requiredConfirmations,
		IncludeImmatureCoinbases: true,
	}

	return c.walletService.UnspentOutputs(context.Background(), req)
}

func (c *WalletRPCClient) signAndPublishTransaction(serializedTx []byte, passphrase string) (string, error) {
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
		tx, err := processTransaction(txDetail)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}

	return transactions, nil
}

func processTransaction(txDetail *walletrpc.TransactionDetails) (*walletcore.Transaction, error) {
	hash, err := chainhash.NewHash(txDetail.Hash)
	if err != nil {
		return nil, err
	}

	_, txFee, txSize, txFeeRate, err := txhelper.MsgTxFeeSizeRate(txDetail.Transaction)
	if err != nil {
		return nil, err
	}

	amount, direction := transactionAmountAndDirection(txDetail)

	tx := &walletcore.Transaction{
		Hash:          hash.String(),
		Amount:        dcrutil.Amount(amount),
		Fee:           txFee,
		FeeRate:       txFeeRate,
		Type:          txhelper.RPCTransactionType(txDetail.TransactionType),
		Direction:     direction,
		Timestamp:     txDetail.Timestamp,
		FormattedTime: time.Unix(txDetail.Timestamp, 0).Format("2006-01-02 15:04:05"),
		Size:          txSize,
	}
	return tx, nil
}

func transactionAmountAndDirection(txDetail *walletrpc.TransactionDetails) (int64, txhelper.TransactionDirection) {
	var outputTotal int64
	for _, credit := range txDetail.Credits {
		outputTotal += int64(credit.Amount)
	}

	var inputTotal int64
	for _, debit := range txDetail.Debits {
		inputTotal += int64(debit.PreviousAmount)
	}

	return txhelper.TransactionAmountAndDirection(inputTotal, outputTotal, int64(txDetail.Fee))
}
