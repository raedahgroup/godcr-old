package commands

import (
	"context"
	"testing"
)

func TestCreateWalletCommand_Run(t *testing.T) {
	type fields struct {
		commanderStub commanderStub
	}
	tests := []struct {
		name    string
		fields  fields
		ctx     context.Context
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := CreateWalletCommand{
				commanderStub: test.fields.commanderStub,
			}
			if err := c.Run(test.ctx); (err != nil) != test.wantErr {
				t.Errorf("CreateWalletCommand.Run() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
