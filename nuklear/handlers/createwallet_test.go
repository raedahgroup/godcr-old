package handlers

import (
	"testing"

	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app"
)

func TestCreateWalletHandler_BeforeRender(t *testing.T) {
	type fields struct {
		err                  error
		isRendering          bool
		passwordInput        nucular.TextEditor
		confirmPasswordInput nucular.TextEditor
		seedBox              nucular.TextEditor
		seed                 string
		hasStoredSeed        bool
		validationErrors     map[string]string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := &CreateWalletHandler{
				err:                  test.fields.err,
				isRendering:          test.fields.isRendering,
				passwordInput:        test.fields.passwordInput,
				confirmPasswordInput: test.fields.confirmPasswordInput,
				seedBox:              test.fields.seedBox,
				seed:                 test.fields.seed,
				hasStoredSeed:        test.fields.hasStoredSeed,
				validationErrors:     test.fields.validationErrors,
			}
			handler.BeforeRender()
		})
	}
}

func TestCreateWalletHandler_Render(t *testing.T) {
	type fields struct {
		err                  error
		isRendering          bool
		passwordInput        nucular.TextEditor
		confirmPasswordInput nucular.TextEditor
		seedBox              nucular.TextEditor
		seed                 string
		hasStoredSeed        bool
		validationErrors     map[string]string
	}
	tests := []struct {
		name       string
		fields     fields
		window     *nucular.Window
		wallet     app.WalletMiddleware
		changePage func(string)
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := &CreateWalletHandler{
				err:                  test.fields.err,
				isRendering:          test.fields.isRendering,
				passwordInput:        test.fields.passwordInput,
				confirmPasswordInput: test.fields.confirmPasswordInput,
				seedBox:              test.fields.seedBox,
				seed:                 test.fields.seed,
				hasStoredSeed:        test.fields.hasStoredSeed,
				validationErrors:     test.fields.validationErrors,
			}
			handler.Render(test.window, test.wallet, test.changePage)
		})
	}
}

func TestCreateWalletHandler_hasErrors(t *testing.T) {
	type fields struct {
		err                  error
		isRendering          bool
		passwordInput        nucular.TextEditor
		confirmPasswordInput nucular.TextEditor
		seedBox              nucular.TextEditor
		seed                 string
		hasStoredSeed        bool
		validationErrors     map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := &CreateWalletHandler{
				err:                  test.fields.err,
				isRendering:          test.fields.isRendering,
				passwordInput:        test.fields.passwordInput,
				confirmPasswordInput: test.fields.confirmPasswordInput,
				seedBox:              test.fields.seedBox,
				seed:                 test.fields.seed,
				hasStoredSeed:        test.fields.hasStoredSeed,
				validationErrors:     test.fields.validationErrors,
			}
			if got := handler.hasErrors(); got != test.want {
				t.Errorf("CreateWalletHandler.hasErrors() = %v, want %v", got, test.want)
			}
		})
	}
}
