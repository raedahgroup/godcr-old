package handlers

import (
	"testing"

	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app/walletcore"
)

func TestTransactionsHandler_BeforeRender(t *testing.T) {
	type fields struct {
		err                    error
		isRendering            bool
		transactions           []*walletcore.Transaction
		hasFetchedTransactions bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &TransactionsHandler{
				err:                    tt.fields.err,
				isRendering:            tt.fields.isRendering,
				transactions:           tt.fields.transactions,
				hasFetchedTransactions: tt.fields.hasFetchedTransactions,
			}
			handler.BeforeRender()
		})
	}
}

func TestTransactionsHandler_Render(t *testing.T) {
	type fields struct {
		err                    error
		isRendering            bool
		transactions           []*walletcore.Transaction
		hasFetchedTransactions bool
	}
	tests := []struct {
		name   string
		fields fields
		window *nucular.Window
		wallet walletcore.Wallet
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &TransactionsHandler{
				err:                    tt.fields.err,
				isRendering:            tt.fields.isRendering,
				transactions:           tt.fields.transactions,
				hasFetchedTransactions: tt.fields.hasFetchedTransactions,
			}
			handler.Render(tt.window, tt.wallet)
		})
	}
}

func TestTransactionsHandler_fetchTransactions(t *testing.T) {
	type fields struct {
		err                    error
		isRendering            bool
		transactions           []*walletcore.Transaction
		hasFetchedTransactions bool
	}
	tests := []struct {
		name   string
		fields fields
		wallet walletcore.Wallet
		window *nucular.Window
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &TransactionsHandler{
				err:                    tt.fields.err,
				isRendering:            tt.fields.isRendering,
				transactions:           tt.fields.transactions,
				hasFetchedTransactions: tt.fields.hasFetchedTransactions,
			}
			handler.fetchTransactions(tt.wallet, tt.window)
		})
	}
}
