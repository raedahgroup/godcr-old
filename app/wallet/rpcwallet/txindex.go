package rpcwallet

import (
	"context"
	"fmt"
	"io"

	"github.com/decred/dcrwallet/rpc/walletrpc"
	"github.com/raedahgroup/dcrlibwallet/txhelper"
)

func (c *WalletRPCClient) indexTransactions(ctx context.Context, startBlockHeight int32, endBlockHeight int32,
	showLog bool, afterIndexing func()) (err error) {

	reportProgress := func(formatting string, values ...interface{}) {
		if showLog {
			fmt.Printf(formatting+".\n", values...)
		}
	}

	defer func() {
		afterIndexing()
		if err != nil {
			reportProgress("Tx indexing error: %v", err)
			return
		}

		// mark current end block height as last index point
		c.txIndexDB.SaveLastIndexPoint(endBlockHeight)

		count, err := c.TransactionCount(nil)
		if err != nil {
			reportProgress("Count tx error: %s", err.Error())
			return
		}
		reportProgress("Transaction indexing finished at %d, %d transaction(s) indexed in total",
			endBlockHeight, count)
	}()

	if startBlockHeight == -1 {
		startBlockHeight, err = c.txIndexDB.ReadIndexingStartBlock()
		if err != nil {
			reportProgress("Error reading block height to start tx indexing :%v", err)
			return err
		}
	}
	if startBlockHeight > endBlockHeight {
		bestBlockHeight, err := c.BestBlock()
		if err != nil {
			endBlockHeight = -1 // up to unmined blocks
		} else {
			endBlockHeight = int32(bestBlockHeight)
		}
	}

	reportProgress("Indexing transactions start height: %d, end height: %d", startBlockHeight, endBlockHeight)

	var totalIndexed int32
	indexTx := func(tx *txhelper.Transaction) error {
		err = c.txIndexDB.SaveOrUpdate(tx)
		if err != nil {
			reportProgress("Save or update tx error :%v", err)
			return err
		}

		totalIndexed++
		if c.syncListener == nil {
			c.syncListener.OnIndexTransactions(totalIndexed)
		}
		return nil
	}

	req := &walletrpc.GetTransactionsRequest{
		StartingBlockHeight: startBlockHeight,
		EndingBlockHeight:   endBlockHeight,
	}

	txStream, err := c.walletService.GetTransactions(ctx, req)
	if err != nil {
		return err
	}

	for {
		in, err := txStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// process mined txs and index
		if in.MinedTransactions != nil {
			for _, txSummary := range in.MinedTransactions.Transactions {
				tx, err := c.decodeTransactionWithTxSummary(ctx, txSummary, in.MinedTransactions.Hash)
				if err == nil {
					err = indexTx(tx)
				}
				if err != nil {
					return err
				}
			}

			err := c.txIndexDB.SaveLastIndexPoint(int32(in.MinedTransactions.Height))
			if err != nil {
				reportProgress("Error setting block height for last indexed tx: ", err)
				return err
			}

			reportProgress("Transaction index caught up to %d", in.MinedTransactions.Height)
		}

		// process unmined txs and index
		if in.UnminedTransactions != nil {
			for _, txSummary := range in.UnminedTransactions {
				tx, err := c.decodeTransactionWithTxSummary(ctx, txSummary, nil)
				if err == nil {
					err = indexTx(tx)
				}
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
