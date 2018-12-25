package app

import "github.com/raedahgroup/dcrcli/app/walletcore"

// WalletMiddleware defines key functions for interacting with a decred wallet
// These functions are implemented by the different mediums that provide access to a decred wallet
type WalletMiddleware interface {
	NetType() string

	WalletExists() (bool, error)

	GenerateNewWalletSeed() (string, error)

	CreateWallet(passphrase, seed string) error

	SyncBlockChain(listener *BlockChainSyncListener) error

	OpenWallet() error

	IsWalletOpen() bool

	walletcore.Wallet
}

// BlockChainSyncListener holds functions that are called during a blockchain sync operation to provide update on the sync operation
type BlockChainSyncListener struct {
	SyncStarted         func()
	SyncEnded           func(err error)
	OnHeadersFetched    func(percentageProgress int64)
	OnDiscoveredAddress func(state string)
	OnRescanningBlocks  func(percentageProgress int64)
}
