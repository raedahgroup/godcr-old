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
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := spvSync{
				client:    test.fields.client,
				bestBlock: test.fields.bestBlock,
				listener:  test.fields.listener,
				netType:   test.fields.netType,
			}
			s.streamBlockchainSyncUpdates(test.showLog)
		})
	}
}
