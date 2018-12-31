package dcrwalletrpc

import (
	"fmt"
	"github.com/decred/dcrwallet/rpc/walletrpc"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletmediums"
)

type spvSync struct {
	client    walletrpc.WalletLoaderService_SpvSyncClient
	bestBlock int64
	listener  *app.BlockChainSyncListener
	netType   string
}

func (s spvSync) streamBlockchainSyncUpdates(showLog bool) {
	logUpdate := func(format string, values ...interface{}) {
		if showLog {
			fmt.Printf(format, values)
		}
	}

	s.listener.SyncStarted()

	for {
		update, err := s.client.Recv()
		if err != nil {
			logUpdate("Blockchain sync failed to start. %s", err.Error())
			s.listener.SyncEnded(err)
			break
		}
		if update.Synced {
			s.listener.SyncEnded(nil)
			break
		}

		switch update.NotificationType {
		case walletrpc.SyncNotificationType_FETCHED_HEADERS_STARTED:
			logUpdate("Blockchain sync in progress. Start fetching headers (1/3)")

		case walletrpc.SyncNotificationType_FETCHED_HEADERS_PROGRESS:
			fetchedPercentage := walletmediums.CalculateBlockSyncProgress(s.netType, s.bestBlock, update.FetchHeaders.LastHeaderTime)
			logUpdate("Blockchain sync in progress. Fetching headers (1/3): %d%%", fetchedPercentage)
			s.listener.OnHeadersFetched(fetchedPercentage)

		case walletrpc.SyncNotificationType_FETCHED_HEADERS_FINISHED:
			logUpdate("Blockchain sync in progress. Done fetching headers (1/3)")

		case walletrpc.SyncNotificationType_DISCOVER_ADDRESSES_STARTED:
			logUpdate("Blockchain sync in progress. Start discovering addresses (2/3)")

		case walletrpc.SyncNotificationType_DISCOVER_ADDRESSES_FINISHED:
			logUpdate("Blockchain sync in progress. Finished discovering addresses (2/3)")

		case walletrpc.SyncNotificationType_RESCAN_STARTED:
			logUpdate("Blockchain sync in progress. Start rescanning blocks (3/3)")

		case walletrpc.SyncNotificationType_RESCAN_PROGRESS:
			scannedPercentage := int64(update.RescanProgress.RescannedThrough) / s.bestBlock * 100
			logUpdate("Blockchain sync in progress. Rescanning blocks (3/3): %d%%", scannedPercentage)
			s.listener.OnHeadersFetched(scannedPercentage)

		case walletrpc.SyncNotificationType_RESCAN_FINISHED:
			logUpdate("Blockchain sync in progress. Done rescanning blocks (3/3)")

		case walletrpc.SyncNotificationType_PEER_CONNECTED:
			logUpdate("New peer %s. Connected to %d peers", update.PeerInformation.Address, update.PeerInformation.PeerCount)

		case walletrpc.SyncNotificationType_PEER_DISCONNECTED:
			logUpdate("Peer disconnected %s. Connected to %d peers", update.PeerInformation.Address, update.PeerInformation.PeerCount)
		}
	}
}
