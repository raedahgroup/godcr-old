package wallet

import (
	"github.com/raedahgroup/dcrlibwallet/defaultsynclistener"
)

// WalletMiddleware defines key functions for interacting with a decred wallet
// These functions are implemented by the different mediums that provide access to a decred wallet
type WalletMiddleware interface {
	GenerateNewWalletSeed() (string, error)

	WalletExists() (bool, error)

	CreateWallet(passphrase, seed string) error

	IsWalletOpen() bool

	SyncBlockChain(showLog bool, syncProgressUpdated func(*defaultsynclistener.ProgressReport))

	RescanBlockChain() error

	WalletConnectionInfo() (info wallet.ConnectionInfo, err error)

	// BestBlock fetches the best block on the network
	BestBlock() (uint32, error)

	// CloseWallet is triggered whenever the godcr program is about to be terminated
	// Usually such termination attempts are halted to allow this function perform shutdown and cleanup operations
	CloseWallet()

	DeleteWallet() error

	Wallet
}
