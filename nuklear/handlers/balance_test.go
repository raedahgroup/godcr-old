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
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := &BalanceHandler{
				err:         test.fields.err,
				isRendering: test.fields.isRendering,
				accounts:    test.fields.accounts,
				detailed:    test.fields.detailed,
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
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := &BalanceHandler{
				err:         test.fields.err,
				isRendering: test.fields.isRendering,
				accounts:    test.fields.accounts,
				detailed:    test.fields.detailed,
			}
			handler.Render(test.window, test.wallet)
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
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := &BalanceHandler{
				err:         test.fields.err,
				isRendering: test.fields.isRendering,
				accounts:    test.fields.accounts,
				detailed:    test.fields.detailed,
			}
			handler.showSimpleView(test.window)
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
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := &BalanceHandler{
				err:         test.fields.err,
				isRendering: test.fields.isRendering,
				accounts:    test.fields.accounts,
				detailed:    test.fields.detailed,
			}
			handler.showTabularView(test.window)
		})
	}
}
