package cli

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/jessevdk/go-flags"
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

// displayAvailableCommandsHelpMessage prints a simple list of available commands when dcrcli is run without any command
func displayAvailableCommandsHelpMessage(parser *flags.Parser) {
	registeredCommands := parser.Commands()
	commandNames := make([]string, 0, len(registeredCommands))
	for _, command := range registeredCommands {
		commandNames = append(commandNames, command.Name)
	}
	sort.Strings(commandNames)
	fmt.Fprintln(os.Stderr, "Available Commands: ", strings.Join(commandNames, ", "))
}

func printErrorAndExit(message string, err error) {
	fmt.Fprintln(os.Stderr, message)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	os.Exit(1)
}
