package libwallet

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/app/wallet"
)

func (lw *LibWallet) WalletExists() (bool, error) {
	return lw.dcrlw.WalletExists()
}

func (lw *LibWallet) CreateWallet(privatePass, seed string) error {
	return lw.dcrlw.CreateWallet(privatePass, seed)
}

// This method may stall if the wallet database is in use by some other process,
// hence the need for ctx, so user can cancel the operation if it's taking too long
// additionally, an error is returned if there's a delay in opening the wallet.
func (lw *LibWallet) OpenWallet(ctx context.Context, publicPass string) error {
	if lw.dcrlw.WalletOpened() {
		return nil
	}

	// wallet database is opened using bolt db by `github.com/decred/dcrwallet/wallet/internal/bdb`
	// bolt db stalls if the database is currently in use by another process,
	// waiting for the other process to release the file.
	// bold db doc advise setting a 1 second timeout to prevent this stalling.
	// see https://github.com/boltdb/bolt#opening-a-database
	walletOpenDelay := time.NewTicker(5 * time.Second)

	loadWalletDone := make(chan error)
	go func() {
		openWalletError := lw.dcrlw.OpenWallet([]byte(publicPass))
		loadWalletDone <- openWalletError
		walletOpenDelay.Stop()
	}()

	select {
	case <-walletOpenDelay.C:
		return fmt.Errorf("wallet database is in use by another process")

	case err := <-loadWalletDone:
		return err

	case <-ctx.Done():
		return ctx.Err()
	}
}

func (lw *LibWallet) AddSyncProgressListener(syncProgressListener dcrlibwallet.SyncProgressListener,
	uniqueIdentifier string) error {
	return lw.dcrlw.AddSyncProgressListener(syncProgressListener, uniqueIdentifier)
}

func (lw *LibWallet) RemoveSyncProgressListener(uniqueIdentifier string) {
	lw.dcrlw.RemoveSyncProgressListener(uniqueIdentifier)
}

func (lw *LibWallet) SpvSync(showLog bool, persistentPeers []string) error {
	if showLog {
		lw.dcrlw.EnableSyncLogs()
	}

	var peerAddresses string
	if persistentPeers != nil && len(persistentPeers) > 0 {
		peerAddresses = strings.Join(persistentPeers, ";")
	}

	return lw.dcrlw.SpvSync(peerAddresses)
}

func (lw *LibWallet) WalletConnectionInfo() (info wallet.ConnectionInfo, err error) {
	accounts, loadAccountErr := lw.AccountsOverview(wallet.DefaultRequiredConfirmations)
	if loadAccountErr != nil {
		err = fmt.Errorf("error fetching account balance: %s", loadAccountErr.Error())
		info.TotalBalance = "0 DCR"
	} else {
		var totalBalance dcrutil.Amount
		for _, acc := range accounts {
			totalBalance += acc.Balance.Total
		}
		info.TotalBalance = totalBalance.String()
	}

	bestBlock, bestBlockErr := lw.BestBlock()
	if bestBlockErr != nil && err != nil {
		err = fmt.Errorf("%s, error in fetching best block %s", err.Error(), bestBlockErr.Error())
	} else if bestBlockErr != nil {
		err = bestBlockErr
	}

	info.LatestBlock = bestBlock
	info.NetworkType = lw.NetType()
	info.PeersConnected = 0 // todo read from lw.dcrlw

	return
}

func (lw *LibWallet) BestBlock() (uint32, error) {
	return uint32(lw.dcrlw.GetBestBlock()), nil
}

func (lw *LibWallet) Shutdown() {
	lw.dcrlw.Shutdown()
}
