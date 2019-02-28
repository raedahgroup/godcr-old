package dcrwalletrpc

import (
	"reflect"
	"testing"
)

func Test_connectionParamsFromDcrwalletConfig(t *testing.T) {
	tests := []struct {
		name          string
		wantAddresses []string
		wantNotls     bool
		wantCertpath  string
		wantErr       bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAddresses, gotNotls, gotCertpath, err := connectionParamsFromDcrwalletConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("connectionParamsFromDcrwalletConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotAddresses, tt.wantAddresses) {
				t.Errorf("connectionParamsFromDcrwalletConfig() gotAddresses = %v, want %v", gotAddresses, tt.wantAddresses)
			}
			if gotNotls != tt.wantNotls {
				t.Errorf("connectionParamsFromDcrwalletConfig() gotNotls = %v, want %v", gotNotls, tt.wantNotls)
			}
			if gotCertpath != tt.wantCertpath {
				t.Errorf("connectionParamsFromDcrwalletConfig() gotCertpath = %v, want %v", gotCertpath, tt.wantCertpath)
			}
		})
	}
}
