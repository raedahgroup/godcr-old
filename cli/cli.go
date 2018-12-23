package cli

import (
	"fmt"
	"github.com/raedahgroup/dcrcli/cli/terminalprompt"
	"github.com/raedahgroup/dcrcli/config"
	"os"
	"strings"
	"sync"

	ws "github.com/raedahgroup/dcrcli/walletsource"
)

type Response struct {
	Columns []string
	Result  [][]interface{}
}

var (
	WalletSource ws.WalletSource
	StdoutWriter = tabWriter(os.Stdout)
)

func CreateWallet() {
	// no need to make the user go through stress of providing following info if wallet already exists
	walletExists, err := WalletSource.WalletExists()
	if err != nil {
		errMsg := fmt.Sprintf("Error checking %s wallet", WalletSource.NetType())
		printErrorAndExit(errMsg, err)
	}

	if walletExists {
		netType := strings.Title(WalletSource.NetType())
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
	seed, err := WalletSource.GenerateNewWalletSeed()
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
	err = WalletSource.CreateWallet(passphrase, seed)
	if err != nil {
		printErrorAndExit("Error creating wallet", err)
	}

	fmt.Println("Your wallet has been created successfully.")
}

// called whenever an action to be executed requires wallet to be loaded
// exits the program is wallet doesn't exist or some other error occurs
func OpenWallet() {
	walletExists, err := WalletSource.WalletExists()
	if err != nil {
		errMsg := fmt.Sprintf("Error checking %s wallet", WalletSource.NetType())
		printErrorAndExit(errMsg, err)
	}

	if !walletExists {
		netType := strings.Title(WalletSource.NetType())
		errMsg := fmt.Sprintf("%s wallet does not exist. Use '%s create' to create a wallet", netType, config.AppName())
		printErrorAndExit(errMsg, nil)
	}

	err = WalletSource.OpenWallet()
	if err != nil {
		errMsg := fmt.Sprintf("Failed to open %s wallet", WalletSource.NetType())
		printErrorAndExit(errMsg, err)
	}
}

// syncBlockChain registers a progress listener with walletsource to download block updates
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

	err = WalletSource.SyncBlockChain(&ws.BlockChainSyncListener{
		SyncStarted: func() {
			fmt.Println("Starting blockchain sync")
		},
		SyncEnded: func(e error) {
			err = e
			wg.Done()
		},
		OnHeadersFetched: func(percentageProgress int64) {
			fmt.Printf("1/3 fetching headers %d%% \n", percentageProgress)
		},
		OnDiscoveredAddress: func(state string) {
			fmt.Printf("2/3 %s discovering addresses\n", state)
		},
		OnRescanningBlocks: func(percentageProgress int64) {
			fmt.Printf("3/3 rescanning blocks %d%% \n", percentageProgress)
		},
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
