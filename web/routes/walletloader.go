package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/raedahgroup/godcr/app/sync"
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
		syncInfo := routes.syncPrivateInfo.Read()

		switch syncInfo.Status {
		case sync.StatusSuccess:
			next.ServeHTTP(res, req)
		case sync.StatusNotStarted:
			errMsg = "Cannot display page. Blockchain hasn't been synced"
		case sync.StatusInProgress:
			var syncInfoMap map[string]interface{}
			syncInfoBytes, _ := json.Marshal(syncInfo)

			jsonDecoder := json.NewDecoder(bytes.NewReader(syncInfoBytes))
			jsonDecoder.UseNumber()

			err := jsonDecoder.Decode(&syncInfoMap)
			if err != nil {
				errMsg = err.Error()
			} else {
				routes.renderSyncPage(syncInfoMap, res)
			}
		case sync.StatusError:
			errMsg = fmt.Sprintf("Cannot display page. Following error occured during sync: %s", syncInfo.Error)
		default:
			errMsg = "Cannot display page. Blockchain sync status cannot be determined"
		}
	})
}

func (routes *Routes) syncBlockchain() {
<<<<<<< HEAD
	err := routes.walletMiddleware.SyncBlockChain(false, func(syncInfo *app.SyncInfoPrivate) {
		currentInfo := routes.syncInfo.Read()
		newInfo := routes.syncInfo.Read()
		if currentInfo.ConnectedPeers != newInfo.ConnectedPeers || !currentInfo.Done && newInfo.Done {
			routes.sendWsConnectionInfoUpdate()
		}
		routes.syncInfo = syncInfo
=======
	err := routes.walletMiddleware.SyncBlockChain(false, func(syncPrivateInfo *sync.PrivateInfo) {
		routes.syncPrivateInfo = syncPrivateInfo
>>>>>>> refactor
	})

	// update sync status
	syncInfo := routes.syncPrivateInfo.Read()

	if err != nil {
		syncInfo.Error = err.Error()
		syncInfo.Done = true
		routes.syncPrivateInfo.Write(syncInfo, sync.StatusError)
	} else {
		routes.syncPrivateInfo.Write(syncInfo, sync.StatusInProgress)
	}
}
