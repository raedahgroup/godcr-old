package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/raedahgroup/dcrlibwallet/defaultsynclistener"
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

		// todo: main.go now requires that the user select a wallet or create one before launching interfaces, so need for this check
		//if !routes.walletExists {
		//	routes.renderNoWalletError(res)
		//	return
		//}

		if !routes.walletMiddleware.IsWalletOpen() {
			errMsg = "Wallet is not open. Restart the server"
			return
		}

		// wallet is open, check if blockchain is synced
		syncProgressReport := routes.syncProgressReport.Read()

		if syncProgressReport.Done {
			// ignore any other reported sync progress
			next.ServeHTTP(res, req)
			return
		}

		syncProgressReport.Status = defaultsynclistener.SyncStatusSuccess

		switch syncProgressReport.Status {
		case defaultsynclistener.SyncStatusSuccess:
			next.ServeHTTP(res, req)
		case defaultsynclistener.SyncStatusInProgress:
			syncInfoMap, err := routes.prepareSyncInfoMap()
			if err != nil {
				errMsg = fmt.Sprintf("Cannot load sync progress page: %s", err.Error())
			} else {
				routes.renderSyncPage(syncInfoMap, res)
			}
		case defaultsynclistener.SyncStatusError:
			errMsg = fmt.Sprintf("Cannot display page. Following error occured during sync: %s", syncProgressReport.Error)
		default:
			errMsg = "Cannot display page. Blockchain sync status cannot be determined"
		}
	})
}

func (routes *Routes) syncBlockChain() {
	routes.walletMiddleware.SyncBlockChain(false, func(report *defaultsynclistener.ProgressReport) {
		routes.syncProgressReport = report
		routes.sendWsSyncProgress()
		routes.sendWsConnectionInfoUpdate()
	})
}

func (routes *Routes) prepareSyncInfoMap() (map[string]interface{}, error) {
	syncInfo := routes.syncProgressReport.Read()
	var syncInfoMap map[string]interface{}

	syncInfoBytes, _ := json.Marshal(syncInfo)
	jsonDecoder := json.NewDecoder(bytes.NewReader(syncInfoBytes))
	jsonDecoder.UseNumber()

	err := jsonDecoder.Decode(&syncInfoMap)
	if err != nil {
		return nil, err
	}

	syncInfoMap["networkType"] = routes.walletMiddleware.NetType()

	if syncInfo.CurrentStep == defaultsynclistener.DiscoveringUsedAddresses {
		// check account discovery progress percentage
		if syncInfo.AddressDiscoveryProgress > 100 {
			syncInfoMap["addressDiscoveryProgress"] = fmt.Sprintf("%d%% (over)", syncInfo.AddressDiscoveryProgress)
		} else {
			syncInfoMap["addressDiscoveryProgress"] = fmt.Sprintf("%d%%", syncInfo.AddressDiscoveryProgress)
		}
	}

	if syncInfo.ConnectedPeers == 1 {
		syncInfoMap["connectedPeers"] = fmt.Sprintf("%d peer", syncInfo.ConnectedPeers)
	} else {
		syncInfoMap["connectedPeers"] = fmt.Sprintf("%d peers", syncInfo.ConnectedPeers)
	}

	return syncInfoMap, nil
}
