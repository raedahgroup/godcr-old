package dcrlibwallet

import (
	"context"
	"fmt"
	"github.com/raedahgroup/godcr/app"
	"time"
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

// This method may stall if the wallet database is in use by some other process,
// hence the need for ctx, so user can cancel the operation if it's taking too long
// additionally, let's notify the user if we sense a delay in opening the wallet
func (lib *DcrWalletLib) OpenWalletIfExist(ctx context.Context) (walletExists bool, err error) {
	walletOpenDelay := time.NewTicker(5 * time.Second)
	go func() {
		<-walletOpenDelay.C
		fmt.Println("It's taking longer than expected to open your wallet. " +
			"The wallet may already be opened by another app.")
	}()

	loadWalletDone := make(chan bool)
	go func() {
		defer func() {
			walletOpenDelay.Stop()
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
