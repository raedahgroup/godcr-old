package dcrlibwallet

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet/blockchainsync"
	"github.com/raedahgroup/dcrlibwallet/utils"
	"github.com/raedahgroup/godcr/app/walletcore"
)

var numberOfPeers int32

func (lib *DcrWalletLib) GenerateNewWalletSeed() (string, error) {
	return utils.GenerateSeed()
}

func (lib *DcrWalletLib) WalletExists() (bool, error) {
	return lib.walletLib.WalletExists()
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

func (lib *DcrWalletLib) IsWalletOpen() bool {
	return lib.walletLib.WalletOpened()
}

func (lib *DcrWalletLib) SyncBlockChain(showLog bool, syncInfoUpdated func(privateSyncInfo *blockchainsync.PrivateSyncInfo, updatedSection string)) error {
	// create wrapper around syncInfoUpdated to store updated peer count before calling main syncInfoUpdated fn
	syncInfoUpdatedWrapper := func(privateSyncInfo *blockchainsync.PrivateSyncInfo, updatedSection string) {
		if updatedSection == blockchainsync.PeersCountUpdate {
			numberOfPeers = privateSyncInfo.Read().ConnectedPeers
		}
		syncInfoUpdated(privateSyncInfo, updatedSection)
	}

	// defaultSyncListener listens for actual sync updates, calculates progress and updates the caller via syncInfoUpdated
	defaultSyncListener := blockchainsync.DefaultSyncProgressListener(lib.activeNet, showLog, lib.walletLib.GetBestBlock, lib.walletLib.GetBestBlockTimeStamp,
		syncInfoUpdatedWrapper)
	lib.walletLib.AddSyncProgressListener(defaultSyncListener)

	return lib.walletLib.SpvSync("")
}

func (lib *DcrWalletLib) SyncBlockChain(showLog bool, syncInfoUpdated func(privateSyncData *sync.PrivateInfo)) error {
	syncListener := NewSyncListener(lib.activeNet, lib.walletLib, showLog, syncInfoUpdated)
	lib.walletLib.AddSyncProgressListener(syncListener)
	return lib.walletLib.SpvSync("")
}

func (lib *DcrWalletLib) RescanBlockChain() error {
	return lib.walletLib.RescanBlocks()
}

func (lib *DcrWalletLib) WalletConnectionInfo() (info walletcore.ConnectionInfo, err error) {
	accounts, loadAccountErr := lib.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if loadAccountErr != nil {
		err = fmt.Errorf("error fetching account balance: %s", loadAccountErr.Error())
		info.TotalBalance = "0 DCR"
	} else {
		var totalBalance dcrutil.Amount
		for _, acc := range accounts {
			totalBalance += acc.Balance.Total
		}
		info.TotalBalance = totalBalance.String()
	}

	bestBlock, bestBlockErr := lib.BestBlock()
	if bestBlockErr != nil && err != nil {
		err = fmt.Errorf("%s, error in fetching best block %s", err.Error(), bestBlockErr.Error())
	} else if bestBlockErr != nil {
		err = bestBlockErr
	}

	info.LatestBlock = bestBlock
	info.NetworkType = lib.NetType()
	info.PeersConnected = numberOfPeers

	return
}

func (lib *DcrWalletLib) BestBlock() (uint32, error) {
	return uint32(lib.walletLib.GetBestBlock()), nil
}

func (lib *DcrWalletLib) CloseWallet() {
	lib.walletLib.Shutdown(false)
}

func (lib *DcrWalletLib) DeleteWallet() error {
	lib.CloseWallet()
	return os.RemoveAll(lib.walletDbDir)
}
