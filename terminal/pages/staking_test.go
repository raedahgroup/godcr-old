package pages

import (
	"reflect"
	"testing"

	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

func TestStakingPage(t *testing.T) {
	tests := []struct {
		name       string
		wallet     walletcore.Wallet
		setFocus   func(p tview.Primitive) *tview.Application
		clearFocus func()
		want       tview.Primitive
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := StakingPage(test.wallet, test.setFocus, test.clearFocus); !reflect.DeepEqual(got, test.want) {
				t.Errorf("StakingPage() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_stakeInfoTable(t *testing.T) {
	tests := []struct {
		name    string
		wallet  walletcore.Wallet
		want    *tview.Table
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := stakeInfoTable(test.wallet)
			if (err != nil) != test.wantErr {
				t.Errorf("stakeInfoTable() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("stakeInfoTable() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_purchaseTicketForm(t *testing.T) {
	tests := []struct {
		name    string
		wallet  walletcore.Wallet
		want    *primitives.Form
		want1   *tview.TextView
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, got1, err := purchaseTicketForm(test.wallet)
			if (err != nil) != test.wantErr {
				t.Errorf("purchaseTicketForm() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("purchaseTicketForm() got = %v, want %v", got, test.want)
			}
			if !reflect.DeepEqual(got1, test.want1) {
				t.Errorf("purchaseTicketForm() got1 = %v, want %v", got1, test.want1)
			}
		})
	}
}

func Test_purchaseTickets(t *testing.T) {
	tests := []struct {
		name             string
		passphrase       string
		numTickets       string
		accountNum       uint32
		spendUnconfirmed bool
		wallet           walletcore.Wallet
		want             []string
		wantErr          bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := purchaseTickets(test.passphrase, test.numTickets, test.accountNum, test.spendUnconfirmed, test.wallet)
			if (err != nil) != test.wantErr {
				t.Errorf("purchaseTickets() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("purchaseTickets() = %v, want %v", got, test.want)
			}
		})
	}
}
