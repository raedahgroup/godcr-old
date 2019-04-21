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
	"strings"
)

type WalletInfo struct {
	DbDir   string
	Network string
	Source  string
}

func DetectWallets(ctx context.Context, cfg *config.Config) (*dcrlibwallet.DcrWalletLib, error) {
	var allDetectedWallets []*WalletInfo
	for _, walletDir := range app.DecredWalletDbDirectories() {
		detectedWallets, err := findWalletsInDirectory(walletDir.Path, walletDir.Source)
		if err != nil {
			return nil, fmt.Errorf("error searching for wallets: %s", err.Error())
		}
		allDetectedWallets = append(allDetectedWallets, detectedWallets...)
	}

	if len(allDetectedWallets) == 0 {
		walletMiddleware, err := askToCreateWallet(ctx, cfg)
		if walletMiddleware != nil {
			promptToSaveDefaultWallet(walletMiddleware.WalletDbDir)
		}
		return walletMiddleware, err
	}

	return listWalletsForSelection(ctx, cfg, allDetectedWallets)
}

func findWalletsInDirectory(walletDir, walletSource string) (wallets []*WalletInfo, err error) {
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

		wallets = append(wallets, &WalletInfo{
			DbDir:   filepath.Dir(path),
			Source:  walletSource,
			Network: netParams.Name,
		})
		return nil
	})
	return
}

func askToCreateWallet(ctx context.Context, cfg *config.Config) (*dcrlibwallet.DcrWalletLib, error) {
	prompt := "No wallets found. Do you want to create a new one?"
	shouldCreateWallet, err := terminalprompt.RequestYesNoConfirmation(prompt, "y")
	if err != nil {
		// There was an error reading input; we cannot proceed.
		return nil, fmt.Errorf("error getting selected account: %s", err.Error())
	}

	if shouldCreateWallet {
		return createWallet(ctx, cfg)
	}

	fmt.Println("Maybe later. Bye.")
	return nil, nil
}

// listWalletsForSelection shows list of detected wallets and asks user to select one, or alternatively, create a new wallet
func listWalletsForSelection(ctx context.Context, cfg *config.Config, allDetectedWallets []*WalletInfo) (*dcrlibwallet.DcrWalletLib, error) {
	// this function will be called when a user responds to the prompt to select wallet
	var selectedWallet *WalletInfo
	validateWalletSelection := func(selection string) error {
		if selection == "" || strings.EqualFold(selection, "c") {
			return nil
		}

		selectedIndex, err := strconv.Atoi(selection)
		if err != nil || selectedIndex < 1 || selectedIndex > len(allDetectedWallets) {
			if len(allDetectedWallets) == 1 {
				return fmt.Errorf("\nInvalid selection. Enter '1' or 'C'.")
			}
			return fmt.Errorf("\nInvalid selection. Enter a number between 1 and %d or enter 'C'.",
				len(allDetectedWallets))
		}

		if selectedIndex <= len(allDetectedWallets) {
			selectedWallet = allDetectedWallets[selectedIndex-1]
		}
		return nil
	}

	fmt.Println("The following wallets were found...")
	for i, wallet := range allDetectedWallets {
		fmt.Printf("(%d) %s\n", i+1, wallet.DbDir)
	}
	fmt.Println("(C)reate a new wallet.")

	for {
		response, err := terminalprompt.RequestInput("Select the wallet to use for this session", validateWalletSelection)
		if err != nil {
			// There was an error reading input; we cannot proceed.
			return nil, fmt.Errorf("error reading your response: %s", err.Error())
		}

		if response != "" {
			break
		}
	}

	// does the user selection correspond to a detected wallet?
	if selectedWallet != nil {
		// if this the only wallet, ask user if to set as default
		if len(allDetectedWallets) == 1 {
			promptToSaveDefaultWallet(selectedWallet.DbDir)
		}
		return dcrlibwallet.Connect(ctx, selectedWallet.DbDir, selectedWallet.Network)
	}

	// user chose to create new wallet
	return createWallet(ctx, cfg)
}

func promptToSaveDefaultWallet(walletDbDir string) {
	prompt := "Would you like to use this wallet by default?"
	setWalletAsDefault, err := terminalprompt.RequestYesNoConfirmation(prompt, "n")
	if err != nil {
		// error reading response, print message and continue to connect to selected wallet
		fmt.Printf("Error reading your response: %s.\n", err.Error())
		return
	}

	if setWalletAsDefault {
		err = config.UpdateConfigFile(func(config *config.ConfFileOptions) {
			config.DefaultWalletDir = walletDbDir
		})
		if err != nil {
			fmt.Printf("Error setting default wallet in config: %s.\n", err.Error())
		} else {
			fmt.Println("Default wallet selected", walletDbDir)
			fmt.Println("You can remove it by editing the config in", config.AppConfigFilePath)
		}
	}
}
