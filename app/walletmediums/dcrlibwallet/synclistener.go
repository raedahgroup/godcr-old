package dcrlibwallet

import (
	"fmt"
	"github.com/raedahgroup/godcr/app/sync"
	"math"
	"time"

	"github.com/decred/dcrwallet/netparams"
	"github.com/raedahgroup/dcrlibwallet"
)

type syncListener struct {
	activeNet *netparams.Params
	walletLib *dcrlibwallet.LibWallet

	netType         string
	showLog         bool
	syncing         bool
	syncInfoUpdated func(*sync.PrivateInfo)

	privateSyncInfo *sync.PrivateInfo
	headersData     *sync.FetchHeadersData

	addressDiscoveryCompleted chan bool

	rescanStartTime int64
}

var numberOfPeers int32

func NewSyncListener(activeNet *netparams.Params, walletLib *dcrlibwallet.LibWallet, showLog bool,
	syncInfoUpdated func(*sync.PrivateInfo)) *syncListener {

	return &syncListener{
		activeNet: activeNet,
		walletLib: walletLib,

		netType:         activeNet.Params.Name,
		showLog:         showLog,
		syncing:         true,
		syncInfoUpdated: syncInfoUpdated,

		privateSyncInfo: sync.NewPrivateInfo(),
		headersData: &sync.FetchHeadersData{
			BeginFetchTimeStamp: -1,
		},
	}
}

// following functions are used to implement dcrlibwallet.SpvSyncResponse interface
func (listener *syncListener) OnPeerConnected(peerCount int32) {
	numberOfPeers = peerCount

	syncInfo := listener.privateSyncInfo.Read()
	syncInfo.ConnectedPeers = peerCount
	listener.privateSyncInfo.Write(syncInfo, syncInfo.Status)

	// notify interface of update
	listener.syncInfoUpdated(listener.privateSyncInfo)

	if listener.showLog && listener.syncing {
		if peerCount == 1 {
			fmt.Printf("Connected to %d peer on %s.\n", peerCount, listener.netType)
		} else {
			fmt.Printf("Connected to %d peers on %s.\n", peerCount, listener.netType)
		}
	}
}

func (listener *syncListener) OnPeerDisconnected(peerCount int32) {
	numberOfPeers = peerCount

	syncInfo := listener.privateSyncInfo.Read()
	syncInfo.ConnectedPeers = peerCount
	listener.privateSyncInfo.Write(syncInfo, syncInfo.Status)

	// notify interface of update
	listener.syncInfoUpdated(listener.privateSyncInfo)

	if listener.showLog && listener.syncing {
		if peerCount == 1 {
			fmt.Printf("Connected to %d peer on %s.\n", peerCount, listener.netType)
		} else {
			fmt.Printf("Connected to %d peers on %s.\n", peerCount, listener.netType)
		}
	}
}

func (listener *syncListener) OnFetchMissingCFilters(missingCFiltersStart, missingCFiltersEnd int32, state string) {
}

func (listener *syncListener) OnFetchedHeaders(fetchedHeadersCount int32, lastHeaderTime int64, state string) {
	syncInfo := listener.privateSyncInfo.Read()

	if !listener.syncing || syncInfo.HeadersFetchTimeTaken != -1 {
		// Ignore this call because this function gets called for each peer and
		// we'd want to ignore those calls as far as the wallet is synced (i.e. !listener.syncing)
		// or headers are completely fetched (i.e. syncInfo.HeadersFetchTimeTaken != -1)
		return
	}

	bestBlockTimeStamp := listener.walletLib.GetBestBlockTimeStamp()
	bestBlock := listener.walletLib.GetBestBlock()
	estimatedFinalBlockHeight := sync.EstimateFinalBlockHeight(listener.netType, bestBlockTimeStamp, bestBlock)

	switch state {
	case dcrlibwallet.START:
		if listener.headersData.BeginFetchTimeStamp != -1 {
			break
		}

		listener.headersData.BeginFetchTimeStamp = time.Now().Unix()
		listener.headersData.StartHeaderHeight = bestBlock
		listener.headersData.CurrentHeaderHeight = listener.headersData.StartHeaderHeight

		syncInfo.TotalHeadersToFetch = int32(estimatedFinalBlockHeight) - listener.headersData.StartHeaderHeight
		syncInfo.CurrentStep = 1

		if listener.showLog {
			fmt.Printf("Step 1 of 3 - fetching %d block headers.\n", syncInfo.TotalHeadersToFetch)
		}

	case dcrlibwallet.PROGRESS:
		headersFetchReport := sync.FetchHeadersProgressReport{
			FetchedHeadersCount:       fetchedHeadersCount,
			LastHeaderTime:            lastHeaderTime,
			EstimatedFinalBlockHeight: estimatedFinalBlockHeight,
		}
		sync.UpdateFetchHeadersProgress(syncInfo, listener.headersData, headersFetchReport)

		if listener.showLog {
			fmt.Printf("Syncing %d%%, %s remaining, fetched %d of %d block headers, %s behind.\n",
				syncInfo.TotalSyncProgress, syncInfo.TotalTimeRemaining,
				syncInfo.FetchedHeadersCount, syncInfo.TotalHeadersToFetch,
				syncInfo.DaysBehind)
		}

	case dcrlibwallet.FINISH:
		syncInfo.HeadersFetchTimeTaken = time.Now().Unix() - listener.headersData.BeginFetchTimeStamp
		syncInfo.TotalHeadersToFetch = -1

		listener.headersData.StartHeaderHeight = -1
		listener.headersData.CurrentHeaderHeight = -1

		if listener.showLog {
			fmt.Println("Fetch headers completed.")
		}
	}

	listener.privateSyncInfo.Write(syncInfo, sync.StatusInProgress)

	// notify ui of updated sync info
	listener.syncInfoUpdated(listener.privateSyncInfo)
}

func (listener *syncListener) OnDiscoveredAddresses(state string) {
	if state == dcrlibwallet.START && listener.addressDiscoveryCompleted == nil {
		if listener.showLog && listener.syncing {
			fmt.Println("Step 2 of 3 - discovering used addresses.")
		}
		listener.addressDiscoveryCompleted = sync.UpdateAddressDiscoveryProgress(listener.privateSyncInfo, listener.showLog,
			listener.syncInfoUpdated)
	} else {
		close(listener.addressDiscoveryCompleted)
		listener.addressDiscoveryCompleted = nil
	}
}

func (listener *syncListener) OnRescan(rescannedThrough int32, state string) {
	if listener.addressDiscoveryCompleted != nil {
		close(listener.addressDiscoveryCompleted)
		listener.addressDiscoveryCompleted = nil
	}

	syncInfo := listener.privateSyncInfo.Read()

	if syncInfo.TotalHeadersToFetch == -1 {
		syncInfo.TotalHeadersToFetch = listener.walletLib.GetBestBlock()
	}

	switch state {
	case dcrlibwallet.START:
		listener.rescanStartTime = time.Now().Unix()
		syncInfo.TotalHeadersToFetch = listener.walletLib.GetBestBlock()
		syncInfo.CurrentStep = 3

		if listener.showLog && listener.syncing {
			fmt.Println("Step 3 of 3 - Rescanning blocks")
		}

	case dcrlibwallet.PROGRESS:
		elapsedRescanTime := time.Now().Unix() - listener.rescanStartTime
		totalElapsedTime := syncInfo.HeadersFetchTimeTaken + syncInfo.TotalDiscoveryTime + elapsedRescanTime

		rescanRate := float64(rescannedThrough) / float64(syncInfo.TotalHeadersToFetch)
		estimatedTotalRescanTime := float64(elapsedRescanTime) / rescanRate
		estimatedTotalSyncTime := syncInfo.HeadersFetchTimeTaken + syncInfo.TotalDiscoveryTime + int64(math.Round(estimatedTotalRescanTime))

		totalProgress := (float64(totalElapsedTime) / float64(estimatedTotalSyncTime)) * 100

		// do not update total time taken and total progress percent if elapsedRescanTime is 0
		// because the estimatedTotalRescanTime will be inaccurate (also 0)
		// which will make the estimatedTotalSyncTime equal to totalElapsedTime
		// giving the wrong impression that the process is complete
		if elapsedRescanTime > 0 {
			syncInfo.TotalTimeRemaining = sync.CalculateTotalTimeRemaining(estimatedTotalRescanTime - float64(elapsedRescanTime))
			syncInfo.TotalSyncProgress = int32(math.Round(totalProgress))
		}

		syncInfo.RescanProgress = int32(math.Round(rescanRate * 100))
		syncInfo.CurrentRescanHeight = rescannedThrough

		if listener.showLog && listener.syncing {
			fmt.Printf("Syncing %d%%, %s remaining, scanning %d of %d block headers.\n",
				syncInfo.TotalSyncProgress, syncInfo.TotalTimeRemaining,
				syncInfo.CurrentRescanHeight, syncInfo.TotalHeadersToFetch)
		}

	case dcrlibwallet.FINISH:
		if listener.showLog && listener.syncing {
			fmt.Println("Block headers scan complete.")
		}
	}

	listener.privateSyncInfo.Write(syncInfo, sync.StatusInProgress)

	// notify ui of updated sync info
	listener.syncInfoUpdated(listener.privateSyncInfo)
}

func (listener *syncListener) OnIndexTransactions(totalIndex int32) {}

func (listener *syncListener) OnSynced(synced bool) {
	if !listener.syncing {
		// ignore subsequent updates
		return
	}

	syncInfo := listener.privateSyncInfo.Read()
	syncInfo.Done = true
	listener.syncing = false

	if !synced {
		syncInfo.Error = "Sync failed or canceled"
		listener.privateSyncInfo.Write(syncInfo, sync.StatusError)
	} else {
		listener.privateSyncInfo.Write(syncInfo, sync.StatusSuccess)
	}

	// notify interface of update
	listener.syncInfoUpdated(listener.privateSyncInfo)
}

// todo sync may not have ended
func (listener *syncListener) OnSyncError(code int, err error) {
	if !listener.syncing {
		// ignore subsequent updates
		return
	}

	syncInfo := listener.privateSyncInfo.Read()
	syncInfo.Done = true
	listener.syncing = false

	syncInfo.Error = fmt.Sprintf("Code: %d, Error: %s", code, err.Error())
	listener.privateSyncInfo.Write(syncInfo, sync.StatusError)

	// notify interface of update
	listener.syncInfoUpdated(listener.privateSyncInfo)
}
