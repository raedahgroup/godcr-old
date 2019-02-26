package handlers

import (
	"testing"

	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app/walletcore"
)

func TestStakeInfoHandler_BeforeRender(t *testing.T) {
	type fields struct {
		err         error
		isRendering bool
		stakeInfo   *walletcore.StakeInfo
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &StakeInfoHandler{
				err:         tt.fields.err,
				isRendering: tt.fields.isRendering,
				stakeInfo:   tt.fields.stakeInfo,
			}
			handler.BeforeRender()
		})
	}
}

func TestStakeInfoHandler_Render(t *testing.T) {
	type fields struct {
		err         error
		isRendering bool
		stakeInfo   *walletcore.StakeInfo
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
			handler := &StakeInfoHandler{
				err:         tt.fields.err,
				isRendering: tt.fields.isRendering,
				stakeInfo:   tt.fields.stakeInfo,
			}
			handler.Render(tt.window, tt.wallet)
		})
	}
}
