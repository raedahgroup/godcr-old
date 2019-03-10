package commands

import (
	"context"
	"testing"

	"github.com/raedahgroup/godcr/app/walletcore"
)

func TestReceiveCommand_Run(t *testing.T) {
	type fields struct {
		commanderStub commanderStub
		Args          ReceiveCommandArgs
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
			receiveCommand := ReceiveCommand{
				commanderStub: test.fields.commanderStub,
				Args:          test.fields.Args,
			}
			if err := receiveCommand.Run(test.ctx, test.wallet); (err != nil) != test.wantErr {
				t.Errorf("ReceiveCommand.Run() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
