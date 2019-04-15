package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/raedahgroup/dcrlibwallet/blockchainsync"
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
		syncInfo := routes.privateSyncInfo.Read()

		switch syncInfo.Status {
		case blockchainsync.StatusSuccess:
			next.ServeHTTP(res, req)
		case blockchainsync.StatusNotStarted:
			errMsg = "Cannot display page. Blockchain hasn't been synced"
		case blockchainsync.StatusInProgress:
			syncInfoMap, err := routes.prepareSyncInfoMap()
			if err != nil {
				errMsg = fmt.Sprintf("Cannot load sync progress page: %s", err.Error())
			} else {
				routes.renderSyncPage(syncInfoMap, res)
			}
		case blockchainsync.StatusError:
			errMsg = fmt.Sprintf("Cannot display page. Following error occured during sync: %s", syncInfo.Error)
		default:
			errMsg = "Cannot display page. Blockchain sync status cannot be determined"
		}
	})
}

func (routes *Routes) syncBlockChain() {
	err := routes.walletMiddleware.SyncBlockChain(false, func(privateSyncInfo *blockchainsync.PrivateSyncInfo, updatedSection string) {
		routes.privateSyncInfo = privateSyncInfo
		routes.sendWsSyncProgress()
		routes.sendWsConnectionInfoUpdate()
	})

	// update sync status
	syncInfo := routes.privateSyncInfo.Read()	

	if err != nil {
		syncInfo.Error = err.Error()
		syncInfo.Done = true
		routes.privateSyncInfo.Write(syncInfo, blockchainsync.StatusError)
	} else {
		routes.privateSyncInfo.Write(syncInfo, blockchainsync.StatusInProgress)
	}
}

func (routes *Routes) prepareSyncInfoMap() (map[string]interface{}, error) {
	syncInfo := routes.privateSyncInfo.Read()
	var syncInfoMap map[string]interface{}

	syncInfoBytes, _ := json.Marshal(syncInfo)
	jsonDecoder := json.NewDecoder(bytes.NewReader(syncInfoBytes))
	jsonDecoder.UseNumber()

	err := jsonDecoder.Decode(&syncInfoMap)
	if err != nil {
		return nil, err
	}

	syncInfoMap["NetworkType"] = routes.walletMiddleware.NetType()

	if syncInfo.CurrentStep == 2 {
		// check account discovery progress percentage
		if syncInfo.AddressDiscoveryProgress > 100 {
			syncInfoMap["AddressDiscoveryProgress"] = fmt.Sprintf("%d%% (over)", syncInfo.AddressDiscoveryProgress)
		} else {
			syncInfoMap["AddressDiscoveryProgress"] = fmt.Sprintf("%d%%", syncInfo.AddressDiscoveryProgress)
		}
	}

	if syncInfo.ConnectedPeers == 1 {
		syncInfoMap["ConnectedPeers"] = fmt.Sprintf("%d peer", syncInfo.ConnectedPeers)
	} else {
		syncInfoMap["ConnectedPeers"] = fmt.Sprintf("%d peers", syncInfo.ConnectedPeers)
	}

	return syncInfoMap, nil
}
