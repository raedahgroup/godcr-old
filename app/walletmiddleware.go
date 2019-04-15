package app

import (
	"context"

	"github.com/raedahgroup/dcrlibwallet/blockchainsync"
	"github.com/raedahgroup/godcr/app/walletcore"
)

// WalletMiddleware defines key functions for interacting with a decred wallet
// These functions are implemented by the different mediums that provide access to a decred wallet
type WalletMiddleware interface {
	GenerateNewWalletSeed() (string, error)

	WalletExists() (bool, error)

	CreateWallet(passphrase, seed string) error

	RescanBlockChain() error

	// OpenWalletIfExist checks if the wallet the user is trying to access exists and opens the wallet
	// This method may stall if the wallet database is in use by some other process,
	// hence the need for ctx, so user can cancel the operation if it's taking too long
	// todo: some wallets may not use default public passphrase,
	// todo: in such cases request public passphrase from user to use
	OpenWalletIfExist(ctx context.Context) (walletExists bool, err error)

	IsWalletOpen() bool

	SyncBlockChain(showLog bool, syncInfoUpdated func(privateSyncInfo *blockchainsync.PrivateSyncInfo, updatedSection string)) error

	WalletConnectionInfo() (info walletcore.ConnectionInfo, err error)

	// BestBlock fetches the best block on the network
	BestBlock() (uint32, error)

	// CloseWallet is triggered whenever the godcr program is about to be terminated
	// Usually such termination attempts are halted to allow this function perform shutdown and cleanup operations
	CloseWallet()

	DeleteWallet() error

	walletcore.Wallet
}
