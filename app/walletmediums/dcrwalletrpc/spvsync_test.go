package dcrwalletrpc

import (
	"testing"

	"github.com/decred/dcrwallet/rpc/walletrpc"
	"github.com/raedahgroup/godcr/app"
)

func Test_spvSync_streamBlockchainSyncUpdates(t *testing.T) {
	type fields struct {
		client    walletrpc.WalletLoaderService_SpvSyncClient
		bestBlock int64
		listener  *app.BlockChainSyncListener
		netType   string
	}
	tests := []struct {
		name    string
		fields  fields
		showLog bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := spvSync{
				client:    tt.fields.client,
				bestBlock: tt.fields.bestBlock,
				listener:  tt.fields.listener,
				netType:   tt.fields.netType,
			}
			s.streamBlockchainSyncUpdates(tt.args.showLog)
		})
	}
}
