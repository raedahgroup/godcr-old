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

	// Catch a ^C interrupt.
	// Make sure that we reset term echo before exiting.
	//signalChannel := make(chan os.Signal, 1)
	//signal.Notify(signalChannel, os.Interrupt)
	//go func() {
	//	<-signalChannel
	//	setTerminalEcho(true)
	//	os.Exit(1)
	//}()

	// disable terminal echo
	//setTerminalEcho(false)

	// Echo is disabled, now grab the data.
	//reader := bufio.NewReader(os.Stdin)
	//text, err := reader.ReadString('\n')

	// re-enable terminal echo
	//setTerminalEcho(true)

	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	password := string(bytePassword)

	if err != nil {
		return "", err
	}
	return strings.TrimRight(password, "\n"), nil
}

//// techEcho() - turns terminal echo on or off.
//func setTerminalEcho(on bool) {
//	// Common settings and variables for both stty calls.
//	attrs := syscall.ProcAttr{
//		Dir:   "",
//		Env:   []string{},
//		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
//		Sys:   nil}
//	var ws syscall.WaitStatus
//	cmd := "echo"
//	if on == false {
//		cmd = "-echo"
//	}
//
//	// Enable/disable echoing.
//	pid, err := syscall.ForkExec(
//		"/bin/stty",
//		[]string{"stty", cmd},
//		&attrs)
//	if err != nil {
//		panic(err)
//	}
//
//	// Wait for the stty process to complete.
//	_, err = syscall.Wait4(pid, &ws, 0, nil)
//	if err != nil {
//		panic(err)
//	}
//}
