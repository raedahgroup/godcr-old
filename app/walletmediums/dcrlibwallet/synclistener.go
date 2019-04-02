package dcrlibwallet

import (
	"fmt"
	"time"

	"github.com/decred/dcrwallet/netparams"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/app"
)

type syncListener struct {
	activeNet *netparams.Params
	walletLib *dcrlibwallet.LibWallet
	showLog bool
	data *syncData
}

type syncData struct {
	netType string
	syncInfo *app.SyncInfoPrivate
	syncInfoUpdated func(*app.SyncInfoPrivate)

	syncing 		bool
	headers       *app.FetchHeadersData
}

var numberOfPeers int32

func NewSyncListener(activeNet *netparams.Params, walletLib *dcrlibwallet.LibWallet, showLog bool,
	syncInfoUpdated func(*app.SyncInfoPrivate)) *syncListener {

	data := &syncData{
		netType: activeNet.Params.Name,
		syncInfo: app.NewSyncInfo(),
		syncInfoUpdated:syncInfoUpdated,
		syncing:true,
		headers: &app.FetchHeadersData{
			BeginFetchTimeStamp: -1,
		},
	}

	return &syncListener{
		activeNet: activeNet,
		walletLib:walletLib,
		showLog:showLog,
		data:data,
	}
}

// following functions are used to implement dcrlibwallet.SpvSyncResponse interface
func (listener *syncListener) OnPeerConnected(peerCount int32)    {
	numberOfPeers = peerCount
	if listener.showLog {
		if peerCount == 1 {
			fmt.Printf("Connected to %d peer on %s.\n", peerCount, listener.data.netType)
		} else {
			fmt.Printf("Connected to %d peers on %s.\n", peerCount, listener.data.netType)
		}
	}

	syncInfo := listener.data.syncInfo.Read()
	syncInfo.ConnectedPeers = peerCount
	listener.data.syncInfo.Write(syncInfo, syncInfo.Status)

	// notify interface of update
	listener.data.syncInfoUpdated(listener.data.syncInfo)
}

func (listener *syncListener) OnPeerDisconnected(peerCount int32) {
	numberOfPeers = peerCount
	if listener.showLog {
		if peerCount == 1 {
			fmt.Printf("Connected to %d peer on %s.\n", peerCount, listener.data.netType)
		} else {
			fmt.Printf("Connected to %d peers on %s.\n", peerCount, listener.data.netType)
		}
	}

	syncInfo := listener.data.syncInfo.Read()
	syncInfo.ConnectedPeers = peerCount
	listener.data.syncInfo.Write(syncInfo, syncInfo.Status)

	// notify interface of update
	listener.data.syncInfoUpdated(listener.data.syncInfo)
}

func (listener *syncListener) OnFetchMissingCFilters(missingCFiltersStart, missingCFiltersEnd int32, state string) {
}

func (listener *syncListener) OnFetchedHeaders(fetchedHeadersCount int32, lastHeaderTime int64, state string) {
	syncInfo := listener.data.syncInfo.Read()

	if !listener.data.syncing || syncInfo.HeadersFetchTimeTaken != -1 {
		// Ignore this call because this function gets called for each peer and
		// we'd want to ignore those calls as far as the wallet is synced (i.e. !listener.data.syncing)
		// or headers are completely fetched (i.e. syncInfo.HeadersFetchTimeTaken != -1)
		return
	}

	bestBlockTimeStamp := listener.walletLib.GetBestBlockTimeStamp()
	bestBlock := listener.walletLib.GetBestBlock()
	estimatedFinalBlockHeight := app.EstimateFinalBlockHeight(listener.data.netType, bestBlockTimeStamp, bestBlock)

	switch state {
	case dcrlibwallet.START:
		if listener.data.headers.BeginFetchTimeStamp != -1 {
			break
		}

		listener.data.headers.BeginFetchTimeStamp = time.Now().Unix()
		listener.data.headers.StartHeaderHeight = bestBlock
		listener.data.headers.CurrentHeaderHeight = listener.data.headers.StartHeaderHeight

		syncInfo.TotalHeadersToFetch = int32(estimatedFinalBlockHeight) - listener.data.headers.StartHeaderHeight
		syncInfo.CurrentStep = 1

		if listener.showLog {
			fmt.Printf("Step 1 of 3 - fetching %d block headers.\n", syncInfo.TotalHeadersToFetch)
		}

	case dcrlibwallet.PROGRESS:
		headersFetchReport := app.FetchHeadersProgressReport{
			FetchedHeadersCount:       fetchedHeadersCount,
			LastHeaderTime:            lastHeaderTime,
			EstimatedFinalBlockHeight: estimatedFinalBlockHeight,
		}
		app.UpdateFetchHeadersProgress(listener.data.headers, headersFetchReport, syncInfo)

		if listener.showLog {
			fmt.Printf("Syncing %d%%, %s remaining, fetched %d of %d block headers, %s behind.\n",
				syncInfo.TotalSyncProgress, syncInfo.TotalTimeRemaining,
				syncInfo.FetchedHeadersCount, syncInfo.TotalHeadersToFetch,
				syncInfo.DaysBehind)
		}

	case dcrlibwallet.FINISH:
		syncInfo.HeadersFetchTimeTaken = time.Now().Unix() - listener.data.headers.BeginFetchTimeStamp
		listener.data.headers.StartHeaderHeight = -1
		listener.data.headers.CurrentHeaderHeight = -1

		if listener.showLog {
			fmt.Println("Fetch headers completed.")
		}
	}

	listener.data.syncInfo.Write(syncInfo, app.SyncStatusInProgress)

	// notify ui of updated sync info
	listener.data.syncInfoUpdated(listener.data.syncInfo)
}

func (listener *syncListener) OnDiscoveredAddresses(state string) {
}

func (listener *syncListener) OnRescan(rescannedThrough int32, state string) {
}

func (listener *syncListener) OnSynced(synced bool) {
	if !listener.data.syncing {
		// ignore subsequent updates
		return
	}

	syncInfo := listener.data.syncInfo.Read()
	syncInfo.Done = true
	listener.data.syncing = false

	if !synced {
		syncInfo.Error = "Sync failed or canceled"
		listener.data.syncInfo.Write(syncInfo, app.SyncStatusError)
	} else {
		listener.data.syncInfo.Write(syncInfo, app.SyncStatusSuccess)
	}

	// notify interface of update
	listener.data.syncInfoUpdated(listener.data.syncInfo)
}

func (listener *syncListener) OnSyncError(code int, err error) {
	if !listener.data.syncing {
		// ignore subsequent updates
		return
	}

	syncInfo := listener.data.syncInfo.Read()
	syncInfo.Done = true
	listener.data.syncing = false

	syncInfo.Error = fmt.Sprintf("Code: %d, Error: %s", code, err.Error())
	listener.data.syncInfo.Write(syncInfo, app.SyncStatusError)

	// notify interface of update
	listener.data.syncInfoUpdated(listener.data.syncInfo)
}
