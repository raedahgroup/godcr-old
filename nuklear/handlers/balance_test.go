package handlers

import (
	"testing"

	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app/walletcore"
)

func TestBalanceHandler_BeforeRender(t *testing.T) {
	type fields struct {
		err         error
		isRendering bool
		accounts    []*walletcore.Account
		detailed    bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &BalanceHandler{
				err:         tt.fields.err,
				isRendering: tt.fields.isRendering,
				accounts:    tt.fields.accounts,
				detailed:    tt.fields.detailed,
			}
			handler.BeforeRender()
		})
	}
}

func TestBalanceHandler_Render(t *testing.T) {
	type fields struct {
		err         error
		isRendering bool
		accounts    []*walletcore.Account
		detailed    bool
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
			handler := &BalanceHandler{
				err:         tt.fields.err,
				isRendering: tt.fields.isRendering,
				accounts:    tt.fields.accounts,
				detailed:    tt.fields.detailed,
			}
			handler.Render(tt.window, tt.wallet)
		})
	}
}

func TestBalanceHandler_showSimpleView(t *testing.T) {
	type fields struct {
		err         error
		isRendering bool
		accounts    []*walletcore.Account
		detailed    bool
	}
	tests := []struct {
		name   string
		fields fields
		window *nucular.Window
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &BalanceHandler{
				err:         tt.fields.err,
				isRendering: tt.fields.isRendering,
				accounts:    tt.fields.accounts,
				detailed:    tt.fields.detailed,
			}
			handler.showSimpleView(tt.window)
		})
	}
}

func TestBalanceHandler_showTabularView(t *testing.T) {
	type fields struct {
		err         error
		isRendering bool
		accounts    []*walletcore.Account
		detailed    bool
	}
	tests := []struct {
		name   string
		fields fields
		window *nucular.Window
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &BalanceHandler{
				err:         tt.fields.err,
				isRendering: tt.fields.isRendering,
				accounts:    tt.fields.accounts,
				detailed:    tt.fields.detailed,
			}
			handler.showTabularView(tt.window)
		})
	}
}
