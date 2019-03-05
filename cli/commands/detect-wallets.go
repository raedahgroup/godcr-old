package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/decred/dcrwallet/netparams"
	"github.com/raedahgroup/dcrlibwallet/util"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/cli/termio/terminalprompt"
	"github.com/raedahgroup/godcr/cli/walletloader"
)

type DetectWalletsCommand struct {
	commanderStub
}

func (detectCmd DetectWalletsCommand) Run(ctx context.Context) error {
	wallets, err := DetectWallets(ctx)
	if err != nil {
		return err
	}

	if wallets == nil || len(wallets) == 0 {
		return nil
	}

	fmt.Printf("Saved information for %d detected wallets\n", len(wallets))
	return nil
}

func DetectWallets(ctx context.Context) ([]*config.WalletInfo, error) {
	var allDetectedWallets []*config.WalletInfo
	for _, walletDir := range app.DecredWalletDbDirectories() {
		detectedWallets, err := findWalletsInDirectory(walletDir.Path, walletDir.Source)
		if err != nil {
			return nil, fmt.Errorf("error searching for wallets: %s", err.Error())
		}
		allDetectedWallets = append(allDetectedWallets, detectedWallets...)
	}

	// ask to create wallet if no wallet is detected
	if len(allDetectedWallets) == 0 {
		createdWallet, err := walletloader.AttemptToCreateWallet(ctx)
		if createdWallet == nil {
			return nil, err
		} else {
			allDetectedWallets = append(allDetectedWallets, createdWallet)
		}
	}

	// mark default wallet
	if len(allDetectedWallets) == 1 {
		allDetectedWallets[0].Default = true
	} else {
		promptToSelectDefaultWallet(allDetectedWallets)
	}

	// update config file with detected wallets info
	err := config.UpdateConfigFile(func(config *config.Config) {
		config.Wallets = allDetectedWallets
	})
	if err != nil {
		return nil, fmt.Errorf("failed to save %d detected wallets: %s", len(allDetectedWallets), err.Error())
	}

	return allDetectedWallets, nil
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
			DbDir:   filepath.Dir(path),
			Source:  walletSource,
			Network: netParams.Name,
		})
		return nil
	})
	return
}

func promptToSelectDefaultWallet(detectedWallets []*config.WalletInfo) []*config.WalletInfo {
	options := make([]string, len(detectedWallets))
	for i, wallet := range detectedWallets {
		options[i] = wallet.Summary()
	}

	invalidSelectionError := fmt.Errorf("invalid selection, select a number between 1 and %d", len(detectedWallets))
	terminalprompt.RequestSelection("Select default wallet", options, func(selection string) error {
		selectedIndex, err := strconv.Atoi(selection)
		if err != nil {
			return invalidSelectionError
		}

		if selectedIndex > len(detectedWallets) {
			return invalidSelectionError
		}

		detectedWallets[selectedIndex-1].Default = true
		return nil
	})

	return detectedWallets
}
