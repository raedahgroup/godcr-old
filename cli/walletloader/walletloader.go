package walletloader

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/cli/termio/terminalprompt"
)

// createWallet creates a new wallet if one doesn't already exist using the WalletMiddleware provided
func CreateWallet(ctx context.Context, walletMiddleware app.WalletMiddleware) (err error) {
	// first check if wallet already exists
	walletExists, err := walletMiddleware.WalletExists()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking %s wallet: %s\n", walletMiddleware.NetType(), err.Error())
		return
	}
	if walletExists {
		netType := strings.Title(walletMiddleware.NetType())
		fmt.Fprintf(os.Stderr, "%s wallet already exists\n", netType)
		return fmt.Errorf("wallet already exists")
	}

	// ask user to enter passphrase twice
	passphrase, err := terminalprompt.RequestInputSecure("Enter private passphrase for new wallet", terminalprompt.EmptyValidator)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %s\n", err.Error())
		return
	}
	confirmPassphrase, err := terminalprompt.RequestInputSecure("Confirm passphrase", terminalprompt.EmptyValidator)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %s\n", err.Error())
		return
	}
	if passphrase != confirmPassphrase {
		fmt.Fprintln(os.Stderr, "Passphrases do not match")
		return fmt.Errorf("passphrases do not match")
	}

	// get seed and display to user
	seed, err := walletMiddleware.GenerateNewWalletSeed()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating seed for new wallet: %s\n", err)
		return
	}
	displayWalletSeed(seed)

	// ask user to back seed up, only proceed after user does so
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
		fmt.Fprintf(os.Stderr, "Error reading input: %s\n", err.Error())
		return
	}

	// user entered "OK" in last prompt, finalize wallet creation
	err = walletMiddleware.CreateWallet(passphrase, seed)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating wallet: %s\n", err.Error())
		return
	}
	fmt.Printf("Decred %s wallet created successfully\n", walletMiddleware.NetType())

	// sync blockchain?
	syncBlockchainPrompt := "Would you like to sync the blockchain now?"
	syncBlockchain, err := terminalprompt.RequestYesNoConfirmation(syncBlockchainPrompt, "Y")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading your response: %s\n", err.Error())
		return err
	}

	if !syncBlockchain {
		return nil
	}

	return SyncBlockChain(ctx, walletMiddleware)
}

// OpenWallet is called whenever an action to be executed requires wallet to be loaded
// notifies the program to exit if wallet doesn't exist or some other error occurs by returning a non-nil error
//
// this method may stall until previous godcr instances are closed (especially in cases of multiple dcrlibwallet instances)
// hence the need for ctx, so user can cancel the operation if it's taking too long
func OpenWallet(ctx context.Context, walletMiddleware app.WalletMiddleware) (walletExists bool, err error) {
	// notify user of the current operation so if takes too long, they have an idea what the cause is
	fmt.Println("Looking for wallets...")

	var errMsg string
	loadWalletDone := make(chan bool)

	go func() {
		defer func() {
			loadWalletDone <- true
		}()

		walletExists, err = walletMiddleware.WalletExists()
		if err != nil {
			errMsg = fmt.Sprintf("Error checking %s wallet", walletMiddleware.NetType())
			return
		}
		if !walletExists {
			return
		}

		err = walletMiddleware.OpenWallet()
		if err != nil {
			errMsg = fmt.Sprintf("Failed to open %s wallet", walletMiddleware.NetType())
		}
	}()

	select {
	case <-loadWalletDone:
		if !walletExists {
			err = attemptToCreateWallet(ctx, walletMiddleware)
			return
		}

		if errMsg != "" {
			fmt.Fprintln(os.Stderr, errMsg)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		return

	case <-ctx.Done():
		err = ctx.Err()
		return
	}
}

// syncBlockChain uses the WalletMiddleware provided to download block updates
// this is a long running operation, listen for ctx.Done and stop processing
func SyncBlockChain(ctx context.Context, walletMiddleware app.WalletMiddleware) error {
	syncDone := make(chan error)
	go func() {
		syncListener := &app.BlockChainSyncListener{
			SyncStarted: func() {
				fmt.Println("Blockchain sync started")
			},
			SyncEnded: func(err error) {
				if err == nil {
					fmt.Println("Blockchain synced successfully")
				} else {
					fmt.Fprintf(os.Stderr, "Blockchain sync completed with error: %s\n", err.Error())
				}
				syncDone <- err
			},
			OnHeadersFetched:    func(percentageProgress int64) {}, // in cli mode, sync updates are logged to terminal, no need to act on this update alert
			OnDiscoveredAddress: func(state string) {},             // in cli mode, sync updates are logged to terminal, no need to act on update alert
			OnRescanningBlocks:  func(percentageProgress int64) {}, // in cli mode, sync updates are logged to terminal, no need to act on update alert
		}

		err := walletMiddleware.SyncBlockChain(syncListener, true)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Blockchain sync failed to start. %s\n", err.Error())
			syncDone <- err
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-syncDone:
		return err
	}
}
