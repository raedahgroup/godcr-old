package dcrlibwallet

import (
	"fmt"
	"time"

	"github.com/decred/dcrwallet/netparams"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletmediums"
)

type syncListener struct {
	data *spvSyncData
}

type spvSyncData struct {
	activeNet *netparams.Params
	walletLib *dcrlibwallet.LibWallet
	syncInfo *app.SyncInfoPrivate

	beginFetchHeaderTimeStamp int64
	syncStartPoint int32
	syncEndPoint   int32
}

func NewSyncListener(activeNet *netparams.Params, walletLib *dcrlibwallet.LibWallet,
	syncInfoUpdated func(*app.SyncInfoPrivate)) *syncListener {

	data := &spvSyncData{
		activeNet: activeNet,
		walletLib:walletLib,
		syncInfo:&app.SyncInfoPrivate{},
		beginFetchHeaderTimeStamp: -1,
	}

	return &syncListener{data:data}
}

func (data *spvSyncData) updateSyncInfo(status app.SyncStatus) {

}

// following functions are used to implement dcrlibwallet.SpvSyncResponse interface
func (listener *syncListener) OnPeerConnected(peerCount int32)    {
	numberOfPeers = peerCount
	syncInfo := listener.data.syncInfo.Read()
	syncInfo.ConnectedPeers = peerCount
	listener.data.syncInfo.Write(syncInfo, app.SyncStatusInProgress)
}

func (listener *syncListener) OnPeerDisconnected(peerCount int32) {
	numberOfPeers = peerCount
	syncInfo := listener.data.syncInfo.Read()
	syncInfo.ConnectedPeers = peerCount
	listener.data.syncInfo.Write(syncInfo, app.SyncStatusInProgress)
}

func (listener *syncListener) OnFetchMissingCFilters(missingCFitlersStart, missingCFitlersEnd int32, state string) {
}

func (listener *syncListener) OnFetchedHeaders(_ int32, lastHeaderTime int64, state string) {
	netType := listener.data.activeNet.Params.Name
	bestBlockTimeStamp := listener.data.walletLib.GetBestBlockTimeStamp()
	estimatedBlocksToFetch := app.EstimateBlocksCount(netType, bestBlockTimeStamp, lastHeaderTime)

	syncStartPoint = listener.data.walletLib.GetBestBlock()
	syncEndPoint = int32(estimatedBlocksToFetch) - syncStartPoint

	switch state {
	case "start":
		if beginFetchHeaderTimeStamp != -1 {
			break
		}
		beginFetchHeaderTimeStamp = time.Now().Unix()


	case "progress":
		response.calculateProgress(lastHeaderTime)
	}
}

func (response *SpvSyncResponse) OnDiscoveredAddresses(state string) {
	response.listener.OnDiscoveredAddress(state)
}

func (response *SpvSyncResponse) OnRescan(rescannedThrough int32, state string) {
	if state == "progress" {
		bestBlock := int64(response.walletLib.GetBestBlock())
		scannedPercentage := int64(rescannedThrough) / bestBlock * 100
		response.listener.OnRescanningBlocks(scannedPercentage)
	}
}

func (response *SpvSyncResponse) OnSynced(synced bool) {
	var err error
	if !synced {
		err = fmt.Errorf("Sync failed")
	}
	response.listener.SyncEnded(err)
}

func (response *SpvSyncResponse) OnSyncError(code int, err error) {
	e := fmt.Errorf("Code: %d, Error: %s", code, err.Error())
	response.listener.SyncEnded(e)
}

func (response *SpvSyncResponse) calculateProgress(lastHeaderTime int64) {
	bestBlock := int64(response.walletLib.GetBestBlock())
	fetchedPercentage := walletmediums.CalculateBlockSyncProgress(response.activeNet.Params.Name, bestBlock, lastHeaderTime)
	response.listener.OnHeadersFetched(fetchedPercentage)
}
