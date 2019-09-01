package wallet

// Wallet defines key functions for interacting with a decred wallet.
// These functions are implemented by the different mediums that provide access to a decred wallet.
type Wallet interface {
	// CreateWallet creates a new wallet using the provided seed and private passphrase.
	// Providing a previously backed-up seed performs wallet restore.
	CreateWallet(seed, privatePass string) error

	// OpenWallet opens an existing wallet database using the provided public passphrase.
	// Supply an empty passphrase to use the default public passphrase.
	OpenWallet(publicPass string) error

	// Shutdown cancels any ongoing operations, unloads the wallet,
	// closes all open databases and other resources.
	Shutdown()

	// SpvSync begins the spv syncing process,
	// providing progress report via the provided progress listener.
	SpvSync(showLog bool)

	// RpcSync begins the syncing process via rpc connection to a daemon,
	// providing progress report via the provided progress listener.
	RpcSync(showLog bool)
}
