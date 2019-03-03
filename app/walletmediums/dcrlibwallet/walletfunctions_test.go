package dcrlibwallet

import (
	"context"
	"reflect"
	"testing"

	"github.com/decred/dcrwallet/netparams"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/dcrlibwallet/txhelper"
	"github.com/raedahgroup/godcr/app/walletcore"
)

func TestDcrWalletLib_AccountBalance(t *testing.T) {
	type fields struct {
		walletLib *dcrlibwallet.LibWallet
		activeNet *netparams.Params
	}
	tests := []struct {
		name                  string
		fields                fields
		accountNumber         uint32
		requiredConfirmations int32
		want                  *walletcore.Balance
		wantErr               bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lib := &DcrWalletLib{
				walletLib: test.fields.walletLib,
				activeNet: test.fields.activeNet,
			}
			got, err := lib.AccountBalance(test.accountNumber, test.requiredConfirmations)
			if (err != nil) != test.wantErr {
				t.Errorf("DcrWalletLib.AccountBalance() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("DcrWalletLib.AccountBalance() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestDcrWalletLib_AccountsOverview(t *testing.T) {
	type fields struct {
		walletLib *dcrlibwallet.LibWallet
		activeNet *netparams.Params
	}
	tests := []struct {
		name                  string
		fields                fields
		requiredConfirmations int32
		want                  []*walletcore.Account
		wantErr               bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lib := &DcrWalletLib{
				walletLib: test.fields.walletLib,
				activeNet: test.fields.activeNet,
			}
			got, err := lib.AccountsOverview(test.requiredConfirmations)
			if (err != nil) != test.wantErr {
				t.Errorf("DcrWalletLib.AccountsOverview() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("DcrWalletLib.AccountsOverview() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestDcrWalletLib_NextAccount(t *testing.T) {
	type fields struct {
		walletLib *dcrlibwallet.LibWallet
		activeNet *netparams.Params
	}
	tests := []struct {
		name        string
		fields      fields
		accountName string
		passphrase  string
		want        uint32
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lib := &DcrWalletLib{
				walletLib: test.fields.walletLib,
				activeNet: test.fields.activeNet,
			}
			got, err := lib.NextAccount(test.accountName, test.passphrase)
			if (err != nil) != test.wantErr {
				t.Errorf("DcrWalletLib.NextAccount() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("DcrWalletLib.NextAccount() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestDcrWalletLib_AccountNumber(t *testing.T) {
	type fields struct {
		walletLib *dcrlibwallet.LibWallet
		activeNet *netparams.Params
	}
	tests := []struct {
		name        string
		fields      fields
		accountName string
		want        uint32
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lib := &DcrWalletLib{
				walletLib: test.fields.walletLib,
				activeNet: test.fields.activeNet,
			}
			got, err := lib.AccountNumber(test.accountName)
			if (err != nil) != test.wantErr {
				t.Errorf("DcrWalletLib.AccountNumber() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("DcrWalletLib.AccountNumber() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestDcrWalletLib_AccountName(t *testing.T) {
	type fields struct {
		walletLib *dcrlibwallet.LibWallet
		activeNet *netparams.Params
	}
	tests := []struct {
		name          string
		fields        fields
		accountNumber uint32
		want          string
		wantErr       bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lib := &DcrWalletLib{
				walletLib: test.fields.walletLib,
				activeNet: test.fields.activeNet,
			}
			got, err := lib.AccountName(test.accountNumber)
			if (err != nil) != test.wantErr {
				t.Errorf("DcrWalletLib.AccountName() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("DcrWalletLib.AccountName() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestDcrWalletLib_AddressInfo(t *testing.T) {
	type fields struct {
		walletLib *dcrlibwallet.LibWallet
		activeNet *netparams.Params
	}
	tests := []struct {
		name    string
		fields  fields
		address string
		want    *txhelper.AddressInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lib := &DcrWalletLib{
				walletLib: test.fields.walletLib,
				activeNet: test.fields.activeNet,
			}
			got, err := lib.AddressInfo(test.address)
			if (err != nil) != test.wantErr {
				t.Errorf("DcrWalletLib.AddressInfo() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("DcrWalletLib.AddressInfo() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestDcrWalletLib_ValidateAddress(t *testing.T) {
	type fields struct {
		walletLib *dcrlibwallet.LibWallet
		activeNet *netparams.Params
	}
	tests := []struct {
		name    string
		fields  fields
		address string
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lib := &DcrWalletLib{
				walletLib: test.fields.walletLib,
				activeNet: test.fields.activeNet,
			}
			got, err := lib.ValidateAddress(test.address)
			if (err != nil) != test.wantErr {
				t.Errorf("DcrWalletLib.ValidateAddress() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("DcrWalletLib.ValidateAddress() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestDcrWalletLib_ReceiveAddress(t *testing.T) {
	type fields struct {
		walletLib *dcrlibwallet.LibWallet
		activeNet *netparams.Params
	}
	tests := []struct {
		name    string
		fields  fields
		account uint32
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lib := &DcrWalletLib{
				walletLib: test.fields.walletLib,
				activeNet: test.fields.activeNet,
			}
			got, err := lib.ReceiveAddress(test.account)
			if (err != nil) != test.wantErr {
				t.Errorf("DcrWalletLib.ReceiveAddress() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("DcrWalletLib.ReceiveAddress() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestDcrWalletLib_GenerateNewAddress(t *testing.T) {
	type fields struct {
		walletLib *dcrlibwallet.LibWallet
		activeNet *netparams.Params
	}
	tests := []struct {
		name    string
		fields  fields
		account uint32
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lib := &DcrWalletLib{
				walletLib: test.fields.walletLib,
				activeNet: test.fields.activeNet,
			}
			got, err := lib.GenerateNewAddress(test.account)
			if (err != nil) != test.wantErr {
				t.Errorf("DcrWalletLib.GenerateNewAddress() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("DcrWalletLib.GenerateNewAddress() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestDcrWalletLib_UnspentOutputs(t *testing.T) {
	type fields struct {
		walletLib *dcrlibwallet.LibWallet
		activeNet *netparams.Params
	}
	tests := []struct {
		name                  string
		fields                fields
		account               uint32
		targetAmount          int64
		requiredConfirmations int32
		want                  []*walletcore.UnspentOutput
		wantErr               bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lib := &DcrWalletLib{
				walletLib: test.fields.walletLib,
				activeNet: test.fields.activeNet,
			}
			got, err := lib.UnspentOutputs(test.account, test.targetAmount, test.requiredConfirmations)
			if (err != nil) != test.wantErr {
				t.Errorf("DcrWalletLib.UnspentOutputs() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("DcrWalletLib.UnspentOutputs() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestDcrWalletLib_SendFromAccount(t *testing.T) {
	type fields struct {
		walletLib *dcrlibwallet.LibWallet
		activeNet *netparams.Params
	}
	tests := []struct {
		name                  string
		fields                fields
		sourceAccount         uint32
		requiredConfirmations int32
		destinations          []txhelper.TransactionDestination
		passphrase            string
		want                  string
		wantErr               bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lib := &DcrWalletLib{
				walletLib: test.fields.walletLib,
				activeNet: test.fields.activeNet,
			}
			got, err := lib.SendFromAccount(test.sourceAccount, test.requiredConfirmations, test.destinations, test.passphrase)
			if (err != nil) != test.wantErr {
				t.Errorf("DcrWalletLib.SendFromAccount() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("DcrWalletLib.SendFromAccount() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestDcrWalletLib_SendFromUTXOs(t *testing.T) {
	type fields struct {
		walletLib *dcrlibwallet.LibWallet
		activeNet *netparams.Params
	}
	tests := []struct {
		name                  string
		fields                fields
		sourceAccount         uint32
		requiredConfirmations int32
		utxoKeys              []string
		txDestinations        []txhelper.TransactionDestination
		changeDestinations    []txhelper.TransactionDestination
		passphrase            string
		want                  string
		wantErr               bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lib := &DcrWalletLib{
				walletLib: test.fields.walletLib,
				activeNet: test.fields.activeNet,
			}
			got, err := lib.SendFromUTXOs(test.sourceAccount, test.requiredConfirmations, test.utxoKeys, test.txDestinations, test.changeDestinations, test.passphrase)
			if (err != nil) != test.wantErr {
				t.Errorf("DcrWalletLib.SendFromUTXOs() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("DcrWalletLib.SendFromUTXOs() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestDcrWalletLib_TransactionHistory(t *testing.T) {
	type fields struct {
		walletLib *dcrlibwallet.LibWallet
		activeNet *netparams.Params
	}
	tests := []struct {
		name    string
		fields  fields
		want    []*walletcore.Transaction
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lib := &DcrWalletLib{
				walletLib: test.fields.walletLib,
				activeNet: test.fields.activeNet,
			}
			got, err := lib.TransactionHistory()
			if (err != nil) != test.wantErr {
				t.Errorf("DcrWalletLib.TransactionHistory() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("DcrWalletLib.TransactionHistory() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestDcrWalletLib_GetTransaction(t *testing.T) {
	type fields struct {
		walletLib *dcrlibwallet.LibWallet
		activeNet *netparams.Params
	}
	tests := []struct {
		name            string
		fields          fields
		transactionHash string
		want            *walletcore.TransactionDetails
		wantErr         bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lib := &DcrWalletLib{
				walletLib: test.fields.walletLib,
				activeNet: test.fields.activeNet,
			}
			got, err := lib.GetTransaction(test.transactionHash)
			if (err != nil) != test.wantErr {
				t.Errorf("DcrWalletLib.GetTransaction() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("DcrWalletLib.GetTransaction() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestDcrWalletLib_StakeInfo(t *testing.T) {
	type fields struct {
		walletLib *dcrlibwallet.LibWallet
		activeNet *netparams.Params
	}
	tests := []struct {
		name    string
		fields  fields
		ctx     context.Context
		want    *walletcore.StakeInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lib := &DcrWalletLib{
				walletLib: test.fields.walletLib,
				activeNet: test.fields.activeNet,
			}
			got, err := lib.StakeInfo(test.ctx)
			if (err != nil) != test.wantErr {
				t.Errorf("DcrWalletLib.StakeInfo() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("DcrWalletLib.StakeInfo() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestDcrWalletLib_PurchaseTickets(t *testing.T) {
	type fields struct {
		walletLib *dcrlibwallet.LibWallet
		activeNet *netparams.Params
	}
	tests := []struct {
		name    string
		fields  fields
		ctx     context.Context
		request dcrlibwallet.PurchaseTicketsRequest
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lib := &DcrWalletLib{
				walletLib: test.fields.walletLib,
				activeNet: test.fields.activeNet,
			}
			got, err := lib.PurchaseTickets(test.ctx, test.request)
			if (err != nil) != test.wantErr {
				t.Errorf("DcrWalletLib.PurchaseTickets() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("DcrWalletLib.PurchaseTickets() = %v, want %v", got, test.want)
			}
		})
	}
}
