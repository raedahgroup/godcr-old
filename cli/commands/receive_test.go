package commands

import (
	"context"
	"testing"

	"github.com/raedahgroup/godcr/app/walletcore"
)

func TestReceiveCommand_Run(t *testing.T) {
	wallet, err := getWalletForTesting()
	if err != nil {
		t.Fatal(err)
	}

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
		{
			name: "receive command",
			fields: fields{
				commanderStub: commanderStub{},
				Args: ReceiveCommandArgs{
					AccountName: "default",
				},
			},
			ctx:     context.Background(),
			wallet:  wallet,
			wantErr: false,
		},
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
