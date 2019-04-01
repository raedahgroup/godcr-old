package dcrlibwallet

import (
	"context"
	"github.com/raedahgroup/godcr/app"
)

func (lib *DcrWalletLib) NetType() string {
	return lib.activeNet.Params.Name
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

func (lib *DcrWalletLib) OpenWalletIfExist(ctx context.Context) (walletExists bool, err error) {
	loadWalletDone := make(chan bool)

	go func() {
		defer func() {
			loadWalletDone <- true
		}()

		walletExists, err = lib.WalletExists()
		if err != nil || !walletExists {
			return
		}

		// open wallet with default public passphrase: "public"
		err = lib.walletLib.OpenWallet([]byte("public"))
	}()

	select {
	case <-loadWalletDone:
		return

	case <-ctx.Done():
		return false, ctx.Err()
	}
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
