package dcrwalletrpc

import (
	"context"
	"fmt"

	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrwallet/rpc/walletrpc"
	"github.com/raedahgroup/dcrlibwallet/txhelper"
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

func (c *WalletRPCClient) decodeTransactionWithTxSummary(ctx context.Context, txSummary *walletrpc.TransactionDetails,
	blockHash []byte) (*txhelper.Transaction, error) {

	var blockHeight int32 = -1
	if blockHash != nil {
		blockInfo, err := c.walletService.BlockInfo(ctx, &walletrpc.BlockInfoRequest{BlockHash: blockHash})
		if err == nil {
			blockHeight = blockInfo.BlockHeight
		}
	}

	accountName := func(accountNumber uint32) string {
		accountName, _ := c.AccountName(accountNumber)
		return accountName
	}

	walletInputs := make([]*txhelper.WalletInput, len(txSummary.Debits))
	for i, input := range txSummary.Debits {
		walletInputs[i] = &txhelper.WalletInput{
			Index:    int32(input.Index),
			AmountIn: int64(input.PreviousAmount),
			WalletAccount: &txhelper.WalletAccount{
				AccountNumber: int32(input.PreviousAccount),
				AccountName:   accountName(input.PreviousAccount),
			},
		}
	}

	walletOutputs := make([]*txhelper.WalletOutput, len(txSummary.Credits))
	for i, output := range txSummary.Credits {
		walletOutputs[i] = &txhelper.WalletOutput{
			Index:     int32(output.Index),
			AmountOut: int64(output.Amount),
			Address:   output.Address,
			WalletAccount: &txhelper.WalletAccount{
				AccountNumber: int32(output.Account),
				AccountName:   accountName(output.Account),
			},
		}
	}

	walletTx := &txhelper.TxInfoFromWallet{
		BlockHeight: blockHeight,
		Timestamp:   txSummary.Timestamp,
		Hex:         fmt.Sprintf("%x", txSummary.Transaction),
		Inputs:      walletInputs,
		Outputs:     walletOutputs,
	}

	return txhelper.DecodeTransaction(walletTx, c.activeNet.Params)
}