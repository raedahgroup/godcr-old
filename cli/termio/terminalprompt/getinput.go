// License: MIT Open Source
// Copyright (c) Joe Linoff 2016
// Go code to prompt for password using only standard packages by utilizing syscall.ForkExec() and syscall.Wait4().
// Correctly resets terminal echo after ^C interrupts.

package terminalprompt

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
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

	// Catch a ^C interrupt.
	// Make sure that we reset term echo before exiting.
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	go func() {
		<-signalChannel
		setTerminalEcho(true)
	}()

	// when this function returns (after reading password input), stop listening for interrupt requests on `signalChannel`
	defer signal.Stop(signalChannel)

	// disable terminal echo
	setTerminalEcho(false)

	// Echo is disabled, now grab the data.
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')

	// re-enable terminal echo
	setTerminalEcho(true)

	if err != nil {
		return "", err
	}

	return strings.TrimRight(text, "\n"), nil
}

func setTerminalEcho(on bool) error {
	var sttyEchoArg string
	if on {
		fmt.Println()
		sttyEchoArg = "echo"
	} else {
		sttyEchoArg = "-echo"
	}

	cmd := exec.Command("stty", sttyEchoArg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()

	return cmd.Wait()
}
