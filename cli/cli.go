package cli

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/raedahgroup/dcrcli/cli/terminalprompt"
	"github.com/raedahgroup/dcrcli/config"
	"github.com/raedahgroup/dcrcli/core"
)

type Response struct {
	Columns []string
	Result  [][]interface{}
}

var (
	Wallet core.Wallet
	StdoutWriter = tabWriter(os.Stdout)
)

func CreateWallet() {
	// no need to make the user go through stress of providing following info if wallet already exists
	walletExists, err := Wallet.WalletExists()
	if err != nil {
		errMsg := fmt.Sprintf("Error checking %s wallet", Wallet.NetType())
		printErrorAndExit(errMsg, err)
	}

	if walletExists {
		netType := strings.Title(Wallet.NetType())
		errMsg := fmt.Sprintf("%s wallet already exists", netType)
		printErrorAndExit(errMsg, nil)
	}

	// ask user to enter passphrase twice
	passphrase, err := terminalprompt.RequestInputSecure("Enter private passphrase for new wallet", terminalprompt.EmptyValidator)
	if err != nil {
		printErrorAndExit("Error reading input", err)
	}

	confirmPassphrase, err := terminalprompt.RequestInputSecure("Confirm passphrase", terminalprompt.EmptyValidator)
	if err != nil {
		printErrorAndExit("Error reading input", err)
	}

	if passphrase != confirmPassphrase {
		printErrorAndExit("Passphrases do not match", nil)
	}

	// get seed and display to user to save/backup
	seed, err := Wallet.GenerateNewWalletSeed()
	if err != nil {
		printErrorAndExit("Error generating seed for new wallet", err)
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
		printErrorAndExit("Error reading input", err)
	}

	// user entered "OK" in last prompt, finalize wallet creation
	err = Wallet.CreateWallet(passphrase, seed)
	if err != nil {
		printErrorAndExit("Error creating wallet", err)
	}

	fmt.Println("Your wallet has been created successfully.")
}

// called whenever an action to be executed requires wallet to be loaded
// exits the program is wallet doesn't exist or some other error occurs
func OpenWallet() {
	walletExists, err := Wallet.WalletExists()
	if err != nil {
		errMsg := fmt.Sprintf("Error checking %s wallet", Wallet.NetType())
		printErrorAndExit(errMsg, err)
	}

	if !walletExists {
		netType := strings.Title(Wallet.NetType())
		errMsg := fmt.Sprintf("%s wallet does not exist. Use '%s create' to create a wallet", netType, config.AppName())
		printErrorAndExit(errMsg, nil)
	}

	err = Wallet.OpenWallet()
	if err != nil {
		errMsg := fmt.Sprintf("Failed to open %s wallet", Wallet.NetType())
		printErrorAndExit(errMsg, err)
	}
}

// syncBlockChain registers a progress listener with core to download block updates
// causes app to exit if an error is encountered
func SyncBlockChain() {
	var err error
	defer func() {
		if err != nil {
			printErrorAndExit("Error syncing blockchain", err)
		} else {
			fmt.Println("Blockchain synced successfully")
		}
	}()

	// use wait group to wait for go routine process to complete before exiting this function
	var wg sync.WaitGroup
	wg.Add(1)

	err = Wallet.SyncBlockChain(&core.BlockChainSyncListener{
		SyncStarted: func() {
			fmt.Println("Starting blockchain sync")
		},
		SyncEnded: func(e error) {
			err = e
			wg.Done()
		},
		OnHeadersFetched: func(percentageProgress int64) {}, // in cli mode, sync updates are logged to terminal, no need to act on this update alert
		OnDiscoveredAddress: func(state string) {}, // in cli mode, sync updates are logged to terminal, no need to act on update alert
		OnRescanningBlocks: func(percentageProgress int64) {}, // in cli mode, sync updates are logged to terminal, no need to act on update alert
	})

	if err != nil {
		// sync go routine failed to start, nothing to wait for
		wg.Done()
	} else {
		// sync in progress, wait for BlockChainSyncListener.OnComplete
		wg.Wait()
	}
}

func printErrorAndExit(message string, err error) {
	fmt.Fprintln(os.Stderr, message)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	os.Exit(1)
}
