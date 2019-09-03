package wallet

import (
	"context"

	"github.com/raedahgroup/dcrlibwallet"
)

const (
	DefaultRequiredConfirmations = 2
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

	// AccountBalance returns account balance for the accountNumbers passed in
	// or for all accounts if no account number is passed in
	AccountBalance(accountNumber uint32, requiredConfirmations int32) (*Balance, error)

	// AccountsOverview returns the name, account number and balance for all accounts in wallet
	AccountsOverview(requiredConfirmations int32) ([]*Account, error)

	// NextAccount adds an account to the wallet using the specified name
	// Returns account number for newly added account
	NextAccount(accountName string, passphrase string) (uint32, error)

	// AccountNumber looks up and returns an account number by the account's unique name
	AccountNumber(accountName string) (uint32, error)

	// AccountNumber returns the name for an account  with the provided account number
	AccountName(accountNumber uint32) (string, error)

	// AddressInfo checks if an address belongs to the wallet to retrieve it's account name
	AddressInfo(address string) (*dcrlibwallet.AddressInfo, error)

	// ValidateAddress checks if an address is valid or not
	ValidateAddress(address string) (bool, error)

	// ReceiveAddress checks if there's a previously generated address that hasn't been used to receive funds and returns it
	// If no unused address exists, it generates a new address to receive funds into specified account
	ReceiveAddress(account uint32) (string, error)

	// GenerateNewAddress generates a new address to receive funds into specified account
	// regardless of whether there was a previously generated address that has not been used
	GenerateNewAddress(account uint32) (string, error)

	// ChangePrivatePassphrase changes the private passphrase from the oldPass to the provided newPass
	ChangePrivatePassphrase(ctx context.Context, oldPass, newPass string) error

	// NetType returns the network type of this wallet
	NetType() string
}
