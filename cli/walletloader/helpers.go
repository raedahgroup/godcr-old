package walletloader

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/app/walletmediums/dcrlibwallet"
	"github.com/raedahgroup/godcr/cli/termio/terminalprompt"
)

func choseNetworkAndCreateMiddleware(ctx context.Context) (app.WalletMiddleware, error) {
	// prompt for new wallet network type and initialize
	network, err := terminalprompt.RequestInput("Which net? (mainnet, testnet)", func(input string) error {
		if strings.EqualFold(input, "mainnet") || strings.EqualFold(input, "testnet") {
			return nil
		}

		return fmt.Errorf("invalid choice, please enter 'mainnet' or 'testnet'")
	})
	if err != nil {
		return nil, fmt.Errorf("error reading your input: %s", err.Error())
	}

	if strings.EqualFold(network, "testnet") {
		network = "testnet3"
	}

	// get user-configured appdata dir to place new wallet into
	cfg, err := config.ReadConfigFile()
	if err != nil {
		return nil, fmt.Errorf("error reading appdata value from config file: %s", err.Error())
	}

	walletInfo := &config.WalletInfo{
		Network: network,
		Source:  "godcr",
		DbDir:   filepath.Join(cfg.AppDataDir, network),
	}
	return dcrlibwallet.Connect(ctx, walletInfo)
}

// displayWalletSeed prints the generated seed for a new wallet
func displayWalletSeed(seed string) {
	fmt.Println("Your wallet generation seed is:")
	fmt.Println("-------------------------------")
	seedWords := strings.Split(seed, " ")
	for i, word := range seedWords {
		fmt.Printf("%s ", word)

		if (i+1)%6 == 0 {
			fmt.Printf("\n")
		}
	}
	fmt.Println("\n-------------------------------")
	fmt.Println("IMPORTANT: Keep the seed in a safe place as you will NOT be able to restore your wallet without it.")
	fmt.Println("Please keep in mind that anyone who has access to the seed can also restore your wallet thereby " +
		"giving them access to all your funds, so it is imperative that you keep it in a secure location.")
}

func AttemptToCreateWallet(ctx context.Context) (*config.WalletInfo, error) {
	createWalletPrompt := "No wallet found. Would you like to create one now?"
	createWallet, err := terminalprompt.RequestYesNoConfirmation(createWalletPrompt, "Y")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading your response: %s", err.Error())
		return nil, err
	}

	if !createWallet {
		fmt.Println("Maybe later. Bye.")
		return nil, nil
	}

	return CreateWallet(ctx)
}
