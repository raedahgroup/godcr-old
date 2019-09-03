package rpcwallet

import (
	"context"
	"fmt"

	"github.com/decred/dcrwallet/rpc/walletrpc"
	"github.com/raedahgroup/dcrlibwallet/txhelper"
)

type TransactionListener interface {
	OnTransaction(transaction *txhelper.Transaction)
	OnTransactionConfirmed(hash string, height int32)
	OnBlockAttached(height int32, timestamp int64)
}

func (c *WalletRPCClient) RegisterTxNotificationListener(listener TransactionListener) {
	c.txNotificationListener = listener
}

func (c *WalletRPCClient) ListenForTxNotification(ctx context.Context) error {
	txNotificationStream, err := c.walletService.TransactionNotifications(ctx, &walletrpc.TransactionNotificationsRequest{})
	if err != nil {
		return fmt.Errorf("cannot start tx notification listener: %s", err.Error())
	}

	go c.startTxNotificationListener(ctx, txNotificationStream)
	return nil
}

func (c *WalletRPCClient) startTxNotificationListener(ctx context.Context,
	txNotificationStream walletrpc.WalletService_TransactionNotificationsClient) {

	for {
		txNotification, err := txNotificationStream.Recv()

		if ctx.Err() != nil {
			// context canceled, stop listening for notifications
			return
		}

		if err != nil {
			// todo use logger, similar logging should be done across rpcwallet and dcrlibwallet functions
			fmt.Printf("error reading tx notification update: %s\n", err.Error())
		}

		// process unmined tx gotten from notification
		for _, txSummary := range txNotification.UnminedTransactions {
			decodedTx, err := c.decodeTransactionWithTxSummary(ctx, txSummary, nil)
			if err != nil {
				continue
			}

			err = c.txIndexDB.SaveOrUpdate(decodedTx)
			if err == nil && c.txNotificationListener != nil {
				c.txNotificationListener.OnTransaction(decodedTx)
				continue
			}
		}

		// process mined tx gotten from notification
		for _, block := range txNotification.AttachedBlocks {
			if c.txNotificationListener != nil {
				c.txNotificationListener.OnBlockAttached(block.Height, block.Timestamp)
			}

			blockHash := block.Hash
			for _, txSummary := range block.Transactions {
				decodedTx, err := c.decodeTransactionWithTxSummary(ctx, txSummary, blockHash)
				if err != nil {
					continue
				}

				err = c.txIndexDB.SaveOrUpdate(decodedTx)
				if err != nil {
					continue
				}

				if c.txNotificationListener != nil {
					c.txNotificationListener.OnTransactionConfirmed(decodedTx.Hash, block.Height)
				}
			}
		}
	}
}
