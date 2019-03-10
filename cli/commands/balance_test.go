package commands

import (
	"context"
	"testing"

	"github.com/raedahgroup/godcr/app/walletcore"
)

func TestBalanceCommand_Run(t *testing.T) {
	type fields struct {
		commanderStub commanderStub
		Detailed      bool
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
			balanceCommand := BalanceCommand{
				commanderStub: test.fields.commanderStub,
				Detailed:      test.fields.Detailed,
			}
			if err := balanceCommand.Run(test.ctx, test.wallet); (err != nil) != test.wantErr {
				t.Errorf("BalanceCommand.Run() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func Test_showDetailedBalance(t *testing.T) {
	tests := []struct {
		name            string
		accountBalances []*walletcore.Account
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			showDetailedBalance(test.accountBalances)
		})
	}
}

func Test_showBalanceSummary(t *testing.T) {
	tests := []struct {
		name     string
		accounts []*walletcore.Account
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			showBalanceSummary(test.accounts)
		})
	}
}
