package commands

import (
	"reflect"
	"testing"

	"github.com/raedahgroup/dcrlibwallet/txhelper"
	"github.com/raedahgroup/godcr/app/walletcore"
)

func Test_selectAccount(t *testing.T) {
	wallet, err := getWalletForTesting()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		wallet  walletcore.Wallet
		want    uint32
		wantErr bool
	}{
		{
			name:    "select account",
			wallet:  wallet,
			want:    0,
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := selectAccount(test.wallet)
			if (err != nil) != test.wantErr {
				t.Errorf("selectAccount() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("selectAccount() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_getSendAmount(t *testing.T) {
	tests := []struct {
		name    string
		want    float64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := getSendAmount()
			if (err != nil) != test.wantErr {
				t.Errorf("getSendAmount() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("getSendAmount() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_getChangeOutputDestinations(t *testing.T) {
	tests := []struct {
		name             string
		wallet           walletcore.Wallet
		totalInputAmount float64
		sourceAccount    uint32
		nUtxoSelection   int
		sendDestinations []txhelper.TransactionDestination
		want             []txhelper.TransactionDestination
		wantErr          bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := getChangeOutputDestinations(test.wallet, test.totalInputAmount, test.sourceAccount, test.nUtxoSelection, test.sendDestinations)
			if (err != nil) != test.wantErr {
				t.Errorf("getChangeOutputDestinations() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("getChangeOutputDestinations() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_getChangeDestinationsFromUser(t *testing.T) {
	tests := []struct {
		name             string
		wallet           walletcore.Wallet
		amountInAtom     int64
		sourceAccount    uint32
		nUtxoSelection   int
		sendDestinations []txhelper.TransactionDestination
		want             []txhelper.TransactionDestination
		wantErr          bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := getChangeDestinationsFromUser(test.wallet, test.amountInAtom, test.sourceAccount, test.nUtxoSelection, test.sendDestinations)
			if (err != nil) != test.wantErr {
				t.Errorf("getChangeDestinationsFromUser() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("getChangeDestinationsFromUser() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_getChangeAmount(t *testing.T) {
	tests := []struct {
		name       string
		prompt     string
		wantAmount float64
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotAmount, err := getChangeAmount(test.prompt)
			if (err != nil) != test.wantErr {
				t.Errorf("getChangeAmount() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if gotAmount != test.wantAmount {
				t.Errorf("getChangeAmount() = %v, want %v", gotAmount, test.wantAmount)
			}
		})
	}
}

func Test_getWalletPassphrase(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := getWalletPassphrase()
			if (err != nil) != test.wantErr {
				t.Errorf("getWalletPassphrase() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("getWalletPassphrase() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_getUtxosForNewTransaction(t *testing.T) {
	tests := []struct {
		name                    string
		utxos                   []*walletcore.UnspentOutput
		sendAmount              float64
		wantSelectedUtxos       []*walletcore.UnspentOutput
		wantTotalAmountSelected float64
		wantErr                 bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotSelectedUtxos, gotTotalAmountSelected, err := getUtxosForNewTransaction(test.utxos, test.sendAmount)
			if (err != nil) != test.wantErr {
				t.Errorf("getUtxosForNewTransaction() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSelectedUtxos, test.wantSelectedUtxos) {
				t.Errorf("getUtxosForNewTransaction() gotSelectedUtxos = %v, want %v", gotSelectedUtxos, test.wantSelectedUtxos)
			}
			if gotTotalAmountSelected != test.wantTotalAmountSelected {
				t.Errorf("getUtxosForNewTransaction() gotTotalAmountSelected = %v, want %v", gotTotalAmountSelected, test.wantTotalAmountSelected)
			}
		})
	}
}

func Test_bestSizedInput(t *testing.T) {
	tests := []struct {
		name            string
		utxos           []*walletcore.UnspentOutput
		sendAmountTotal float64
		want            []*walletcore.UnspentOutput
		want1           float64
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, got1 := bestSizedInput(test.utxos, test.sendAmountTotal)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("bestSizedInput() got = %v, want %v", got, test.want)
			}
			if got1 != test.want1 {
				t.Errorf("bestSizedInput() got1 = %v, want %v", got1, test.want1)
			}
		})
	}
}
