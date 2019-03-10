package runner

import (
	"context"
	"reflect"
	"testing"

	flags "github.com/jessevdk/go-flags"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
)

func TestNew(t *testing.T) {
	type args struct {
		parser           *flags.Parser
		ctx              context.Context
		walletMiddleware app.WalletMiddleware
	}
	tests := []struct {
		name string
		args args
		want *CommandRunner
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.parser, tt.args.ctx, tt.args.walletMiddleware); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommandRunner_Run(t *testing.T) {
	type fields struct {
		parser           *flags.Parser
		ctx              context.Context
		walletMiddleware app.WalletMiddleware
	}
	type args struct {
		command flags.Commander
		args    []string
		options config.CliOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := &CommandRunner{
				parser:           tt.fields.parser,
				ctx:              tt.fields.ctx,
				walletMiddleware: tt.fields.walletMiddleware,
			}
			if err := runner.Run(tt.args.command, tt.args.args, tt.args.options); (err != nil) != tt.wantErr {
				t.Errorf("CommandRunner.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommandRunner_RunNoneWalletCommands(t *testing.T) {
	type fields struct {
		parser           *flags.Parser
		ctx              context.Context
		walletMiddleware app.WalletMiddleware
	}
	type args struct {
		command flags.Commander
		args    []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := &CommandRunner{
				parser:           tt.fields.parser,
				ctx:              tt.fields.ctx,
				walletMiddleware: tt.fields.walletMiddleware,
			}
			if err := runner.RunNoneWalletCommands(tt.args.command, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("CommandRunner.RunNoneWalletCommands() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
