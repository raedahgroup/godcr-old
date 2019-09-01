package libwallet

import (
	"strings"

	"github.com/raedahgroup/dcrlibwallet"
)

func (lw *LibWallet) AddSyncProgressListener(syncProgressListener dcrlibwallet.SyncProgressListener,
	uniqueIdentifier string) error {
	return lw.dcrlw.AddSyncProgressListener(syncProgressListener, uniqueIdentifier)
}

func (lw *LibWallet) RemoveSyncProgressListener(uniqueIdentifier string) {
	lw.dcrlw.RemoveSyncProgressListener(uniqueIdentifier)
}

func (lw *LibWallet) SpvSync(showLog bool, persistentPeers []string) error {
	if showLog {
		lw.dcrlw.EnableSyncLogs()
	}

	var peerAddresses string
	if persistentPeers != nil && len(persistentPeers) > 0 {
		peerAddresses = strings.Join(persistentPeers, ";")
	}

	return lw.dcrlw.SpvSync(peerAddresses)
}
