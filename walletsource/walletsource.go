package walletsource

// WalletSource interface defines the key functions that are implemented
// by the different mediums for connecting to a dcr wallet
// Individual mediums may expose other functions but must implement these
type WalletSource interface {
	NetType() string

	WalletExists() (bool, error)

	GenerateNewWalletSeed() (string, error)

	CreateWallet(passphrase, seed string) error

	SyncBlockChain(listener *BlockChainSyncListener) error

	OpenWallet() error

	IsWalletOpen() bool

	// Balance returns account balance for the accountNumbers passed in
	// or for all accounts if no account number is passed in
	AccountBalance(accountNumber uint32) (*Balance, error)

	// AccountsOverview returns the name, account number and balance for all accounts in wallet
	AccountsOverview() ([]*Account, error)

	// NextAccount adds an account to the wallet using the specified name
	// Returns account number for newly added account
	NextAccount(accountName string, passphrase string) (uint32, error)

	// AccountNumber looks up and returns an account number by the account's unique name
	AccountNumber(accountName string) (uint32, error)

	// GenerateReceiveAddress generates an address to receive funds into specified account
	GenerateReceiveAddress(account uint32) (string, error)

	// ValidateAddress checks if an address is valid or not
	ValidateAddress(address string) (bool, error)

	// UnspentOutputs lists all unspent outputs in the specified account that sum up to `targetAmount`
	// If `targetAmount` is 0, all unspent outputs in account are returned
	UnspentOutputs(account uint32, targetAmount int64) ([]*UnspentOutput, error)

	// SendFromAccount sends funds to the destination address
	// by automatically selecting 1 or more unspent outputs from the specified account
	// Returns the transaction hash as string if successful
	SendFromAccount(amountInDCR float64, sourceAccount uint32, destinationAddress, passphrase string) (string, error)

	// UTXOSend sends funds to the destination address using unspent outputs matching the keys sent in []utxoKeys
	// Returns the transaction hash as string if successful
	SendFromUTXOs(utxoKeys []string, amountInDCR float64, sourceAccount uint32, destinationAddress, passphrase string) (string, error)

	// TransactionHistory
	TransactionHistory() ([]*Transaction, error)
}
