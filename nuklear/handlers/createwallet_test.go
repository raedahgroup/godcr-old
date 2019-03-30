package handlers

import (
	"testing"

	"github.com/aarzilli/nucular"
)

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
