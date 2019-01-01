package runner

import (
	"context"
	"github.com/jessevdk/go-flags"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
)

// WalletCommandRunner defines the Run method that cli commands that interact with the decred wallet must implement
// to have access to walletcore.Wallet at execution time
type WalletCommandRunner interface {
	Run(wallet walletcore.Wallet) error
	flags.Commander
}

// WalletMiddlewareCommandRunner defines the Run method that cli commands must implement to have access to app.WalletMiddleware
// in order to perform wallet creation/opening/closing and blockchain syncing operations at execution time
type WalletMiddlewareCommandRunner interface {
	Run(ctx context.Context, walletMiddleware app.WalletMiddleware) error
	flags.Commander
}

// ParserCommandRunner defines the Run method that cli commands that depends on
// flags.Parser can implement to have it injected at run time
type ParserCommandRunner interface {
	Run(parser *flags.Parser) error
	flags.Commander
}
