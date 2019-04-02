package dcrlibwallet

import (
	"fmt"
	"math"
	"time"

	"github.com/decred/dcrwallet/netparams"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/app"
)

type syncListener struct {
	activeNet *netparams.Params
	walletLib *dcrlibwallet.LibWallet
	data *syncData
}

type syncData struct {
	syncInfo *app.SyncInfoPrivate
	syncInfoUpdated func(*app.SyncInfoPrivate)

	syncing 		bool
	headers       fetchHeadersData
}

type fetchHeadersData struct {
	startHeaderHeight int32
	currentHeaderHeight int32
	beginFetchTimeStamp int64
	totalFetchTime  int64
}

const (
	RescanPercentage = 0.1
	DiscoveryPercentage = 0.8
)

var numberOfPeers int32

func NewSyncListener(activeNet *netparams.Params, walletLib *dcrlibwallet.LibWallet, syncInfoUpdated func(*app.SyncInfoPrivate)) *syncListener {
	data := &syncData{
		syncInfo:&app.SyncInfoPrivate{},
		syncInfoUpdated:syncInfoUpdated,
		syncing:true,
		headers: fetchHeadersData{
			beginFetchTimeStamp: -1,
			totalFetchTime:-1,
		},
	}

	return &syncListener{
		activeNet: activeNet,
		walletLib:walletLib,
		data:data,
	}
}

// following functions are used to implement dcrlibwallet.SpvSyncResponse interface
func (listener *syncListener) OnPeerConnected(peerCount int32)    {
	numberOfPeers = peerCount

	syncInfo := listener.data.syncInfo.Read()
	syncInfo.ConnectedPeers = peerCount

	if listener.data.syncing {
		listener.data.syncInfo.Write(syncInfo, app.SyncStatusInProgress)
	} else {
		listener.data.syncInfo.Write(syncInfo, syncInfo.Status)
	}

	// notify interface of update
	listener.data.syncInfoUpdated(listener.data.syncInfo)
}

func (listener *syncListener) OnPeerDisconnected(peerCount int32) {
	numberOfPeers = peerCount

	syncInfo := listener.data.syncInfo.Read()
	syncInfo.ConnectedPeers = peerCount

	if listener.data.syncing {
		listener.data.syncInfo.Write(syncInfo, app.SyncStatusInProgress)
	} else {
		listener.data.syncInfo.Write(syncInfo, syncInfo.Status)
	}

	// notify interface of update
	listener.data.syncInfoUpdated(listener.data.syncInfo)
}

func (listener *syncListener) OnFetchMissingCFilters(missingCFiltersStart, missingCFiltersEnd int32, state string) {
}

func (listener *syncListener) OnFetchedHeaders(fetchedHeadersCount int32, lastHeaderTime int64, state string) {
	if !listener.data.syncing || listener.data.headers.totalFetchTime != -1 {
		// Ignore this call because this function gets called for each peer and
		// we'd want to ignore those calls as far as the wallet is synced ()
		// or headers are completely fetched (listener.data.headers.totalFetchTime != -1)
		return
	}

	netType := listener.activeNet.Params.Name
	bestBlockTimeStamp := listener.walletLib.GetBestBlockTimeStamp()
	bestBlock := listener.walletLib.GetBestBlock()
	estimatedBlocksToFetch := app.EstimateBlocksCount(netType, bestBlockTimeStamp, bestBlock)

	// update sync info
	syncInfo := listener.data.syncInfo.Read()
	syncInfo.CurrentStep = 1

	switch state {
	case dcrlibwallet.START:
		if listener.data.headers.beginFetchTimeStamp != -1 {
			break
		}

		syncEndPoint := int32(estimatedBlocksToFetch) - listener.data.headers.startHeaderHeight
		syncInfo.TotalHeadersToFetch = syncEndPoint

		listener.data.headers.beginFetchTimeStamp = time.Now().Unix()
		listener.data.headers.startHeaderHeight = bestBlock
		listener.data.headers.currentHeaderHeight = listener.data.headers.startHeaderHeight

	case dcrlibwallet.PROGRESS:
		// increment current block height value
		listener.data.headers.currentHeaderHeight += fetchedHeadersCount

		// calculate percentage progress and eta
		totalFetchedHeaders := listener.data.headers.currentHeaderHeight
		if listener.data.headers.startHeaderHeight > 0 {
			totalFetchedHeaders -= listener.data.headers.startHeaderHeight
		}

		syncEndPoint := int32(estimatedBlocksToFetch) - listener.data.headers.startHeaderHeight
		headersFetchingRate := float64(totalFetchedHeaders) / float64(syncEndPoint)

		timeTakenSoFar := time.Now().Unix() - listener.data.headers.beginFetchTimeStamp
		estimatedTotalHeadersFetchTime := math.Round(float64(timeTakenSoFar) / headersFetchingRate)

		// 10% of estimated fetch time is used for estimating rescan time while 80% is used for estimating address discovery time
		estimatedRescanTime := estimatedTotalHeadersFetchTime * RescanPercentage
		estimatedDiscoveryTime := estimatedTotalHeadersFetchTime * DiscoveryPercentage
		estimatedTotalSyncTime := estimatedTotalHeadersFetchTime + estimatedRescanTime + estimatedDiscoveryTime

		totalTimeRemaining := (int64(estimatedTotalSyncTime) - timeTakenSoFar) / 60
		totalSyncProgress := (float64(timeTakenSoFar) / float64(estimatedTotalSyncTime)) * 100.0

		syncInfo.FetchedHeadersCount = totalFetchedHeaders
		syncInfo.TotalHeadersToFetch = syncEndPoint
		syncInfo.HeadersFetchProgress = int32(math.Round(headersFetchingRate * 100))
		syncInfo.TotalTimeRemaining = fmt.Sprintf("%d min", totalTimeRemaining)
		syncInfo.TotalSyncProgress = int32(math.Round(totalSyncProgress))

		// calculate block header time difference
		hoursBehind := float64(time.Now().Unix() - lastHeaderTime) / 60
		daysBehind := int(math.Round(hoursBehind / 24))
		if daysBehind < 1 {
			syncInfo.DaysBehind = "<1 day"
		} else if daysBehind == 1 {
			syncInfo.DaysBehind = "1 day"
		} else {
			syncInfo.DaysBehind = fmt.Sprintf("%d days", daysBehind)
		}

	case dcrlibwallet.FINISH:
		listener.data.headers.totalFetchTime = time.Now().Unix() - listener.data.headers.beginFetchTimeStamp
		listener.data.headers.startHeaderHeight = -1
		listener.data.headers.currentHeaderHeight = -1
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
