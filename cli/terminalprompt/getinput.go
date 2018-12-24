// License: MIT Open Source
// Copyright (c) Joe Linoff 2016
// Go code to prompt for password using only standard packages by utilizing syscall.ForkExec() and syscall.Wait4().
// Correctly resets terminal echo after ^C interrupts.

package terminalprompt

import (
	"bufio"
	"fmt"
	"os"
	//"os/signal"
	"strings"
	"syscall"
	"golang.org/x/crypto/ssh/terminal"
)

// getTextInput - Prompt for text input.
func getTextInput(prompt string) (string, error) {
	fmt.Print(prompt)

	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')

	if err != nil {
		return "", err
	}
	return strings.TrimRight(text, "\n"), nil
}

// getPasswordInput - Prompt for password.
func getPasswordInput(prompt string) (string, error) {
	fmt.Print(prompt)

	reader := bufio.NewReader(os.Stdin)
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		panic(err)
	}
	fmt.Println("\nPassword typed:" + string(bytePassword))
	text, err := reader.ReadString('\n')


	if err != nil {
		return "", err
	}

	return strings.TrimRight(text, "\n"), nil
}

