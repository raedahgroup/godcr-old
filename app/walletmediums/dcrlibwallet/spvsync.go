package dcrlibwallet

import (
	"fmt"

	"github.com/decred/dcrwallet/netparams"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletmediums"
)

type SpvSyncResponse struct {
	activeNet *netparams.Params
	walletLib *dcrlibwallet.LibWallet
	listener  *app.BlockChainSyncListener
}

// following functions are used to implement dcrlibwallet.SpvSyncResponse interface
func (response SpvSyncResponse) OnPeerConnected(peerCount int32)    {}
func (response SpvSyncResponse) OnPeerDisconnected(peerCount int32) {}
func (response SpvSyncResponse) OnFetchMissingCFilters(missingCFitlersStart, missingCFitlersEnd int32, state string) {
}
func (response SpvSyncResponse) OnFetchedHeaders(_ int32, lastHeaderTime int64, state string) {
	if state == "progress" {
		response.calculateProgress(lastHeaderTime)
	}
}
func (response SpvSyncResponse) OnDiscoveredAddresses(state string) {
	response.listener.OnDiscoveredAddress(state)
}
func (response SpvSyncResponse) OnRescan(rescannedThrough int32, state string) {
	if state == "progress" {
		bestBlock := int64(response.walletLib.GetBestBlock())
		scannedPercentage := int64(rescannedThrough) / bestBlock * 100
		response.listener.OnRescanningBlocks(scannedPercentage)
	}
}
func (response SpvSyncResponse) OnSynced(synced bool) {
	var err error
	if !synced {
		err = fmt.Errorf("Sync failed")
	}
	response.listener.SyncEnded(err)
}
func (response SpvSyncResponse) OnSyncError(code int, err error) {
	e := fmt.Errorf("Code: %d, Error: %s", code, err.Error())
	response.listener.SyncEnded(e)
}

func (response SpvSyncResponse) calculateProgress(lastHeaderTime int64) {
	bestBlock := int64(response.walletLib.GetBestBlock())
	fetchedPercentage := walletmediums.CalculateBlockSyncProgress(response.activeNet.Params.Name, bestBlock, lastHeaderTime)
	response.listener.OnHeadersFetched(fetchedPercentage)
}
