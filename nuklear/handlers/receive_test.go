package handlers

import (
	"testing"

	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/helpers"
)

func TestReceiveHandler_BeforeRender(t *testing.T) {
	type fields struct {
		err                   error
		isRendering           bool
		accounts              []*walletcore.Account
		selectedAccountIndex  int
		selectedAccountNumber uint32
		generatedAddress      string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &ReceiveHandler{
				err:                   tt.fields.err,
				isRendering:           tt.fields.isRendering,
				accounts:              tt.fields.accounts,
				selectedAccountIndex:  tt.fields.selectedAccountIndex,
				selectedAccountNumber: tt.fields.selectedAccountNumber,
				generatedAddress:      tt.fields.generatedAddress,
			}
			handler.BeforeRender()
		})
	}
}

func TestReceiveHandler_Render(t *testing.T) {
	type fields struct {
		err                   error
		isRendering           bool
		accounts              []*walletcore.Account
		selectedAccountIndex  int
		selectedAccountNumber uint32
		generatedAddress      string
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
			handler := &ReceiveHandler{
				err:                   tt.fields.err,
				isRendering:           tt.fields.isRendering,
				accounts:              tt.fields.accounts,
				selectedAccountIndex:  tt.fields.selectedAccountIndex,
				selectedAccountNumber: tt.fields.selectedAccountNumber,
				generatedAddress:      tt.fields.generatedAddress,
			}
			handler.Render(tt.window, tt.wallet)
		})
	}
}

func TestReceiveHandler_RenderAddress(t *testing.T) {
	type fields struct {
		err                   error
		isRendering           bool
		accounts              []*walletcore.Account
		selectedAccountIndex  int
		selectedAccountNumber uint32
		generatedAddress      string
	}
	tests := []struct {
		name   string
		fields fields
		window *helpers.Window
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &ReceiveHandler{
				err:                   tt.fields.err,
				isRendering:           tt.fields.isRendering,
				accounts:              tt.fields.accounts,
				selectedAccountIndex:  tt.fields.selectedAccountIndex,
				selectedAccountNumber: tt.fields.selectedAccountNumber,
				generatedAddress:      tt.fields.generatedAddress,
			}
			handler.RenderAddress(tt.window)
		})
	}
}
