package dcrlibwallet

import (
	"time"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/dcrlibwallet/txhelper"
	"github.com/raedahgroup/godcr/app/walletcore"
)

func processAndAppendTransactions(rawTxs []*dcrlibwallet.Transaction, processedTxs []*walletcore.Transaction) (
	[]*walletcore.Transaction, error) {

	for _, tx := range rawTxs {
		_, txFee, txSize, txFeeRate, err := txhelper.MsgTxFeeSizeRate(tx.Transaction)
		if err != nil {
			return nil, err
		}

		processedTxs = append(processedTxs, &walletcore.Transaction{
			Hash:          tx.Hash,
			Amount:        dcrutil.Amount(tx.Amount),
			Fee:           txFee,
			FeeRate:       txFeeRate,
			Size:          txSize,
			Type:          walletcore.FormatTxType(tx.Type),
			Direction:     tx.Direction,
			Timestamp:     tx.Timestamp,
			FormattedTime: time.Unix(tx.Timestamp, 0).Format("2006-01-02 15:04:05"),
		})
	}

	return processedTxs, nil
}
