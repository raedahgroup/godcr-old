package walletloader

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/decred/dcrwallet/netparams"
	"github.com/raedahgroup/dcrlibwallet/utils"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/app/walletmediums/dcrlibwallet"
	"github.com/raedahgroup/godcr/cli/termio/terminalprompt"
)

func DetectWallets(ctx context.Context, cfg *config.Config) (*dcrlibwallet.DcrWalletLib, error) {
	var allDetectedWallets []*config.WalletInfo
	for _, walletDir := range app.DecredWalletDbDirectories() {
		detectedWallets, err := findWalletsInDirectory(walletDir.Path, walletDir.Source)
		if err != nil {
			return nil, fmt.Errorf("error searching for wallets: %s", err.Error())
		}
		allDetectedWallets = append(allDetectedWallets, detectedWallets...)
	}

	// show list of detected wallets to user to select one, or alternatively, create a new wallet
	walletList := make([]string, len(allDetectedWallets))
	for i, wallet := range allDetectedWallets {
		walletList[i] = wallet.Summary()
	}
	walletList = append(walletList, "Create new wallet")

	// this function will be called when a user responds to the prompt to select wallet
	var selectedWallet *config.WalletInfo
	validateWalletSelection := func(selection string) error {
		selectedIndex, err := strconv.Atoi(selection)
		if err != nil || selectedIndex < 1 || selectedIndex > len(allDetectedWallets) + 1 {
			return fmt.Errorf("invalid selection, select a number between 1 and %d",
				len(allDetectedWallets) + 1)
		}

		if selectedIndex <= len(allDetectedWallets) {
			selectedWallet = allDetectedWallets[selectedIndex - 1]
		}

		return nil
	}

	_, err := terminalprompt.RequestSelection("Select wallet", walletList, validateWalletSelection)
	if err != nil {
		// There was an error reading input; we cannot proceed.
		return nil, fmt.Errorf("error getting selected account: %s", err.Error())
	}

	// does the user selection correspond to a detected wallet?
	if selectedWallet != nil {
		return dcrlibwallet.Connect(ctx, selectedWallet.DbDir, selectedWallet.Network)
	}

	// user chose to create new wallet
	return createWallet(ctx, cfg)
}

func findWalletsInDirectory(walletDir, walletSource string) (wallets []*config.WalletInfo, err error) {
	// netType checks if the name of the directory where a wallet.db file was found is the name of a known/supported network type
	// dcrwallet, decredition and dcrlibwallet place wallet db files in "mainnet" or "testnet3" directories
	// returns nil if the directory used does not correspond to a known/supported network type
	detectNetParams := func(path string) *netparams.Params {
		walletDbDir := filepath.Dir(path)
		netType := filepath.Base(walletDbDir)
		return utils.NetParams(netType)
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
