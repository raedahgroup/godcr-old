package walletloader

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/cli/termio/terminalprompt"
)

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

func attemptToCreateWallet(ctx context.Context, walletMiddleware app.WalletMiddleware) error {
	createWalletPrompt := "No wallet found. Would you like to create one now?"
	createWallet, err := terminalprompt.RequestYesNoConfirmation(createWalletPrompt, "Y")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading your response: %s", err.Error())
		return err
	}

	if !createWallet {
		fmt.Println("Maybe later. Bye.")
		return nil
	}

	return CreateWallet(ctx, walletMiddleware)
}
