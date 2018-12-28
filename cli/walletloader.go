package cli

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/raedahgroup/dcrcli/app"
	"github.com/raedahgroup/dcrcli/cli/terminalprompt"
)

// createWallet creates a new wallet if one doesn't already exist using the WalletMiddleware provided
func createWallet(walletMiddleware app.WalletMiddleware) (err error) {
	// first check if wallet already exists
	walletExists, err := walletMiddleware.WalletExists()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking %s wallet: %s\n", walletMiddleware.NetType(), err.Error())
		return
	}
	if walletExists {
		netType := strings.Title(walletMiddleware.NetType())
		fmt.Fprintf(os.Stderr, "%s wallet already exists", netType)
		return fmt.Errorf("wallet already exists")
	}

	// ask user to enter passphrase twice
	passphrase, err := terminalprompt.RequestInputSecure("Enter private passphrase for new wallet", terminalprompt.EmptyValidator)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %s", err.Error())
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
		fmt.Fprintf(os.Stderr, "Error reading input: %s", err.Error())
		return
	}

	// user entered "OK" in last prompt, finalize wallet creation
	err = walletMiddleware.CreateWallet(passphrase, seed)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating wallet: %s", err.Error())
		return
	}

	fmt.Println("Your wallet has been created successfully")

	// perform first blockchain sync after creating wallet
	return syncBlockChain(walletMiddleware)
}

// openWallet is called whenever an action to be executed requires wallet to be loaded
// exits the program if wallet doesn't exist or some other error occurs
//
// this method may stall until previous dcrcli instances are closed (especially in cases of multiple mobilewallet instances)
// hence the need for ctx, so user can cancel the operation if it's taking too long
func openWallet(ctx context.Context, walletMiddleware app.WalletMiddleware) error {
	// notify user of the current operation so if takes too long, they have an idea what the cause is
	fmt.Println("Opening wallet...")

	var err error
	var errMsg string
	loadWalletDone := make(chan bool)

	go func() {
		defer func() {
			loadWalletDone <- true
		}()

		var walletExists bool
		walletExists, err = walletMiddleware.WalletExists()
		if err != nil {
			errMsg = fmt.Sprintf("Error checking %s wallet", walletMiddleware.NetType())
			return
		}

		if !walletExists {
			netType := strings.Title(walletMiddleware.NetType())
			errMsg = fmt.Sprintf("%s wallet does not exist. Create it using '%s --createwallet'", netType, app.Name())
			return
		}

		err = walletMiddleware.OpenWallet()
		if err != nil {
			errMsg = fmt.Sprintf("Failed to open %s wallet", walletMiddleware.NetType())
		}
	}()

	select {
	case <-loadWalletDone:
		if errMsg != "" {
			fmt.Fprintln(os.Stderr, errMsg)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		return err

	case <-ctx.Done():
		return ctx.Err()
	}
}

// syncBlockChain uses the WalletMiddleware provided to download block updates
func syncBlockChain(walletMiddleware app.WalletMiddleware) (err error) {
	// use wait group to wait for go routine process to complete before exiting this function
	var wg sync.WaitGroup
	wg.Add(1)

	err = walletMiddleware.SyncBlockChain(&app.BlockChainSyncListener{
		SyncStarted: func() {
			fmt.Println("Blockchain sync started")
		},
		SyncEnded: func(e error) {
			err = e
			if err == nil {
				fmt.Println("Blockchain synced successfully")
			} else {
				fmt.Fprintf(os.Stderr, "Blockchain sync completed with error: %s", err.Error())
			}
			wg.Done()
		},
		OnHeadersFetched:    func(percentageProgress int64) {}, // in cli mode, sync updates are logged to terminal, no need to act on this update alert
		OnDiscoveredAddress: func(state string) {},             // in cli mode, sync updates are logged to terminal, no need to act on update alert
		OnRescanningBlocks:  func(percentageProgress int64) {}, // in cli mode, sync updates are logged to terminal, no need to act on update alert
	}, true)

	if err != nil {
		// sync go routine failed to start, nothing to wait for
		wg.Done()
	} else {
		// sync in progress, wait for BlockChainSyncListener.OnComplete
		wg.Wait()
	}
	return
}
