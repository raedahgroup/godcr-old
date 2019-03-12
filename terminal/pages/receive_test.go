package pages

import (
	"reflect"
	"testing"

	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/rivo/tview"
	qrcode "github.com/skip2/go-qrcode"
)

func TestReceivePage(t *testing.T) {
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
			if got := ReceivePage(test.wallet, test.setFocus, test.clearFocus); !reflect.DeepEqual(got, test.want) {
				t.Errorf("ReceivePage() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_generateAddress(t *testing.T) {
	tests := []struct {
		name          string
		wallet        walletcore.Wallet
		accountNumber uint32
		want          string
		want1         *qrcode.QRCode
		wantErr       bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, got1, err := generateAddress(test.wallet, test.accountNumber)
			if (err != nil) != test.wantErr {
				t.Errorf("generateAddress() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("generateAddress() got = %v, want %v", got, test.want)
			}
			if !reflect.DeepEqual(got1, test.want1) {
				t.Errorf("generateAddress() got1 = %v, want %v", got1, test.want1)
			}
		})
	}
}
