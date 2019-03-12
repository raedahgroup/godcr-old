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
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := &PurchaseTicketsHandler{
				err:                   test.fields.err,
				isRendering:           test.fields.isRendering,
				numTicketsInput:       test.fields.numTicketsInput,
				numTicketsInputErrStr: test.fields.numTicketsInputErrStr,
				spendUnconfirmed:      test.fields.spendUnconfirmed,
				accountNumbers:        test.fields.accountNumbers,
				accountOverviews:      test.fields.accountOverviews,
				selectedAccountIndex:  test.fields.selectedAccountIndex,
				isSubmitting:          test.fields.isSubmitting,
				ticketHashes:          test.fields.ticketHashes,
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
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := &PurchaseTicketsHandler{
				err:                   test.fields.err,
				isRendering:           test.fields.isRendering,
				numTicketsInput:       test.fields.numTicketsInput,
				numTicketsInputErrStr: test.fields.numTicketsInputErrStr,
				spendUnconfirmed:      test.fields.spendUnconfirmed,
				accountNumbers:        test.fields.accountNumbers,
				accountOverviews:      test.fields.accountOverviews,
				selectedAccountIndex:  test.fields.selectedAccountIndex,
				isSubmitting:          test.fields.isSubmitting,
				ticketHashes:          test.fields.ticketHashes,
			}
			handler.Render(test.window, test.wallet)
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
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := &PurchaseTicketsHandler{
				err:                   test.fields.err,
				isRendering:           test.fields.isRendering,
				numTicketsInput:       test.fields.numTicketsInput,
				numTicketsInputErrStr: test.fields.numTicketsInputErrStr,
				spendUnconfirmed:      test.fields.spendUnconfirmed,
				accountNumbers:        test.fields.accountNumbers,
				accountOverviews:      test.fields.accountOverviews,
				selectedAccountIndex:  test.fields.selectedAccountIndex,
				isSubmitting:          test.fields.isSubmitting,
				ticketHashes:          test.fields.ticketHashes,
			}
			handler.fetchAccounts(test.wallet)
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
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := &PurchaseTicketsHandler{
				err:                   test.fields.err,
				isRendering:           test.fields.isRendering,
				numTicketsInput:       test.fields.numTicketsInput,
				numTicketsInputErrStr: test.fields.numTicketsInputErrStr,
				spendUnconfirmed:      test.fields.spendUnconfirmed,
				accountNumbers:        test.fields.accountNumbers,
				accountOverviews:      test.fields.accountOverviews,
				selectedAccountIndex:  test.fields.selectedAccountIndex,
				isSubmitting:          test.fields.isSubmitting,
				ticketHashes:          test.fields.ticketHashes,
			}
			handler.validateAndSubmit(test.window, test.wallet)
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
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := &PurchaseTicketsHandler{
				err:                   test.fields.err,
				isRendering:           test.fields.isRendering,
				numTicketsInput:       test.fields.numTicketsInput,
				numTicketsInputErrStr: test.fields.numTicketsInputErrStr,
				spendUnconfirmed:      test.fields.spendUnconfirmed,
				accountNumbers:        test.fields.accountNumbers,
				accountOverviews:      test.fields.accountOverviews,
				selectedAccountIndex:  test.fields.selectedAccountIndex,
				isSubmitting:          test.fields.isSubmitting,
				ticketHashes:          test.fields.ticketHashes,
			}
			handler.submit(test.passphrase, test.window, test.wallet)
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
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := &PurchaseTicketsHandler{
				err:                   test.fields.err,
				isRendering:           test.fields.isRendering,
				numTicketsInput:       test.fields.numTicketsInput,
				numTicketsInputErrStr: test.fields.numTicketsInputErrStr,
				spendUnconfirmed:      test.fields.spendUnconfirmed,
				accountNumbers:        test.fields.accountNumbers,
				accountOverviews:      test.fields.accountOverviews,
				selectedAccountIndex:  test.fields.selectedAccountIndex,
				isSubmitting:          test.fields.isSubmitting,
				ticketHashes:          test.fields.ticketHashes,
			}
			handler.resetForm()
		})
	}
}
