package commands

import (
	"context"
	"testing"

	"github.com/raedahgroup/godcr/app/walletcore"
)

func TestPurchaseTicketsCommand_Run(t *testing.T) {
	type fields struct {
		commanderStub    commanderStub
		MinConfirmations uint32
		TicketAddress    string
		NumTickets       uint32
		PoolAddress      string
		PoolFees         float64
		Expiry           uint32
		TxFee            int64
		TicketFee        int64
		PayFrom          string
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
			ptc := PurchaseTicketsCommand{
				commanderStub:    test.fields.commanderStub,
				MinConfirmations: test.fields.MinConfirmations,
				TicketAddress:    test.fields.TicketAddress,
				NumTickets:       test.fields.NumTickets,
				PoolAddress:      test.fields.PoolAddress,
				PoolFees:         test.fields.PoolFees,
				Expiry:           test.fields.Expiry,
				TxFee:            test.fields.TxFee,
				TicketFee:        test.fields.TicketFee,
				PayFrom:          test.fields.PayFrom,
			}
			if err := ptc.Run(test.ctx, test.wallet); (err != nil) != test.wantErr {
				t.Errorf("PurchaseTicketsCommand.Run() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
