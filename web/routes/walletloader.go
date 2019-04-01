package routes

import (
	"fmt"
	"github.com/raedahgroup/godcr/app"
	"net/http"
)

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

		if !routes.walletExists {
			routes.renderNoWalletError(res)
			return
		}

		if !routes.walletMiddleware.IsWalletOpen() {
			errMsg = "Wallet is not open. Restart the server"
			return
		}

		// wallet is open, check if blockchain is synced
		syncInfo := routes.syncInfo.Read()

		switch syncInfo.Status {
		case app.SyncStatusSuccess:
			next.ServeHTTP(res, req)
		case app.SyncStatusNotStarted:
			errMsg = "Cannot display page. Blockchain hasn't been synced"
		case app.SyncStatusInProgress:
			routes.renderSyncPage(syncInfo, res)
		case app.SyncStatusError:
			errMsg = fmt.Sprintf("Cannot display page. Following error occured during sync: %s", syncInfo.Error)
		default:
			errMsg = "Cannot display page. Blockchain sync status cannot be determined"
		}
	})
}

func (routes *Routes) syncBlockchain() {
	syncInfo := routes.syncInfo.Read()
	updateSyncInfo := func(status app.SyncStatus) {
		routes.syncInfo.Write(syncInfo, status)
	}

	err := routes.walletMiddleware.SyncBlockChain(&app.BlockChainSyncListener{
		SyncStarted: func() {
			updateSyncInfo(app.SyncStatusInProgress)
		},
		SyncEnded: func(err error) {
			routes.sendWsConnectionInfoUpdate()
			if syncInfo.Done {
				// ignore subsequent sync ended updates after the sync is already set to done
				return
			}

			syncInfo.Done = true
			if err != nil {
				syncInfo.Error = err.Error()
				updateSyncInfo(app.SyncStatusError)
			} else {
				updateSyncInfo(app.SyncStatusSuccess)
			}
		},
		OnHeadersFetched: func(percentageProgress int64) {
			syncInfo.CurrentBlockHeight = int(percentageProgress)
			updateSyncInfo(app.SyncStatusInProgress)
		},
		OnDiscoveredAddress: func(_ string) {
			updateSyncInfo(app.SyncStatusInProgress)
		},
		OnRescanningBlocks: func(percentageProgress int64) {
			updateSyncInfo(app.SyncStatusInProgress)
		},
		OnPeersUpdated: func(peerCount int32) {
			syncInfo.ConnectedPeers = peerCount
			updateSyncInfo(app.SyncStatusInProgress)
			routes.sendWsConnectionInfoUpdate()
		},
	}, false)

	if err != nil {
		syncInfo.Error = err.Error()
		updateSyncInfo(app.SyncStatusError)
	}
}
