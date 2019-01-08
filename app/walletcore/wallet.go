package walletcore

import (
	"context"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/dcrlibwallet/txhelper"
)

// Wallet defines key functions for performing operations on a decred wallet
// These functions are implemented by the different mediums that provide access to a decred wallet
type Wallet interface {
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

	// AccountNumber returns the name for an account  with the provided account number
	AccountName(accountNumber uint32) (string, error)

	// AddressInfo checks if an address belongs to the wallet to retrieve it's account name
	AddressInfo(address string) (*txhelper.AddressInfo, error)

	// ValidateAddress checks if an address is valid or not
	ValidateAddress(address string) (bool, error)

	// GenerateReceiveAddress generates an address to receive funds into specified account
	GenerateReceiveAddress(account uint32) (string, error)

	// UnspentOutputs lists all unspent outputs in the specified account that sum up to `targetAmount`
	// If `targetAmount` is 0, all unspent outputs in account are returned
	UnspentOutputs(account uint32, targetAmount int64) ([]*UnspentOutput, error)

	// SendFromAccount sends funds to 1 or more destination addresses, each with a specified amount
	// The inputs to the transaction are automatically selected from all unspent outputs in the account
	// Returns the transaction hash as string if successful
	SendFromAccount(sourceAccount uint32, destinations []txhelper.TransactionDestination, passphrase string) (string, error)

	// SendFromUTXOs sends funds to 1 or more destination addresses, each with a specified amount
	// SendFromUTXOs also sends any change amount that arises from the transaction to the provided changeDestinations
	// The inputs to the transaction are unspent outputs in the account, matching the keys sent in []utxoKeys
	// Returns the transaction hash as string if successful
	SendFromUTXOs(sourceAccount uint32, utxoKeys []string, txDestinations []txhelper.TransactionDestination, changeDestinations []txhelper.TransactionDestination, passphrase string) (string, error)

	// TransactionHistory
	TransactionHistory() ([]*Transaction, error)

	// GetTransaction returns information about the transaction witht the given hash.
	// An error is returned if the no transaction with the given hash is found.
	GetTransaction(transactionHash string) (*TransactionDetails, error)

	// StakeInfo returns information about wallet stakes, tickets and their statuses.
	StakeInfo(ctx context.Context) (*StakeInfo, error)

	// PurchaseTicket is used to purchase tickets.
	PurchaseTicket(ctx context.Context, request dcrlibwallet.PurchaseTicketsRequest) (ticketHashes []string, err error)
}
