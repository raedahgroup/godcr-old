package app

import (
	"github.com/raedahgroup/dcrlibwallet/defaultsynclistener"
	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/app/walletcore"
)

// WalletMiddleware defines key functions for interacting with a decred wallet
// These functions are implemented by the different mediums that provide access to a decred wallet
type WalletMiddleware interface {
	GenerateNewWalletSeed() (string, error)

	WalletExists() (bool, error)

	CreateWallet(passphrase, seed string) error

	IsWalletOpen() bool

	SpvSync(showLog bool, syncProgressUpdated func(*defaultsynclistener.ProgressReport))

	// todo implement for dcrwallet rpc connections
	RpcSync(showLog bool, dcrdConfig config.DcrdRpcConfig, syncProgressUpdated func(*defaultsynclistener.ProgressReport))

	// todo should introduce SpvRescan and RpcRescan
	RescanBlockChain() error

	WalletConnectionInfo() (info walletcore.ConnectionInfo, err error)

	// BestBlock fetches the best block on the network
	BestBlock() (uint32, error)

	// CloseWallet is triggered whenever the godcr program is about to be terminated
	// Usually such termination attempts are halted to allow this function perform shutdown and cleanup operations
	CloseWallet()

	DeleteWallet() error

	walletcore.Wallet
}
