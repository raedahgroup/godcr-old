package dcrlibwallet

import (
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/dcrlibwallet/txhelper"
	"github.com/raedahgroup/dcrlibwallet/utils"
	"github.com/raedahgroup/godcr/app/walletcore"
)

func (lib *DcrWalletLib) processAndAppendTransactions(rawTxs []*dcrlibwallet.Transaction, processedTxs []*walletcore.Transaction) (
	[]*walletcore.Transaction, error) {

	for _, tx := range rawTxs {
		_, txFee, txSize, txFeeRate, err := txhelper.MsgTxFeeSizeRate(tx.Hex)
		if err != nil {
			return nil, err
		}

		// todo this is not very performant, fetching tx details for each tx in history simply to get tx status...
		var status string
		txDetails, err := lib.GetTransaction(tx.Hash)
		if err != nil {
			status = "Unknown"
		} else if txDetails.Confirmations >= walletcore.DefaultRequiredConfirmations {
			status = "Confirmed"
		} else {
			status = "Unconfirmed"
		}

		processedTxs = append(processedTxs, &walletcore.Transaction{
			Hash:          tx.Hash,
			Amount:        dcrutil.Amount(tx.Amount).String(),
			RawAmount:     tx.Amount,
			Fee:           txFee.String(),
			FeeRate:       txFeeRate,
			Size:          txSize,
			Type:          tx.Type,
			Direction:     tx.Direction,
			Status:        status,
			Timestamp:     tx.Timestamp,
			FormattedTime: utils.ExtractDateOrTime(tx.Timestamp),
		})
	}

	return processedTxs, nil
}
