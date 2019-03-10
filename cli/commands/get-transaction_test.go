package commands

import (
	"context"
	"testing"

	"github.com/raedahgroup/godcr/app/walletcore"
)

func TestShowTransactionCommand_Run(t *testing.T) {
	type fields struct {
		commanderStub commanderStub
		Detailed      bool
		Args          ShowTransactionCommandArgs
	}
	tests := []struct {
		name    string
		fields  fields
		ctx     context.Context
		wallet  walletcore.Wallet
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			showTxCommand := ShowTransactionCommand{
				commanderStub: test.fields.commanderStub,
				Detailed:      test.fields.Detailed,
				Args:          test.fields.Args,
			}
			if err := showTxCommand.Run(test.ctx, test.wallet); (err != nil) != test.wantErr {
				t.Errorf("ShowTransactionCommand.Run() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
