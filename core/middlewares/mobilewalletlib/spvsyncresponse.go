package mobilewalletlib

import (
	"fmt"
	"time"

	"github.com/decred/dcrwallet/netparams"
	"github.com/raedahgroup/dcrcli/core"
	"github.com/raedahgroup/mobilewallet"
)

const (
	MainNetTargetTimePerBlock = 300
	TestNetTargetTimePerBlock = 120
)

type SpvSyncResponse struct {
	activeNet *netparams.Params
	walletLib *mobilewallet.LibWallet
	listener  *core.BlockChainSyncListener
}

// following functions are used to implement mobilewallet.SpvSyncResponse interface
func (response SpvSyncResponse) OnPeerConnected(peerCount int32)    {}
func (response SpvSyncResponse) OnPeerDisconnected(peerCount int32) {}
func (response SpvSyncResponse) OnFetchMissingCFilters(missingCFitlersStart, missingCFitlersEnd int32, state string) {}
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
	e := fmt.Errorf("Error syncing. Code: %d, Error: %s", code, err.Error())
	response.listener.SyncEnded(e)
}

func (response SpvSyncResponse) calculateProgress(lastHeaderTime int64) {
	var targetTimePerBlock int64
	if response.activeNet.Params.Name == "mainnet" {
		targetTimePerBlock = MainNetTargetTimePerBlock
	} else {
		targetTimePerBlock = TestNetTargetTimePerBlock
	}

	bestBlock := int64(response.walletLib.GetBestBlock())
	estimatedBlocks := ((time.Now().Unix() - lastHeaderTime) / targetTimePerBlock) + bestBlock
	fetchedPercentage := bestBlock / estimatedBlocks * 100

	if fetchedPercentage >= 100 {
		fetchedPercentage = 100
	}

	response.listener.OnHeadersFetched(fetchedPercentage)
}
