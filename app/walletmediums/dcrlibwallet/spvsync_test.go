package dcrlibwallet

import (
	"testing"

	"github.com/decred/dcrwallet/netparams"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/app"
)

func TestSpvSyncResponse_OnPeerConnected(t *testing.T) {
	type fields struct {
		activeNet *netparams.Params
		walletLib *dcrlibwallet.LibWallet
		listener  *app.BlockChainSyncListener
	}
	tests := []struct {
		name      string
		fields    fields
		peerCount int32
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := SpvSyncResponse{
				activeNet: test.fields.activeNet,
				walletLib: test.fields.walletLib,
				listener:  test.fields.listener,
			}
			response.OnPeerConnected(test.peerCount)
		})
	}
}

func TestSpvSyncResponse_OnPeerDisconnected(t *testing.T) {
	type fields struct {
		activeNet *netparams.Params
		walletLib *dcrlibwallet.LibWallet
		listener  *app.BlockChainSyncListener
	}
	tests := []struct {
		name      string
		fields    fields
		peerCount int32
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := SpvSyncResponse{
				activeNet: test.fields.activeNet,
				walletLib: test.fields.walletLib,
				listener:  test.fields.listener,
			}
			response.OnPeerDisconnected(test.peerCount)
		})
	}
}

func TestSpvSyncResponse_OnFetchMissingCFilters(t *testing.T) {
	type fields struct {
		activeNet *netparams.Params
		walletLib *dcrlibwallet.LibWallet
		listener  *app.BlockChainSyncListener
	}
	tests := []struct {
		name                 string
		fields               fields
		missingCFitlersStart int32
		missingCFitlersEnd   int32
		state                string
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := SpvSyncResponse{
				activeNet: test.fields.activeNet,
				walletLib: test.fields.walletLib,
				listener:  test.fields.listener,
			}
			response.OnFetchMissingCFilters(test.missingCFitlersStart, test.missingCFitlersEnd, test.state)
		})
	}
}

func TestSpvSyncResponse_OnFetchedHeaders(t *testing.T) {
	type fields struct {
		activeNet *netparams.Params
		walletLib *dcrlibwallet.LibWallet
		listener  *app.BlockChainSyncListener
	}
	tests := []struct {
		name           string
		fields         fields
		in0            int32
		lastHeaderTime int64
		state          string
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := SpvSyncResponse{
				activeNet: test.fields.activeNet,
				walletLib: test.fields.walletLib,
				listener:  test.fields.listener,
			}
			response.OnFetchedHeaders(test.in0, test.lastHeaderTime, test.state)
		})
	}
}

func TestSpvSyncResponse_OnDiscoveredAddresses(t *testing.T) {
	type fields struct {
		activeNet *netparams.Params
		walletLib *dcrlibwallet.LibWallet
		listener  *app.BlockChainSyncListener
	}
	tests := []struct {
		name   string
		fields fields
		state  string
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := SpvSyncResponse{
				activeNet: test.fields.activeNet,
				walletLib: test.fields.walletLib,
				listener:  test.fields.listener,
			}
			response.OnDiscoveredAddresses(test.state)
		})
	}
}

func TestSpvSyncResponse_OnRescan(t *testing.T) {
	type fields struct {
		activeNet *netparams.Params
		walletLib *dcrlibwallet.LibWallet
		listener  *app.BlockChainSyncListener
	}
	tests := []struct {
		name             string
		fields           fields
		rescannedThrough int32
		state            string
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := SpvSyncResponse{
				activeNet: test.fields.activeNet,
				walletLib: test.fields.walletLib,
				listener:  test.fields.listener,
			}
			response.OnRescan(test.rescannedThrough, test.state)
		})
	}
}

func TestSpvSyncResponse_OnSynced(t *testing.T) {
	type fields struct {
		activeNet *netparams.Params
		walletLib *dcrlibwallet.LibWallet
		listener  *app.BlockChainSyncListener
	}
	tests := []struct {
		name   string
		fields fields
		synced bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := SpvSyncResponse{
				activeNet: test.fields.activeNet,
				walletLib: test.fields.walletLib,
				listener:  test.fields.listener,
			}
			response.OnSynced(test.synced)
		})
	}
}

func TestSpvSyncResponse_OnSyncError(t *testing.T) {
	type fields struct {
		activeNet *netparams.Params
		walletLib *dcrlibwallet.LibWallet
		listener  *app.BlockChainSyncListener
	}
	tests := []struct {
		name   string
		fields fields
		code   int
		err    error
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := SpvSyncResponse{
				activeNet: test.fields.activeNet,
				walletLib: test.fields.walletLib,
				listener:  test.fields.listener,
			}
			response.OnSyncError(test.code, test.err)
		})
	}
}

func TestSpvSyncResponse_calculateProgress(t *testing.T) {
	type fields struct {
		activeNet *netparams.Params
		walletLib *dcrlibwallet.LibWallet
		listener  *app.BlockChainSyncListener
	}
	tests := []struct {
		name           string
		fields         fields
		lastHeaderTime int64
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := SpvSyncResponse{
				activeNet: test.fields.activeNet,
				walletLib: test.fields.walletLib,
				listener:  test.fields.listener,
			}
			response.calculateProgress(test.lastHeaderTime)
		})
	}
}
