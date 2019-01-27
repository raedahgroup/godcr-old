package routes

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
)

type Blockchain struct {
	sync.RWMutex
	_status walletcore.SyncStatus
	_report string
}

func (routes *Routes) walletLoaderMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return routes.walletLoaderFn(next)
	}
}

// walletLoaderFn checks if wallet is not open, attempts to open it and also perform sync the blockchain
// an error page is displayed and the actual route handler is not called, if ...
// - wallet doesn't exist (hasn't been created)
// - wallet exists but is not open
// - wallet is open but blockchain isn't synced
func (routes *Routes) walletLoaderFn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// render error on page if errMsg != ""
		var errMsg string
		defer func() {
			if errMsg != "" {
				routes.renderError(errMsg, res)
			}
		}()

		// check if wallet exists
		walletExists, err := routes.walletMiddleware.WalletExists()
		if err != nil {
			errMsg = fmt.Sprintf("Error checking for wallet: %s", err.Error())
			return
		}
		if !walletExists {
			routes.renderNoWalletError(res)
			return
		}

		if !routes.walletMiddleware.IsWalletOpen() {
			errMsg = "Wallet is not open. Restart the server"
			return
		}

		// wallet is open, check if blockchain is synced
		blockchainSyncStatus := routes.blockchain.status()
		blockchainSyncStatus = syncStatusSuccess
		switch blockchainSyncStatus {
		case walletcore.SyncStatusSuccess:
			next.ServeHTTP(res, req)
		case walletcore.SyncStatusNotStarted:
			errMsg = "Cannot display page. Blockchain hasn't been synced"
		case walletcore.SyncStatusInProgress:
			errMsg = fmt.Sprintf("%s. Refresh after a while to access this page", routes.blockchain.report())
		case walletcore.SyncStatusError:
			errMsg = fmt.Sprintf("Cannot display page. %s", routes.blockchain.report())
		default:
			errMsg = "Cannot display page. Blockchain sync status cannot be determined"
		}
	})
}

func (routes *Routes) syncBlockchain() {
	updateStatus := routes.blockchain.updateStatus

	err := routes.walletMiddleware.SyncBlockChain(&app.BlockChainSyncListener{
		SyncStarted: func() {
			updateStatus("Blockchain sync started...", walletcore.SyncStatusInProgress)
		},
		SyncEnded: func(err error) {
			if err != nil {
				updateStatus(fmt.Sprintf("Blockchain sync completed with error: %s", err.Error()), walletcore.SyncStatusError)
			} else {
				updateStatus("Blockchain sync completed successfully", walletcore.SyncStatusSuccess)
			}
		},
		OnHeadersFetched: func(percentageProgress int64) {
			updateStatus(fmt.Sprintf("Blockchain sync in progress. Fetching headers (1/3): %d%%", percentageProgress), walletcore.SyncStatusInProgress)
		},
		OnDiscoveredAddress: func(_ string) {
			updateStatus("Blockchain sync in progress. Discovering addresses (2/3)", walletcore.SyncStatusInProgress)
		},
		OnRescanningBlocks: func(percentageProgress int64) {
			updateStatus(fmt.Sprintf("Blockchain sync in progress. Rescanning blocks (3/3): %d%%", percentageProgress), walletcore.SyncStatusInProgress)
		},
	}, false)

	if err != nil {
		updateStatus(fmt.Sprintf("Blockchain sync failed to start. %s", err.Error()), walletcore.SyncStatusError)
	}
}

func (b *Blockchain) updateStatus(report string, status walletcore.SyncStatus) {
	b.Lock()
	b._status = status
	b._report = report
	b.Unlock()
}

func (b *Blockchain) status() walletcore.SyncStatus {
	b.RLock()
	defer b.RUnlock()
	return b._status
}

func (b *Blockchain) report() string {
	b.RLock()
	defer b.RUnlock()
	return b._report
}
