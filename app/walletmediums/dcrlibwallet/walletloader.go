package dcrlibwallet

import (
	"fmt"
	"github.com/raedahgroup/godcr/app"
)

func (lib *DcrWalletLib) NetType() string {
	if lib.activeNet.Params.Name != "mainnet" {
		// could be testnet3 or testnet, return "testnet" for both cases
		return "testnet"
	} else {
		return lib.activeNet.Params.Name
	}
}

func (lib *DcrWalletLib) WalletExists() (bool, error) {
	return lib.walletLib.WalletExists()
}

func (lib *DcrWalletLib) GenerateNewWalletSeed() (string, error) {
	return lib.walletLib.GenerateSeed()
}

func (lib *DcrWalletLib) CreateWallet(passphrase, seed string) error {
	return lib.walletLib.CreateWallet(passphrase, seed)
}

func (lib *DcrWalletLib) OpenWallet() error {
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

func (lib *DcrWalletLib) CloseWallet() {
	lib.walletLib.Shutdown(false)
}

func (lib *DcrWalletLib) IsWalletOpen() bool {
	return lib.walletLib.WalletOpened()
}

func (lib *DcrWalletLib) SyncBlockChain(listener *app.BlockChainSyncListener, showLog bool) error {
	if showLog {
		lib.walletLib.SetLogLevel("info")

		// create wrapper around sync ended listener to deactivate logging after syncing ends
		originalSyncEndedListener := listener.SyncEnded
		syncEndedListener := func(err error) {
			lib.walletLib.SetLogLevel("off")
			originalSyncEndedListener(err)
		}
		listener.SyncEnded = syncEndedListener
	}

	syncResponse := SpvSyncResponse{
		walletLib: lib.walletLib,
		listener:  listener,
		activeNet: lib.activeNet,
	}
	lib.walletLib.AddSyncResponse(syncResponse)

	err := lib.walletLib.SpvSync("")
	if err != nil {
		lib.walletLib.SetLogLevel("off")
		return err
	}

	listener.SyncStarted()
	return nil
}
