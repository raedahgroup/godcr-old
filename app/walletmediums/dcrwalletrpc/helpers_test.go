package dcrwalletrpc

import (
	"reflect"
	"testing"

	"github.com/decred/dcrd/chaincfg"
	"github.com/decred/dcrwallet/rpc/walletrpc"
	"github.com/raedahgroup/dcrlibwallet/txhelper"
	"github.com/raedahgroup/godcr/app/walletcore"
)

func TestWalletRPCClient_unspentOutputStream(t *testing.T) {
	type fields struct {
		walletLoader  walletrpc.WalletLoaderServiceClient
		walletService walletrpc.WalletServiceClient
		activeNet     *chaincfg.Params
		walletOpen    bool
	}
	tests := []struct {
		name                  string
		fields                fields
		account               uint32
		targetAmount          int64
		requiredConfirmations int32
		want                  walletrpc.WalletService_UnspentOutputsClient
		wantErr               bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &WalletRPCClient{
				walletLoader:  tt.fields.walletLoader,
				walletService: tt.fields.walletService,
				activeNet:     tt.fields.activeNet,
				walletOpen:    tt.fields.walletOpen,
			}
			got, err := c.unspentOutputStream(tt.account, tt.targetAmount, tt.requiredConfirmations)
			if (err != nil) != tt.wantErr {
				t.Errorf("WalletRPCClient.unspentOutputStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WalletRPCClient.unspentOutputStream() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWalletRPCClient_signAndPublishTransaction(t *testing.T) {
	type fields struct {
		walletLoader  walletrpc.WalletLoaderServiceClient
		walletService walletrpc.WalletServiceClient
		activeNet     *chaincfg.Params
		walletOpen    bool
	}
	tests := []struct {
		name         string
		fields       fields
		serializedTx []byte
		passphrase   string
		want         string
		wantErr      bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &WalletRPCClient{
				walletLoader:  tt.fields.walletLoader,
				walletService: tt.fields.walletService,
				activeNet:     tt.fields.activeNet,
				walletOpen:    tt.fields.walletOpen,
			}
			got, err := c.signAndPublishTransaction(tt.serializedTx, tt.passphrase)
			if (err != nil) != tt.wantErr {
				t.Errorf("WalletRPCClient.signAndPublishTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("WalletRPCClient.signAndPublishTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_processTransactions(t *testing.T) {
	tests := []struct {
		name               string
		transactionDetails []*walletrpc.TransactionDetails
		want               []*walletcore.Transaction
		wantErr            bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := processTransactions(tt.transactionDetails)
			if (err != nil) != tt.wantErr {
				t.Errorf("processTransactions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processTransactions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_processTransaction(t *testing.T) {
	tests := []struct {
		name     string
		txDetail *walletrpc.TransactionDetails
		want     *walletcore.Transaction
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := processTransaction(tt.txDetail)
			if (err != nil) != tt.wantErr {
				t.Errorf("processTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_transactionAmountAndDirection(t *testing.T) {
	tests := []struct {
		name     string
		txDetail *walletrpc.TransactionDetails
		want     int64
		want1    txhelper.TransactionDirection
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := transactionAmountAndDirection(tt.txDetail)
			if got != tt.want {
				t.Errorf("transactionAmountAndDirection() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("transactionAmountAndDirection() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
