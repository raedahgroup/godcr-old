package dcrlibwallet

import (
	"fmt"
	"os"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/dcrlibwallet/defaultsynclistener"
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

func (lib *DcrWalletLib) IsWalletOpen() bool {
	return lib.walletLib.WalletOpened()
}

func (lib *DcrWalletLib) SyncBlockChain(showLog bool, syncProgressUpdated func(*defaultsynclistener.ProgressReport)) {
	// create wrapper around syncProgressUpdated to store updated peer count before calling main syncInfoUpdated fn
	syncInfoUpdatedWrapper := func(progressReport *defaultsynclistener.ProgressReport, op defaultsynclistener.SyncOp) {
		if op == defaultsynclistener.PeersCountUpdate {
			numberOfPeers = progressReport.Read().ConnectedPeers
		}
		syncProgressUpdated(progressReport)
	}

	// syncListener listens for actual sync updates, calculates progress and updates the caller via syncInfoUpdated
	syncListener := defaultsynclistener.DefaultSyncProgressListener(lib.NetType(), showLog,
		lib.walletLib.GetBestBlock, lib.walletLib.GetBestBlockTimeStamp, syncInfoUpdatedWrapper)
	lib.walletLib.AddSyncProgressListener(syncListener)

	err := lib.walletLib.SpvSync("")
	if err != nil {
		syncListener.OnSyncError(dcrlibwallet.ErrorCodeUnexpectedError, err)
	}
}

func (lib *DcrWalletLib) RescanBlockChain() error {
	return lib.walletLib.RescanBlocks(0)
}

func (lib *DcrWalletLib) IsRescanning() bool {
	return lib.walletLib.IsRescanning()
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
	return os.RemoveAll(lib.WalletDbDir)
}
