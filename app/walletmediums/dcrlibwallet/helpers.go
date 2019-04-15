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

	bestBlockHeight := lib.walletLib.GetBestBlock()

	for _, tx := range rawTxs {
		_, txFee, txSize, txFeeRate, err := txhelper.MsgTxFeeSizeRate(tx.Hex)
		if err != nil {
			return nil, err
		}

		_, status := walletcore.TxStatus(tx.BlockHeight, bestBlockHeight)

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
