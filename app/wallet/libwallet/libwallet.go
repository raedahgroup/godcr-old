package libwallet

import (
	"context"
	"time"
	"fmt"

	"github.com/decred/dcrwallet/netparams"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/dcrlibwallet/utils"
)

// LibWallet implements `wallet.Wallet` using `dcrlibwallet.LibWallet`
// as medium for connecting to a decred wallet.
type LibWallet struct {
	WalletDbDir string // todo confirm if this is needed
	activeNet   *netparams.Params
	lw   *dcrlibwallet.LibWallet
}

// todo: correct this doc
// Init opens connection to the wallet database via dcrlibwallet and returns an instance of LibWallet
func Init(walletDbDir, networkType string) (*LibWallet, error) {
	activeNet := utils.NetParams(networkType)
	if activeNet == nil {
		return nil, fmt.Errorf("unsupported wallet: %s", networkType)
	}

	dcrlibwallet.SetLogLevels("off")
	lw, err := dcrlibwallet.NewLibWalletWithDbPath(walletDbDir, activeNet)
	if err != nil {
		return nil, err
	}

	return &LibWallet{
		WalletDbDir: walletDbDir,
		lw:   lw,
		activeNet:   activeNet,
	}, nil
}

// This method may stall if the wallet database is in use by some other process,
// hence the need for ctx, so user can cancel the operation if it's taking too long
// additionally, let's notify the user if we sense a delay in opening the wallet
func openWalletIfExist(ctx context.Context, walletLib *dcrlibwallet.LibWallet) error {
	// wallet database is opened using bolt db by `github.com/decred/dcrwallet/wallet/internal/bdb`
	// bolt db stalls if the database is currently in use by another process,
	// waiting for the other process to release the file.
	// bold db doc advise setting a 1 second timeout to prevent this stalling.
	// see https://github.com/boltdb/bolt#opening-a-database
	walletOpenDelay := time.NewTicker(5 * time.Second)

	loadWalletDone := make(chan error)
	go func() {
		var openWalletError error
		defer func() {
			loadWalletDone <- openWalletError
			walletOpenDelay.Stop()
		}()

		walletExists, openWalletError := walletLib.WalletExists()
		if openWalletError != nil || !walletExists {
			return
		}

		// open wallet with default public passphrase: "public"
		openWalletError = walletLib.OpenWallet([]byte("public"))
	}()

	select {
	case <-walletOpenDelay.C:
		return fmt.Errorf("\nWallet database is in use by another process.")

	case err := <-loadWalletDone:
		return err

	case <-ctx.Done():
		return ctx.Err()
	}
}
