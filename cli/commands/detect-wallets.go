package commands

import (
	"fmt"
	"github.com/decred/dcrwallet/netparams"
	"github.com/raedahgroup/dcrlibwallet/util"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
	"os"
	"path/filepath"
)

type DetectWalletsCommand struct {}

func (detectCmd DetectWalletsCommand) Execute(args []string) error {
	var allDetectedWallets []*config.WalletInfo
	for _, walletDir := range app.DecredWalletDbDirectories() {
		detectedWallets, err := findWalletsInDirectory(walletDir.Path, walletDir.Source)
		if err != nil {
			return fmt.Errorf("error searching for wallets: %s", err.Error())
		}
		allDetectedWallets = append(allDetectedWallets, detectedWallets...)
	}

	return config.SaveDetectedWalletsInfo(allDetectedWallets)
}

func findWalletsInDirectory(walletDir, walletSource string) (wallets []*config.WalletInfo, err error) {
	// netType checks if the name of the directory where a wallet.db file was found is the name of a known/supported network type
	// dcrwallet, decredition and dcrlibwallet place wallet db files in "mainnet" or "testnet3" directories
	// returns nil if the directory used does not correspond to a known/supported network type
	detectNetParams := func(path string) *netparams.Params {
		walletDbDir := filepath.Dir(path)
		netType := filepath.Base(walletDbDir)
		return util.NetParams(netType)
	}

	err = filepath.Walk(walletDir, func(path string, file os.FileInfo, err error) error {
		if err != nil || file.IsDir() || file.Name() != app.WalletDbFileName {
			return nil
		}

		netParams := detectNetParams(path)
		if netParams == nil {
			return nil
		}

		wallets = append(wallets, &config.WalletInfo{
			DbPath: path,
			Source: walletSource,
			NetType: netParams.Name,
		})
		return nil
	})
	return
}
