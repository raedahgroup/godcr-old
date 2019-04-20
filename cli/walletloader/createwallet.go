package walletloader

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/raedahgroup/dcrlibwallet/utils"
	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/app/walletmediums/dcrlibwallet"
	"github.com/raedahgroup/godcr/cli/termio/terminalprompt"
)

// createWallet creates a new wallet using the dcrlibwallet WalletMiddleware.
// User is prompted to select the network type for the wallet to be created.
// If no wallet for that type already exist, user is asked to provide a private passphrase for the wallet.
// After which, a new wallet seed is generated and shown to the user and user is asked to save the wallet seed.
func createWallet(ctx context.Context, cfg *config.Config) (dcrlibwalletMiddleware *dcrlibwallet.DcrWalletLib, err error) {
	newWalletNetwork, err := requestNetworkTypeForNewWallet()
	if err != nil {
		return
	}

	// create dcrlibwallet wallet middleware and check if wallet of this type already exist
	dcrlibwalletMiddleware, err = prepareMiddlewareToCreateNewWallet(ctx, newWalletNetwork)
	if err != nil {
		return
	}

	newWalletPassphrase, err := requestNewWalletPassphrase()
	if err != nil {
		return
	}

	// get and display new wallet seed
	seed, err := generateNewWalletSeedAndAskUserToBackup()
	if err != nil {
		return
	}

	// user says they have backed up the generated wallet seed, finalize wallet creation
	err = dcrlibwalletMiddleware.CreateWallet(newWalletPassphrase, seed)
	if err != nil {
		return nil, fmt.Errorf("\nError creating wallet: %s.", err.Error())
	}
	fmt.Printf("Decred %s wallet created successfully.\n", dcrlibwalletMiddleware.NetType())

	var syncBlockchain bool

	if cfg.InterfaceMode != "cli" {
		// do not attempt to sync on cli if the user requested a different interface to be launched
		syncBlockchain = false
	} else if cfg.SyncBlockchain {
		// no need to ask user if to sync since `--sync` was already specified
		syncBlockchain = true
	} else {
		syncBlockchainPrompt := "Would you like to sync the blockchain now?"
		syncBlockchain, err = terminalprompt.RequestYesNoConfirmation(syncBlockchainPrompt, "Y")
		if err != nil {
			return dcrlibwalletMiddleware, fmt.Errorf("\nError reading your response: %s.", err.Error())
		}
	}

	if !syncBlockchain {
		return dcrlibwalletMiddleware, nil
	}

	return dcrlibwalletMiddleware, SyncBlockChain(ctx, dcrlibwalletMiddleware)
}

func requestNetworkTypeForNewWallet() (string, error) {
	// this function will be called when a user responds to the prompt to specify network type for new wallet
	checkNetworkTypeSelection := func(input string) error {
		if strings.EqualFold(input, "mainnet") || strings.EqualFold(input, "testnet") {
			return nil
		}
		return fmt.Errorf("invalid choice, please enter 'mainnet' or 'testnet'")
	}

	// prompt user to select network type for new wallet
	newWalletNetwork, err := terminalprompt.RequestInput("Which net? (mainnet, testnet)", checkNetworkTypeSelection)
	if err != nil {
		return "", fmt.Errorf("\nError getting network type for new wallet: %s.", err.Error())
	}
	if strings.EqualFold(newWalletNetwork, "testnet") {
		newWalletNetwork = "testnet3"
	}

	return newWalletNetwork, nil
}

// prepareMiddlewareToCreateNewWallet reads appdata dir from godcr.conf to use as wallet db dir.
// Also ensures that a wallet of the specified type does not already exist in the appdata dir.
func prepareMiddlewareToCreateNewWallet(ctx context.Context, newWalletNetwork string) (*dcrlibwallet.DcrWalletLib, error) {
	// get appdata dir from config to place new wallet into
	cfg, err := config.ReadConfigFile()
	if err != nil {
		return nil, fmt.Errorf("\nError reading config file to determine directory to place new wallet: %s.", err.Error())
	}
	walletDbDir := filepath.Join(cfg.AppDataDir, newWalletNetwork)
	walletMiddleware, err := dcrlibwallet.Connect(ctx, walletDbDir, newWalletNetwork)
	if err != nil {
		return nil, err
	}

	// check if wallet already exists for selected network type
	walletExists, err := walletMiddleware.WalletExists()
	if err != nil {
		return nil, fmt.Errorf("\nError checking if %s wallet already exist: %s.", walletMiddleware.NetType(), err.Error())
	}
	if walletExists {
		netType := strings.Title(walletMiddleware.NetType())
		return nil, fmt.Errorf("\n%s wallet already exist at %s.", netType, walletDbDir)
	}

	return walletMiddleware, nil
}

// requestNewWalletPassphrase asks user to enter private passphrase for new wallet twice.
// Prompt is repeated if both entered passphrases don't match.
func requestNewWalletPassphrase() (string, error) {
	for {
		passphrase, err := terminalprompt.RequestInputSecure("Enter private passphrase for new wallet", terminalprompt.EmptyValidator)
		if err != nil {
			return "", fmt.Errorf("\nError reading new wallet passphase: %s.", err.Error())
		}
		confirmPassphrase, err := terminalprompt.RequestInputSecure("Confirm passphrase", terminalprompt.EmptyValidator)
		if err != nil {
			return "", fmt.Errorf("\nError reading new wallet confirm passphase: %s.", err.Error())
		}

		if passphrase != confirmPassphrase {
			fmt.Println("Passphrases don't match, try again.")
			continue
		}

		return passphrase, nil
	}
}

func generateNewWalletSeedAndAskUserToBackup() (string, error) {
	// generate seed
	seed, err := utils.GenerateSeed()
	if err != nil {
		return "", fmt.Errorf("\nError generating seed for new wallet: %s.", err)
	}

	// display seed
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

	// ask user to back seed up
	// this function should only return after user backs up,
	// or if there's an error reading user's response to the "have you backed up?" question
	backupPrompt := `Enter "OK" to continue. This assumes you have stored the seed in a safe and secure location`
	backupValidator := func(userResponse string) error {
		userResponse = strings.TrimSpace(userResponse)
		userResponse = strings.Trim(userResponse, `"`)
		if strings.EqualFold("OK", userResponse) {
			return nil
		} else {
			return fmt.Errorf("invalid response, try again")
		}
	}
	_, err = terminalprompt.RequestInput(backupPrompt, backupValidator)
	if err != nil {
		return "", fmt.Errorf("\nError reading your response: %s.", err.Error())
	}

	return seed, nil
}
