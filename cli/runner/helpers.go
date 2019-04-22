package runner

import (
	"context"
	"fmt"

	flags "github.com/jessevdk/go-flags"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/cli/walletloader"
)

// CommandRequiresWallet checks if the command passed implements WalletCommandRunner, WalletMiddlewareCommandRunner
// or any other *Runner interface that requires access to a decred wallet
func CommandRequiresWallet(command flags.Commander) bool {
	if _, requiresWallet := command.(WalletCommandRunner); requiresWallet {
		return true
	}
	if _, requiresWallet := command.(WalletMiddlewareCommandRunner); requiresWallet {
		return true
	}
	return false
}

// prepareWallet ensures a wallet exists that be used and performs sync operations if requested
func prepareWallet(ctx context.Context, middleware app.WalletMiddleware, options config.CliOptions) (walletExists bool, err error) {
	// todo: main.go now requires that the user select a wallet or create one before launching interfaces, so need for this check
	//walletExists, err = middleware.WalletExists()
	//if err != nil || !walletExists {
	//	return
	//}
	walletExists = true

	if options.SyncBlockchain {
		err = walletloader.SyncBlockChain(ctx, middleware)
		if err != nil {
			return
		}
	}
	return
}

// brokenCommandError returns an error message for a command that does not have an Execute method
func brokenCommandError(command *flags.Command) error {
	return fmt.Errorf("The command %s was not properly setup.\n"+
		"Please report this bug at https://github.com/raedahgroup/godcr/issues",
		commandName(command))
}

func commandName(command *flags.Command) string {
	name := command.Name
	if command.Active != nil {
		return fmt.Sprintf("%s %s", name, commandName(command.Active))
	}
	return name
}
