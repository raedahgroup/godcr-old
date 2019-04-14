package dcrwalletrpc

import (
	"fmt"
	"github.com/decred/dcrwallet/rpc/walletrpc"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletmediums"
	"github.com/decred/dcrwallet/netparams"
	"github.com/raedahgroup/godcr/app/sync"
)

var numberOfPeers int32

type syncListener struct {
	activeNet *netparams.Params
	spvSyncClient    walletrpc.WalletLoaderService_SpvSyncClient

	netType         string
	showLog         bool
	syncing         bool
	syncInfoUpdated func(*sync.PrivateInfo)

	privateSyncInfo *sync.PrivateInfo
	headersData     *sync.FetchHeadersData

	addressDiscoveryCompleted chan bool

	rescanStartTime int64
}

func (listener syncListener) streamBlockChainSyncUpdates(showLog bool) {
	logUpdate := func(format string, values ...interface{}) {
		if showLog {
			fmt.Printf(format, values)
		}
	}

	for {
		syncUpdate, err := listener.spvSyncClient.Recv()

		// read current sync info to get ready for updating
		syncInfo := listener.privateSyncInfo.Read()

		if err != nil {
			syncInfo.Error = err.Error()
			listener.privateSyncInfo.Write(syncInfo, sync.StatusError)
		} else if syncUpdate.Synced {
			listener.privateSyncInfo.Write(syncInfo, sync.StatusSuccess)
		}
		if err != nil || syncUpdate.Synced {
			listener.syncing = false
			syncInfo.Done = true

			// notify interface of sync update and exit loop
			listener.syncInfoUpdated(listener.privateSyncInfo)
			break
		}

		switch syncUpdate.NotificationType {
		case walletrpc.SyncNotificationType_FETCHED_HEADERS_STARTED:
			logUpdate("Blockchain sync in progress. Start fetching headers (1/3)")

		case walletrpc.SyncNotificationType_FETCHED_HEADERS_PROGRESS:
			fetchedPercentage := walletmediums.CalculateBlockSyncProgress(listener.netType, listener.bestBlock, syncUpdate.FetchHeaders.LastHeaderTime)
			logUpdate("Blockchain sync in progress. Fetching headers (1/3): %d%%", fetchedPercentage)
			listener.listener.OnHeadersFetched(fetchedPercentage)

		case walletrpc.SyncNotificationType_FETCHED_HEADERS_FINISHED:
			logUpdate("Blockchain sync in progress. Done fetching headers (1/3)")

		case walletrpc.SyncNotificationType_DISCOVER_ADDRESSES_STARTED:
			logUpdate("Blockchain sync in progress. Start discovering addresses (2/3)")

		case walletrpc.SyncNotificationType_DISCOVER_ADDRESSES_FINISHED:
			logUpdate("Blockchain sync in progress. Finished discovering addresses (2/3)")

		case walletrpc.SyncNotificationType_RESCAN_STARTED:
			logUpdate("Blockchain sync in progress. Start rescanning blocks (3/3)")

		case walletrpc.SyncNotificationType_RESCAN_PROGRESS:
			scannedPercentage := int64(syncUpdate.RescanProgress.RescannedThrough) / listener.bestBlock * 100
			logUpdate("Blockchain sync in progress. Rescanning blocks (3/3): %d%%", scannedPercentage)
			listener.listener.OnHeadersFetched(scannedPercentage)

		case walletrpc.SyncNotificationType_RESCAN_FINISHED:
			logUpdate("Blockchain sync in progress. Done rescanning blocks (3/3)")

		case walletrpc.SyncNotificationType_PEER_CONNECTED:
			listener.listener.OnPeersUpdated(syncUpdate.PeerInformation.PeerCount)
			logUpdate("New peer %listener. Connected to %d peers", syncUpdate.PeerInformation.Address, syncUpdate.PeerInformation.PeerCount)
			// numberOfPeers needs to be updated before send OnPeerConnected
			numberOfPeers = syncUpdate.PeerInformation.PeerCount
			//listener.listener.OnPeerConnected(syncUpdate.PeerInformation.PeerCount)

		case walletrpc.SyncNotificationType_PEER_DISCONNECTED:
			listener.listener.OnPeersUpdated(syncUpdate.PeerInformation.PeerCount)
			logUpdate("Peer disconnected %listener. Connected to %d peers", syncUpdate.PeerInformation.Address, syncUpdate.PeerInformation.PeerCount)
			numberOfPeers = syncUpdate.PeerInformation.PeerCount
			//listener.listener.OnPeerDisconnected(syncUpdate.PeerInformation.PeerCount)
		}
	}
}
