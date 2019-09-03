package wallet

import (
	"github.com/raedahgroup/dcrlibwallet"
)

const (
	DefaultRequiredConfirmations = 2
)

// Wallet defines key functions for interacting with a decred wallet.
// These functions are implemented by the different mediums that provide access to a decred wallet.
type Wallet interface {
	loaderFunctions
	syncFunctions
	accountFunctions
	addressFunctions
}

// loaderFunctions holds definition for functions relating to creating,
// opening, closing and deleting a wallet.
type loaderFunctions interface {
	// WalletExists returns whether a file at the wallet database's file path exists.
	// If so, OpenWallet should be used to open the existing wallet, or CreateWallet to create a new wallet.
	WalletExists() (bool, error)

	// CreateWallet creates a new wallet using the provided seed and private passphrase.
	// Providing a previously backed-up seed performs wallet restore.
	CreateWallet(privatePass, seed string) error

	// OpenWallet opens an existing wallet database using the provided public passphrase.
	// Supply an empty passphrase to use the default public passphrase.
	// If the wallet database is already open, no error is returned.
	// This is a potentially long-running operation and care should be taken
	// to allow users cancel the op if it takes too long.
	OpenWallet(publicPass string) error

	// ChangePrivatePassphrase changes the wallet's private/spending passphrase.
	ChangePrivatePassphrase(oldPass, newPass string) error

	// NetType returns the network type of this wallet.
	NetType() string

	// Shutdown cancels any ongoing operations, unloads the wallet,
	// closes all open databases and other resources.
	Shutdown()
}

// syncFunctions holds definition for functions relating to synchronizing
// an open wallet with the network backend via Spv or Rpc.
type syncFunctions interface {
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

	// BestBlock returns the height of the tip-most block in the main
	// chain that the wallet is synchronized to.
	BestBlock() (uint32, error)
}

// accountFunctions defines functions for working with accounts in a wallet.
type accountFunctions interface {
	// Accounts returns the name, account number and balance for all accounts in the wallet.
	Accounts(requiredConfirmations int32) ([]*dcrlibwallet.Account, error)

	// AccountBalance returns account balance for the account with the specified `accountNumber`.
	AccountBalance(accountNumber uint32, requiredConfirmations int32) (*dcrlibwallet.Balance, error)

	// CreateAccount adds an account to the wallet using the specified name.
	// Returns account number for newly added account if successful.
	CreateAccount(accountName string, privatePass string) (uint32, error)

	// AccountNumber looks up and returns an account number by the account's unique name.
	AccountNumber(accountName string) (uint32, error)

	// AccountNumber returns the name for an account  with the provided `accountNumber`.
	AccountName(accountNumber uint32) (string, error)
}

// addressFunctions defines functions for working with addresses in a wallet.
type addressFunctions interface {
	// AddressInfo checks if an address belongs to the wallet and
	// returns details for the account.
	AddressInfo(address string) (*dcrlibwallet.AddressInfo, error)

	// ValidateAddress checks if an address is valid for the network type of the wallet.
	// The address does not need to belong to the wallet.
	ValidateAddress(address string) (bool, error)

	// CurrentReceiveAddress returns the last generated address that has not been used in
	// any transaction. If no unused address exists, it generates and returns a new address.
	CurrentReceiveAddress(account uint32) (string, error)

	// GenerateNewAddress returns the next address of an account branch regardless of
	// whether or not a previously generated unused address exists.
	GenerateNewAddress(account uint32) (string, error)
}
