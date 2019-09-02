package wallet

import (
	"context"

	"github.com/raedahgroup/dcrlibwallet"
)

// Wallet defines key functions for interacting with a decred wallet.
// These functions are implemented by the different mediums that provide access to a decred wallet.
type Wallet interface {
	// WalletExists returns whether a file at the wallet database's file path exists.
	// If so, OpenWallet should be used to open the existing wallet, or CreateWallet to create a new wallet.
	WalletExists() (bool, error)

	// CreateWallet creates a new wallet using the provided seed and private passphrase.
	// Providing a previously backed-up seed performs wallet restore.
	CreateWallet(privatePass, seed string) error

	// OpenWallet opens an existing wallet database using the provided public passphrase.
	// Supply an empty passphrase to use the default public passphrase.
	// If the wallet database is already open, no error is returned.
	OpenWallet(ctx context.Context, publicPass string) error

	// Shutdown cancels any ongoing operations, unloads the wallet,
	// closes all open databases and other resources.
	Shutdown()

	// AddSyncProgressListener registers a set of callback functions to receive sync progress updates.
	// Multiple listeners can be registered, each with a unique id.
	// Registering a listener for a unique id with a listener already registered returns an error.
	// If a sync operation is ongoing when a listener is registered,
	// the listener receives the current sync progress report as at the time of registering the listener.
	AddSyncProgressListener(syncProgressListener dcrlibwallet.SyncProgressListener, uniqueIdentifier string) error

	// RemoveSyncProgressListener de-registers a progress listener for the specified unique id.
	// Useful if the progress listener no longer serves any purpose.
	RemoveSyncProgressListener(uniqueIdentifier string)

	// SpvSync begins the spv syncing process,
	// providing progress report to all registered progress listeners.
	SpvSync(showLog bool, persistentPeers []string) error
}
