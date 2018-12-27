package mobilewalletlib

import (
	"fmt"
	"github.com/raedahgroup/dcrcli/app"
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

func (lib *MobileWalletLib) CloseWallet() {
	lib.walletLib.Shutdown()
}

func (lib *MobileWalletLib) IsWalletOpen() bool {
	return lib.walletLib.WalletOpened()
}

func (lib *MobileWalletLib) SyncBlockChain(listener *app.BlockChainSyncListener) error {
	// create wrapper around sync ended listener to deactivate logging after syncing ends and publish any pending transactions
	originalSyncEndedListener := listener.SyncEnded
	syncEndedListener := func(err error) {
		lib.walletLib.SetLogLevel("off")
		lib.walletLib.PublishUnminedTransactions()
		originalSyncEndedListener(err)
	}
	listener.SyncEnded = syncEndedListener

	syncResponse := SpvSyncResponse{
		walletLib: lib.walletLib,
		listener:  listener,
		activeNet: lib.activeNet,
	}
	lib.walletLib.AddSyncResponse(syncResponse)

	// log info messages to show progress report on terminal
	lib.walletLib.SetLogLevel("info")

	err := lib.walletLib.SpvSync("")
	if err != nil {
		lib.walletLib.SetLogLevel("off")
		return err
	}

	listener.SyncStarted()
	return nil
}
