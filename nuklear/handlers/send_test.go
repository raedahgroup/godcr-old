package handlers

import (
	"testing"

	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app/walletcore"
)

func TestSendHandler_BeforeRender(t *testing.T) {
	type fields struct {
		err          error
		hasRendered  bool
		transactions []*walletcore.Transaction
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &SendHandler{
				err:          tt.fields.err,
				hasRendered:  tt.fields.hasRendered,
				transactions: tt.fields.transactions,
			}
			handler.BeforeRender()
		})
	}
}

func TestSendHandler_Render(t *testing.T) {
	type fields struct {
		err          error
		hasRendered  bool
		transactions []*walletcore.Transaction
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
			handler := &SendHandler{
				err:          tt.fields.err,
				hasRendered:  tt.fields.hasRendered,
				transactions: tt.fields.transactions,
			}
			handler.Render(tt.args.w, tt.args.wallet)
		})
	}
}
