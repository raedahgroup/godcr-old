package commands

import (
	"context"
	"testing"

	"github.com/raedahgroup/dcrlibwallet/txhelper"
	"github.com/raedahgroup/godcr/app/walletcore"
)

func TestSendCommand_Run(t *testing.T) {
	type fields struct {
		commanderStub    commanderStub
		SpendUnconfirmed bool
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
			s := SendCommand{
				commanderStub:    test.fields.commanderStub,
				SpendUnconfirmed: test.fields.SpendUnconfirmed,
			}
			if err := s.Run(test.ctx, test.wallet); (err != nil) != test.wantErr {
				t.Errorf("SendCommand.Run() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func TestSendCustomCommand_Run(t *testing.T) {
	type fields struct {
		commanderStub    commanderStub
		SpendUnconfirmed bool
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
			s := SendCustomCommand{
				commanderStub:    test.fields.commanderStub,
				SpendUnconfirmed: test.fields.SpendUnconfirmed,
			}
			if err := s.Run(test.ctx, test.wallet); (err != nil) != test.wantErr {
				t.Errorf("SendCustomCommand.Run() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func Test_send(t *testing.T) {
	tests := []struct {
		name             string
		wallet           walletcore.Wallet
		spendUnconfirmed bool
		custom           bool
		wantErr          bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := send(test.wallet, test.spendUnconfirmed, test.custom); (err != nil) != test.wantErr {
				t.Errorf("send() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func Test_completeCustomSend(t *testing.T) {
	tests := []struct {
		name                  string
		wallet                walletcore.Wallet
		sourceAccount         uint32
		sendDestinations      []txhelper.TransactionDestination
		sendAmountTotal       float64
		requiredConfirmations int32
		want                  string
		wantErr               bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := completeCustomSend(test.wallet, test.sourceAccount, test.sendDestinations, test.sendAmountTotal, test.requiredConfirmations)
			if (err != nil) != test.wantErr {
				t.Errorf("completeCustomSend() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("completeCustomSend() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_completeNormalSend(t *testing.T) {
	tests := []struct {
		name                  string
		wallet                walletcore.Wallet
		sourceAccount         uint32
		sendDestinations      []txhelper.TransactionDestination
		requiredConfirmations int32
		want                  string
		wantErr               bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := completeNormalSend(test.wallet, test.sourceAccount, test.sendDestinations, test.requiredConfirmations)
			if (err != nil) != test.wantErr {
				t.Errorf("completeNormalSend() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("completeNormalSend() = %v, want %v", got, test.want)
			}
		})
	}
}
