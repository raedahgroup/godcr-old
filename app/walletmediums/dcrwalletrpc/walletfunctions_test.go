package dcrwalletrpc

import (
	"context"
	"reflect"
	"testing"

	"github.com/decred/dcrwallet/netparams"
	"github.com/decred/dcrwallet/rpc/walletrpc"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/dcrlibwallet/txhelper"
	"github.com/raedahgroup/godcr/app/walletcore"
)

func TestWalletRPCClient_AccountBalance(t *testing.T) {
	type fields struct {
		walletLoader  walletrpc.WalletLoaderServiceClient
		walletService walletrpc.WalletServiceClient
		activeNet     *netparams.Params
		walletOpen    bool
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
			c := &WalletRPCClient{
				walletLoader:  test.fields.walletLoader,
				walletService: test.fields.walletService,
				activeNet:     test.fields.activeNet,
				walletOpen:    test.fields.walletOpen,
			}
			got, err := c.AccountBalance(test.accountNumber, test.requiredConfirmations)
			if (err != nil) != test.wantErr {
				t.Errorf("WalletRPCClient.AccountBalance() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("WalletRPCClient.AccountBalance() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestWalletRPCClient_AccountsOverview(t *testing.T) {
	type fields struct {
		walletLoader  walletrpc.WalletLoaderServiceClient
		walletService walletrpc.WalletServiceClient
		activeNet     *netparams.Params
		walletOpen    bool
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
			c := &WalletRPCClient{
				walletLoader:  test.fields.walletLoader,
				walletService: test.fields.walletService,
				activeNet:     test.fields.activeNet,
				walletOpen:    test.fields.walletOpen,
			}
			got, err := c.AccountsOverview(test.requiredConfirmations)
			if (err != nil) != test.wantErr {
				t.Errorf("WalletRPCClient.AccountsOverview() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("WalletRPCClient.AccountsOverview() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestWalletRPCClient_NextAccount(t *testing.T) {
	type fields struct {
		walletLoader  walletrpc.WalletLoaderServiceClient
		walletService walletrpc.WalletServiceClient
		activeNet     *netparams.Params
		walletOpen    bool
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
			c := &WalletRPCClient{
				walletLoader:  test.fields.walletLoader,
				walletService: test.fields.walletService,
				activeNet:     test.fields.activeNet,
				walletOpen:    test.fields.walletOpen,
			}
			got, err := c.NextAccount(test.accountName, test.passphrase)
			if (err != nil) != test.wantErr {
				t.Errorf("WalletRPCClient.NextAccount() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("WalletRPCClient.NextAccount() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestWalletRPCClient_AccountNumber(t *testing.T) {
	type fields struct {
		walletLoader  walletrpc.WalletLoaderServiceClient
		walletService walletrpc.WalletServiceClient
		activeNet     *netparams.Params
		walletOpen    bool
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
			c := &WalletRPCClient{
				walletLoader:  test.fields.walletLoader,
				walletService: test.fields.walletService,
				activeNet:     test.fields.activeNet,
				walletOpen:    test.fields.walletOpen,
			}
			got, err := c.AccountNumber(test.accountName)
			if (err != nil) != test.wantErr {
				t.Errorf("WalletRPCClient.AccountNumber() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("WalletRPCClient.AccountNumber() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestWalletRPCClient_AccountName(t *testing.T) {
	type fields struct {
		walletLoader  walletrpc.WalletLoaderServiceClient
		walletService walletrpc.WalletServiceClient
		activeNet     *netparams.Params
		walletOpen    bool
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
			c := &WalletRPCClient{
				walletLoader:  test.fields.walletLoader,
				walletService: test.fields.walletService,
				activeNet:     test.fields.activeNet,
				walletOpen:    test.fields.walletOpen,
			}
			got, err := c.AccountName(test.accountNumber)
			if (err != nil) != test.wantErr {
				t.Errorf("WalletRPCClient.AccountName() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("WalletRPCClient.AccountName() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestWalletRPCClient_AddressInfo(t *testing.T) {
	type fields struct {
		walletLoader  walletrpc.WalletLoaderServiceClient
		walletService walletrpc.WalletServiceClient
		activeNet     *netparams.Params
		walletOpen    bool
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
			c := &WalletRPCClient{
				walletLoader:  test.fields.walletLoader,
				walletService: test.fields.walletService,
				activeNet:     test.fields.activeNet,
				walletOpen:    test.fields.walletOpen,
			}
			got, err := c.AddressInfo(test.address)
			if (err != nil) != test.wantErr {
				t.Errorf("WalletRPCClient.AddressInfo() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("WalletRPCClient.AddressInfo() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestWalletRPCClient_ValidateAddress(t *testing.T) {
	type fields struct {
		walletLoader  walletrpc.WalletLoaderServiceClient
		walletService walletrpc.WalletServiceClient
		activeNet     *netparams.Params
		walletOpen    bool
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
			c := &WalletRPCClient{
				walletLoader:  test.fields.walletLoader,
				walletService: test.fields.walletService,
				activeNet:     test.fields.activeNet,
				walletOpen:    test.fields.walletOpen,
			}
			got, err := c.ValidateAddress(test.address)
			if (err != nil) != test.wantErr {
				t.Errorf("WalletRPCClient.ValidateAddress() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("WalletRPCClient.ValidateAddress() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestWalletRPCClient_ReceiveAddress(t *testing.T) {
	type fields struct {
		walletLoader  walletrpc.WalletLoaderServiceClient
		walletService walletrpc.WalletServiceClient
		activeNet     *netparams.Params
		walletOpen    bool
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
			c := &WalletRPCClient{
				walletLoader:  test.fields.walletLoader,
				walletService: test.fields.walletService,
				activeNet:     test.fields.activeNet,
				walletOpen:    test.fields.walletOpen,
			}
			got, err := c.ReceiveAddress(test.account)
			if (err != nil) != test.wantErr {
				t.Errorf("WalletRPCClient.ReceiveAddress() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("WalletRPCClient.ReceiveAddress() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestWalletRPCClient_GenerateNewAddress(t *testing.T) {
	type fields struct {
		walletLoader  walletrpc.WalletLoaderServiceClient
		walletService walletrpc.WalletServiceClient
		activeNet     *netparams.Params
		walletOpen    bool
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
			c := &WalletRPCClient{
				walletLoader:  test.fields.walletLoader,
				walletService: test.fields.walletService,
				activeNet:     test.fields.activeNet,
				walletOpen:    test.fields.walletOpen,
			}
			got, err := c.GenerateNewAddress(test.account)
			if (err != nil) != test.wantErr {
				t.Errorf("WalletRPCClient.GenerateNewAddress() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("WalletRPCClient.GenerateNewAddress() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestWalletRPCClient_UnspentOutputs(t *testing.T) {
	type fields struct {
		walletLoader  walletrpc.WalletLoaderServiceClient
		walletService walletrpc.WalletServiceClient
		activeNet     *netparams.Params
		walletOpen    bool
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
			c := &WalletRPCClient{
				walletLoader:  test.fields.walletLoader,
				walletService: test.fields.walletService,
				activeNet:     test.fields.activeNet,
				walletOpen:    test.fields.walletOpen,
			}
			got, err := c.UnspentOutputs(test.account, test.targetAmount, test.requiredConfirmations)
			if (err != nil) != test.wantErr {
				t.Errorf("WalletRPCClient.UnspentOutputs() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("WalletRPCClient.UnspentOutputs() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestWalletRPCClient_SendFromAccount(t *testing.T) {
	type fields struct {
		walletLoader  walletrpc.WalletLoaderServiceClient
		walletService walletrpc.WalletServiceClient
		activeNet     *netparams.Params
		walletOpen    bool
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
			c := &WalletRPCClient{
				walletLoader:  test.fields.walletLoader,
				walletService: test.fields.walletService,
				activeNet:     test.fields.activeNet,
				walletOpen:    test.fields.walletOpen,
			}
			got, err := c.SendFromAccount(test.sourceAccount, test.requiredConfirmations, test.destinations, test.passphrase)
			if (err != nil) != test.wantErr {
				t.Errorf("WalletRPCClient.SendFromAccount() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("WalletRPCClient.SendFromAccount() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestWalletRPCClient_SendFromUTXOs(t *testing.T) {
	type fields struct {
		walletLoader  walletrpc.WalletLoaderServiceClient
		walletService walletrpc.WalletServiceClient
		activeNet     *netparams.Params
		walletOpen    bool
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
			c := &WalletRPCClient{
				walletLoader:  test.fields.walletLoader,
				walletService: test.fields.walletService,
				activeNet:     test.fields.activeNet,
				walletOpen:    test.fields.walletOpen,
			}
			got, err := c.SendFromUTXOs(test.sourceAccount, test.requiredConfirmations, test.utxoKeys, test.txDestinations, test.changeDestinations, test.passphrase)
			if (err != nil) != test.wantErr {
				t.Errorf("WalletRPCClient.SendFromUTXOs() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("WalletRPCClient.SendFromUTXOs() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestWalletRPCClient_TransactionHistory(t *testing.T) {
	type fields struct {
		walletLoader  walletrpc.WalletLoaderServiceClient
		walletService walletrpc.WalletServiceClient
		activeNet     *netparams.Params
		walletOpen    bool
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
			c := &WalletRPCClient{
				walletLoader:  test.fields.walletLoader,
				walletService: test.fields.walletService,
				activeNet:     test.fields.activeNet,
				walletOpen:    test.fields.walletOpen,
			}
			got, err := c.TransactionHistory()
			if (err != nil) != test.wantErr {
				t.Errorf("WalletRPCClient.TransactionHistory() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("WalletRPCClient.TransactionHistory() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestWalletRPCClient_GetTransaction(t *testing.T) {
	type fields struct {
		walletLoader  walletrpc.WalletLoaderServiceClient
		walletService walletrpc.WalletServiceClient
		activeNet     *netparams.Params
		walletOpen    bool
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
			c := &WalletRPCClient{
				walletLoader:  test.fields.walletLoader,
				walletService: test.fields.walletService,
				activeNet:     test.fields.activeNet,
				walletOpen:    test.fields.walletOpen,
			}
			got, err := c.GetTransaction(test.transactionHash)
			if (err != nil) != test.wantErr {
				t.Errorf("WalletRPCClient.GetTransaction() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("WalletRPCClient.GetTransaction() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestWalletRPCClient_StakeInfo(t *testing.T) {
	type fields struct {
		walletLoader  walletrpc.WalletLoaderServiceClient
		walletService walletrpc.WalletServiceClient
		activeNet     *netparams.Params
		walletOpen    bool
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
			c := &WalletRPCClient{
				walletLoader:  test.fields.walletLoader,
				walletService: test.fields.walletService,
				activeNet:     test.fields.activeNet,
				walletOpen:    test.fields.walletOpen,
			}
			got, err := c.StakeInfo(test.ctx)
			if (err != nil) != test.wantErr {
				t.Errorf("WalletRPCClient.StakeInfo() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("WalletRPCClient.StakeInfo() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestWalletRPCClient_PurchaseTickets(t *testing.T) {
	type fields struct {
		walletLoader  walletrpc.WalletLoaderServiceClient
		walletService walletrpc.WalletServiceClient
		activeNet     *netparams.Params
		walletOpen    bool
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
			c := &WalletRPCClient{
				walletLoader:  test.fields.walletLoader,
				walletService: test.fields.walletService,
				activeNet:     test.fields.activeNet,
				walletOpen:    test.fields.walletOpen,
			}
			got, err := c.PurchaseTicket(test.ctx, test.request)
			if (err != nil) != test.wantErr {
				t.Errorf("WalletRPCClient.PurchaseTickets() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("WalletRPCClient.PurchaseTickets() = %v, want %v", got, test.want)
			}
		})
	}
}
