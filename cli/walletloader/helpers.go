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
	createWalletPrompt := "No wallet found. Would you like to create one now? (y/N)"
	validateUserResponse := func(userResponse string) error {
		userResponse = strings.TrimSpace(userResponse)
		userResponse = strings.Trim(userResponse, `"`)
		if userResponse == "" || strings.EqualFold("y", userResponse) || strings.EqualFold("N", userResponse) {
			return nil
		} else {
			return fmt.Errorf("invalid option, try again")
		}
	}
	userResponse, err := terminalprompt.RequestInput(createWalletPrompt, validateUserResponse)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading your response: %s", err.Error())
		return err
	}

	if userResponse == "" || strings.EqualFold("N", userResponse) {
		fmt.Println("Maybe later. Bye.")
		return nil
	}

	return CreateWallet(ctx, walletMiddleware)
}
