package terminal

import (
	"context"
	"reflect"
	"testing"

	"github.com/raedahgroup/godcr/app"
	"github.com/rivo/tview"
)

func TestStartTerminalApp(t *testing.T) {
	tests := []struct {
		name             string
		ctx              context.Context
		walletMiddleware app.WalletMiddleware
		wantErr          bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := StartTerminalApp(test.ctx, test.walletMiddleware); (err != nil) != test.wantErr {
				t.Errorf("StartTerminalApp() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func Test_terminalLayout(t *testing.T) {
	tests := []struct {
		name             string
		tviewApp         *tview.Application
		walletMiddleware app.WalletMiddleware
		want             tview.Primitive
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := terminalLayout(test.tviewApp, test.walletMiddleware); !reflect.DeepEqual(got, test.want) {
				t.Errorf("terminalLayout() = %v, want %v", got, test.want)
			}
		})
	}
}
