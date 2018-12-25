package mobilewalletlib

import (
	"fmt"

	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrcli/core"
)

func (lib *MobileWalletLib) NetType() string {
	if lib.activeNet.Params.Name != "mainnet" {
		// could be testnet3 or testnet, return "testnet" for both cases
		return "testnet"
	} else {
		return lib.activeNet.Params.Name
	}
}

func (lib *MobileWalletLib) WalletExists() (bool, error) {
	return lib.walletLib.WalletExists()
}

func (lib *MobileWalletLib) GenerateNewWalletSeed() (string, error) {
	return lib.walletLib.GenerateSeed()
}

func (lib *MobileWalletLib) CreateWallet(passphrase, seed string) error {
	return lib.walletLib.CreateWallet(passphrase, seed)
}

func (lib *MobileWalletLib) OpenWallet() error {
	walletExists, err := lib.WalletExists()
	if err != nil {
		return err
	}

	if !walletExists {
		return fmt.Errorf("Wallet does not exist. Please create a wallet first")
	}

	// open wallet with default public passphrase: "public"
	return lib.walletLib.OpenWallet([]byte("public"))
}

func (lib *MobileWalletLib) IsWalletOpen() bool {
	return lib.walletLib.WalletOpened()
}

func (lib *MobileWalletLib) SyncBlockChain(listener *core.BlockChainSyncListener) error {
	// create wrapper around sync ended listener to deactivate logging
	originalSyncEndListener := listener.SyncEnded
	syncCompletedListener := func(err error) {
		lib.walletLib.SetLogLevel("off")
		originalSyncEndListener(err)
	}
	listener.SyncEnded = syncCompletedListener

	syncResponse := SpvSyncResponse{
		walletLib: lib.walletLib,
		listener:  listener,
		activeNet: lib.activeNet,
	}
	lib.walletLib.AddSyncResponse(syncResponse)

	// log info messages to so progress report is visible on terminal
	lib.walletLib.SetLogLevel("info")

	err := lib.walletLib.SpvSync("")
	if err != nil {
		lib.walletLib.SetLogLevel("off")
		return err
	}

	listener.SyncStarted()
	return nil
}

func (lib *MobileWalletLib) AccountBalance(accountNumber uint32) (*core.Balance, error) {
	// pass 0 as requiredConfirmations
	balance, err := lib.walletLib.GetAccountBalance(accountNumber, 0)
	if err != nil {
		return nil, err
	}

	return &core.Balance{
		Total:     dcrutil.Amount(balance.Total),
		Spendable: dcrutil.Amount(balance.Spendable),
		LockedByTickets: dcrutil.Amount(balance.LockedByTickets),
		VotingAuthority: dcrutil.Amount(balance.VotingAuthority),
		Unconfirmed:     dcrutil.Amount(balance.UnConfirmed),
	}, nil
}

func (lib *MobileWalletLib) AccountsOverview() ([]*core.Account, error) {
	// pass 0 as requiredConfirmations
	accounts, err := lib.walletLib.GetAccountsRaw(0)
	if err != nil {
		return nil, fmt.Errorf("error fetching accounts: %s", err.Error())
	}

	accountsOverview := make([]*core.Account, 0, len(accounts.Acc))

	for _, acc := range accounts.Acc {
		accountNumber := uint32(acc.Number)

		balance, err := lib.AccountBalance(accountNumber)
		if err != nil {
			return nil, err
		}

		// skip zero-balance imported accounts
		if acc.Name == "imported" && balance.Total == 0 {
			continue
		}

		account := &core.Account{
			Name:    acc.Name,
			Number:  accountNumber,
			Balance: balance,
		}
		accountsOverview = append(accountsOverview, account)
	}

	return accountsOverview, nil
}

func (lib *MobileWalletLib) NextAccount(accountName string, passphrase string) (uint32, error) {
	return 0, fmt.Errorf("not yet implemented")
}

func (lib *MobileWalletLib) AccountNumber(accountName string) (uint32, error) {
	return 0, fmt.Errorf("not yet implemented")
}

func (lib *MobileWalletLib) GenerateReceiveAddress(account uint32) (string, error) {
	return lib.walletLib.CurrentAddress(int32(account))
}

func (lib *MobileWalletLib) ValidateAddress(address string) (bool, error) {
	return lib.walletLib.IsAddressValid(address), nil
}

func (lib *MobileWalletLib) UnspentOutputs(account uint32, targetAmount int64) ([]*core.UnspentOutput, error) {
	return nil, fmt.Errorf("not yet implemented")
}

func (lib *MobileWalletLib) SendFromAccount(amountInDCR float64, sourceAccount uint32, destinationAddress, passphrase string) (string, error) {
	// convert amount from float64 DCR to int64 Atom
	amountInAtom, err := dcrutil.NewAmount(amountInDCR)
	if err != nil {
		return "", err
	}
	amount := int64(amountInAtom)

	txHash, err := lib.walletLib.SendTransaction([]byte(passphrase), destinationAddress, amount,
		int32(sourceAccount), 0, false)

	if err != nil {
		return "", nil
	}

	transactionHash, err := chainhash.NewHash(txHash)
	if err != nil {
		return "", fmt.Errorf("error parsing successful transaction hash: %s", err.Error())
	}

	return transactionHash.String(), nil
}

func (lib *MobileWalletLib) SendFromUTXOs(utxoKeys []string, dcrAmount float64, account uint32, destAddress, passphrase string) (string, error) {
	return "", fmt.Errorf("not yet implemented")
}

func (lib *MobileWalletLib) TransactionHistory() ([]*core.Transaction, error) {
	return nil, fmt.Errorf("not yet implemented")
}
