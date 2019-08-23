package dcrlibwallet

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/dcrlibwallet/defaultsynclistener"
	"github.com/raedahgroup/dcrlibwallet/utils"
	"github.com/raedahgroup/godcr/app/config"
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

func (lib *DcrWalletLib) SpvSync(showLog bool, syncProgressUpdated func(*defaultsynclistener.ProgressReport)) {
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

func (lib *DcrWalletLib) RpcSync(showLog bool, dcrdConfig config.DcrdRpcConfig, syncProgressUpdated func(*defaultsynclistener.ProgressReport)) {
	// create wrapper around syncProgressUpdated to store updated peer count before calling main syncInfoUpdated fn
	syncInfoUpdatedWrapper := func(progressReport *defaultsynclistener.ProgressReport, op defaultsynclistener.SyncOp) {
		if op == defaultsynclistener.PeersCountUpdate {
			// todo peers connection/disconnection are not broadcasted on rpc sync
			numberOfPeers = progressReport.Read().ConnectedPeers
		}
		syncProgressUpdated(progressReport)
	}

	// syncListener listens for actual sync updates, calculates progress and updates the caller via syncInfoUpdated
	syncListener := defaultsynclistener.DefaultSyncProgressListener(lib.NetType(), showLog,
		lib.walletLib.GetBestBlock, lib.walletLib.GetBestBlockTimeStamp, syncInfoUpdatedWrapper)
	lib.walletLib.AddSyncProgressListener(syncListener)

	cert, err := ioutil.ReadFile(dcrdConfig.DcrdCert)
	if err != nil {
		err = fmt.Errorf("error reading dcrd cert file at %s: %s\n", dcrdConfig.DcrdCert, err.Error())
		syncListener.OnSyncError(dcrlibwallet.ErrorCodeUnexpectedError, err)
		return
	}

	err = lib.walletLib.RpcSync(dcrdConfig.DcrdHost, dcrdConfig.DcrdUser, dcrdConfig.DcrdPassword, cert)
	if err != nil {
		// todo any errors at this point won't propagate to the UI because sync has not started
		// see implementation of syncListener.OnSyncError for more
		fmt.Println("rpc sync error", err.Error())
		syncListener.OnSyncError(dcrlibwallet.ErrorCodeUnexpectedError, err)
	}
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
	return os.RemoveAll(lib.WalletDbDir)
}
