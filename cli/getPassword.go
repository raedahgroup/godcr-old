// License: MIT Open Source
// Copyright (c) Joe Linoff 2016
// Go code to prompt for password using only standard packages by utilizing syscall.ForkExec() and syscall.Wait4().
// Correctly resets terminal echo after ^C interrupts.

package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// getPassword - Prompt for password.
func getPassword(prompt string) (string, error) {
	fmt.Println(prompt)

	// Catch a ^C interrupt.
	// Make sure that we reset term echo before exiting.
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	go func() {
		for _ = range signalChannel {
			fmt.Println("\n^C interrupt.")
			termEcho(true)
			os.Exit(1)
		}
	}()

	// disable terminal echo
	termEcho(false)

	// Echo is disabled, now grab the data.
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')

	// re-enable terminal echo
	termEcho(true)
	fmt.Println("")

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(text), nil
}

// techEcho() - turns terminal echo on or off.
func termEcho(on bool) {
	// Common settings and variables for both stty calls.
	attrs := syscall.ProcAttr{
		Dir:   "",
		Env:   []string{},
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
		Sys:   nil}
	var ws syscall.WaitStatus
	cmd := "echo"
	if on == false {
		cmd = "-echo"
	}

	// Enable/disable echoing.
	pid, err := syscall.ForkExec(
		"/bin/stty",
		[]string{"stty", cmd},
		&attrs)
	if err != nil {
		panic(err)
	}

	// Wait for the stty process to complete.
	_, err = syscall.Wait4(pid, &ws, 0, nil)
	if err != nil {
		panic(err)
	}
}
