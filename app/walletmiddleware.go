package app

import (
	"context"
	"sync"

	"github.com/raedahgroup/godcr/app/walletcore"
)

// WalletMiddleware defines key functions for interacting with a decred wallet
// These functions are implemented by the different mediums that provide access to a decred wallet
type WalletMiddleware interface {
	WalletExists() (bool, error)

	GenerateNewWalletSeed() (string, error)

	CreateWallet(passphrase, seed string) error

	SyncBlockChainOld(listener *BlockChainSyncListener, showLog bool) error

	SyncBlockChain(syncInfo *SyncInfoPrivate, syncInfoUpdated func()) error

	RescanBlockChain() error

	// OpenWalletIfExist checks if the wallet the user is trying to access exists and opens the wallet
	// This method may stall if the wallet database is in use by some other process,
	// hence the need for ctx, so user can cancel the operation if it's taking too long
	// todo: some wallets may not use default public passphrase,
	// todo: in such cases request public passphrase from user to use in OpenWallet
	OpenWalletIfExist(ctx context.Context) (walletExists bool, err error)

	// CloseWallet is triggered whenever the godcr program is about to be terminated
	// Usually such termination attempts are halted to allow this function perform shutdown and cleanup operations
	CloseWallet()

	DeleteWallet() error

	IsWalletOpen() bool

	WalletConnectionInfo() (info walletcore.ConnectionInfo, err error)

	// BestBlock fetches the best block on the network
	BestBlock() (uint32, error)

	// GetConnectedPeersCount returns the number of connected peers
	GetConnectedPeersCount() int32

	walletcore.Wallet
}

// BlockChainSyncListener holds functions that are called during a blockchain sync operation to provide update on the sync operation
type BlockChainSyncListener struct {
	SyncStarted         func()
	SyncEnded           func(err error)
	OnHeadersFetched    func(percentageProgress int64)
	OnDiscoveredAddress func(state string)
	OnRescanningBlocks  func(percentageProgress int64)
	OnPeersUpdated      func(peerCount int32)
}

// SyncInfoPrivate holds information about a sync op in private variables
// to prevent reading/writing the values directly during a sync op.
type SyncInfoPrivate struct {
	sync.RWMutex
	status             SyncStatus
	bestBlockHeight    int
	currentBlockHeight int
	daysBehind         int
	connectedPeers     int32
	error              string
	done               bool
}

// syncInfo holds information about an ongoing sync op for display on the different UIs.
// Not to be used directly but with `SyncInfoPrivate`
type syncInfo struct {
	Status             SyncStatus
	BestBlockHeight    int
	CurrentBlockHeight int
	DaysBehind         int
	ConnectedPeers     int32
	Error              string
	Done               bool
}

// Read returns the current sync op info from private variables after locking the mutex for reading
func (s *SyncInfoPrivate) Read() *syncInfo {
	s.RLock()
	defer s.RUnlock()

	return &syncInfo{
		s.status,
		s.bestBlockHeight,
		s.currentBlockHeight,
		s.daysBehind,
		s.connectedPeers,
		s.error,
		s.done,
	}
}

// Write saves info for ongoing sync op to private variables after locking mutex for writing
func (s *SyncInfoPrivate) Write(info *syncInfo, status SyncStatus) {
	s.Lock()
	defer s.Unlock()

	s.status = status
	s.bestBlockHeight = info.BestBlockHeight
	s.currentBlockHeight = info.CurrentBlockHeight
	s.daysBehind = info.DaysBehind
	s.connectedPeers = info.ConnectedPeers
	s.error = info.Error
	s.done = info.Done
}

type SyncStatus uint8

const (
	SyncStatusNotStarted SyncStatus = iota
	SyncStatusSuccess
	SyncStatusError
	SyncStatusInProgress
)
