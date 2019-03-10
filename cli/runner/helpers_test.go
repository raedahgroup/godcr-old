package runner

import (
	"context"
	"testing"

	flags "github.com/jessevdk/go-flags"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
)

func TestCommandRequiresWallet(t *testing.T) {
	tests := []struct {
		name    string
		command flags.Commander
		want    bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := CommandRequiresWallet(test.command); got != test.want {
				t.Errorf("CommandRequiresWallet() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_prepareWallet(t *testing.T) {
	tests := []struct {
		name             string
		ctx              context.Context
		middleware       app.WalletMiddleware
		options          config.CliOptions
		wantWalletExists bool
		wantErr          bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotWalletExists, err := prepareWallet(test.ctx, test.middleware, test.options)
			if (err != nil) != test.wantErr {
				t.Errorf("prepareWallet() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if gotWalletExists != test.wantWalletExists {
				t.Errorf("prepareWallet() = %v, want %v", gotWalletExists, test.wantWalletExists)
			}
		})
	}
}

func Test_brokenCommandError(t *testing.T) {
	tests := []struct {
		name    string
		command *flags.Command
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := brokenCommandError(test.command); (err != nil) != test.wantErr {
				t.Errorf("brokenCommandError() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func Test_commandName(t *testing.T) {
	tests := []struct {
		name    string
		command *flags.Command
		want    string
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := commandName(test.command); got != test.want {
				t.Errorf("commandName() = %v, want %v", got, test.want)
			}
		})
	}
}
