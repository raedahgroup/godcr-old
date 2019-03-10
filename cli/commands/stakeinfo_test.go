package commands

import (
	"context"
	"testing"

	"github.com/raedahgroup/godcr/app/walletcore"
)

func TestStakeInfoCommand_Run(t *testing.T) {
	type fields struct {
		commanderStub commanderStub
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
			g := StakeInfoCommand{
				commanderStub: test.fields.commanderStub,
			}
			if err := g.Run(test.ctx, test.wallet); (err != nil) != test.wantErr {
				t.Errorf("StakeInfoCommand.Run() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
