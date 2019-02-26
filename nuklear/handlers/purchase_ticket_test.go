package handlers

import (
	"testing"

	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app/walletcore"
)

func TestPurchaseTicketsHandler_BeforeRender(t *testing.T) {
	type fields struct {
		err                   error
		isRendering           bool
		numTicketsInput       nucular.TextEditor
		numTicketsInputErrStr string
		spendUnconfirmed      bool
		accountNumbers        []uint32
		accountOverviews      []string
		selectedAccountIndex  int
		isSubmitting          bool
		ticketHashes          []string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &PurchaseTicketsHandler{
				err:                   tt.fields.err,
				isRendering:           tt.fields.isRendering,
				numTicketsInput:       tt.fields.numTicketsInput,
				numTicketsInputErrStr: tt.fields.numTicketsInputErrStr,
				spendUnconfirmed:      tt.fields.spendUnconfirmed,
				accountNumbers:        tt.fields.accountNumbers,
				accountOverviews:      tt.fields.accountOverviews,
				selectedAccountIndex:  tt.fields.selectedAccountIndex,
				isSubmitting:          tt.fields.isSubmitting,
				ticketHashes:          tt.fields.ticketHashes,
			}
			handler.BeforeRender()
		})
	}
}

func TestPurchaseTicketsHandler_Render(t *testing.T) {
	type fields struct {
		err                   error
		isRendering           bool
		numTicketsInput       nucular.TextEditor
		numTicketsInputErrStr string
		spendUnconfirmed      bool
		accountNumbers        []uint32
		accountOverviews      []string
		selectedAccountIndex  int
		isSubmitting          bool
		ticketHashes          []string
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
			handler := &PurchaseTicketsHandler{
				err:                   tt.fields.err,
				isRendering:           tt.fields.isRendering,
				numTicketsInput:       tt.fields.numTicketsInput,
				numTicketsInputErrStr: tt.fields.numTicketsInputErrStr,
				spendUnconfirmed:      tt.fields.spendUnconfirmed,
				accountNumbers:        tt.fields.accountNumbers,
				accountOverviews:      tt.fields.accountOverviews,
				selectedAccountIndex:  tt.fields.selectedAccountIndex,
				isSubmitting:          tt.fields.isSubmitting,
				ticketHashes:          tt.fields.ticketHashes,
			}
			handler.Render(tt.window, tt.wallet)
		})
	}
}

func TestPurchaseTicketsHandler_fetchAccounts(t *testing.T) {
	type fields struct {
		err                   error
		isRendering           bool
		numTicketsInput       nucular.TextEditor
		numTicketsInputErrStr string
		spendUnconfirmed      bool
		accountNumbers        []uint32
		accountOverviews      []string
		selectedAccountIndex  int
		isSubmitting          bool
		ticketHashes          []string
	}
	tests := []struct {
		name   string
		fields fields
		wallet walletcore.Wallet
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &PurchaseTicketsHandler{
				err:                   tt.fields.err,
				isRendering:           tt.fields.isRendering,
				numTicketsInput:       tt.fields.numTicketsInput,
				numTicketsInputErrStr: tt.fields.numTicketsInputErrStr,
				spendUnconfirmed:      tt.fields.spendUnconfirmed,
				accountNumbers:        tt.fields.accountNumbers,
				accountOverviews:      tt.fields.accountOverviews,
				selectedAccountIndex:  tt.fields.selectedAccountIndex,
				isSubmitting:          tt.fields.isSubmitting,
				ticketHashes:          tt.fields.ticketHashes,
			}
			handler.fetchAccounts(tt.wallet)
		})
	}
}

func TestPurchaseTicketsHandler_validateAndSubmit(t *testing.T) {
	type fields struct {
		err                   error
		isRendering           bool
		numTicketsInput       nucular.TextEditor
		numTicketsInputErrStr string
		spendUnconfirmed      bool
		accountNumbers        []uint32
		accountOverviews      []string
		selectedAccountIndex  int
		isSubmitting          bool
		ticketHashes          []string
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
			handler := &PurchaseTicketsHandler{
				err:                   tt.fields.err,
				isRendering:           tt.fields.isRendering,
				numTicketsInput:       tt.fields.numTicketsInput,
				numTicketsInputErrStr: tt.fields.numTicketsInputErrStr,
				spendUnconfirmed:      tt.fields.spendUnconfirmed,
				accountNumbers:        tt.fields.accountNumbers,
				accountOverviews:      tt.fields.accountOverviews,
				selectedAccountIndex:  tt.fields.selectedAccountIndex,
				isSubmitting:          tt.fields.isSubmitting,
				ticketHashes:          tt.fields.ticketHashes,
			}
			handler.validateAndSubmit(tt.window, tt.wallet)
		})
	}
}

func TestPurchaseTicketsHandler_submit(t *testing.T) {
	type fields struct {
		err                   error
		isRendering           bool
		numTicketsInput       nucular.TextEditor
		numTicketsInputErrStr string
		spendUnconfirmed      bool
		accountNumbers        []uint32
		accountOverviews      []string
		selectedAccountIndex  int
		isSubmitting          bool
		ticketHashes          []string
	}
	tests := []struct {
		name       string
		fields     fields
		passphrase string
		window     *nucular.Window
		wallet     walletcore.Wallet
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &PurchaseTicketsHandler{
				err:                   tt.fields.err,
				isRendering:           tt.fields.isRendering,
				numTicketsInput:       tt.fields.numTicketsInput,
				numTicketsInputErrStr: tt.fields.numTicketsInputErrStr,
				spendUnconfirmed:      tt.fields.spendUnconfirmed,
				accountNumbers:        tt.fields.accountNumbers,
				accountOverviews:      tt.fields.accountOverviews,
				selectedAccountIndex:  tt.fields.selectedAccountIndex,
				isSubmitting:          tt.fields.isSubmitting,
				ticketHashes:          tt.fields.ticketHashes,
			}
			handler.submit(tt.passphrase, tt.window, tt.wallet)
		})
	}
}

func TestPurchaseTicketsHandler_resetForm(t *testing.T) {
	type fields struct {
		err                   error
		isRendering           bool
		numTicketsInput       nucular.TextEditor
		numTicketsInputErrStr string
		spendUnconfirmed      bool
		accountNumbers        []uint32
		accountOverviews      []string
		selectedAccountIndex  int
		isSubmitting          bool
		ticketHashes          []string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &PurchaseTicketsHandler{
				err:                   tt.fields.err,
				isRendering:           tt.fields.isRendering,
				numTicketsInput:       tt.fields.numTicketsInput,
				numTicketsInputErrStr: tt.fields.numTicketsInputErrStr,
				spendUnconfirmed:      tt.fields.spendUnconfirmed,
				accountNumbers:        tt.fields.accountNumbers,
				accountOverviews:      tt.fields.accountOverviews,
				selectedAccountIndex:  tt.fields.selectedAccountIndex,
				isSubmitting:          tt.fields.isSubmitting,
				ticketHashes:          tt.fields.ticketHashes,
			}
			handler.resetForm()
		})
	}
}
