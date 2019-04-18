// License: MIT Open Source
// Copyright (c) Joe Linoff 2016
// Go code to prompt for password using only standard packages by utilizing syscall.ForkExec() and syscall.Wait4().
// Correctly resets terminal echo after ^C interrupts.

package terminalprompt

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/raedahgroup/godcr/cli/termio"
	"golang.org/x/crypto/ssh/terminal"
)

type FdReader interface {
	io.Reader
	Fd() uintptr
}

var getch = func(r io.Reader) (byte, error) {
	buf := make([]byte, 1)
	if n, err := r.Read(buf); n == 0 || err != nil {
		if err != nil {
			return 0, err
		}
		return 0, io.EOF
	}
	return buf[0], nil
}

var (
	maxLength            = 512
	ErrInterrupted       = errors.New("interrupted")
	ErrMaxLengthExceeded = fmt.Errorf("maximum byte limit (%v) exceeded", maxLength)
)



type terminalState struct {
	state *terminal.State
}

// getTextInput - Prompt for text input.
func getTextInput(prompt string) (string, error) {
	// printing the prompt with tabWriter to ensure adequate formatting of tabulated list of options
	tabWriter := termio.StdoutWriter
	fmt.Fprint(tabWriter, prompt)
	tabWriter.Flush()

	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')

	if err != nil {
		return "", err
	}

	text = strings.TrimSuffix(text, "\n")
	text = strings.TrimSuffix(text, "\r")

	return text, nil
}

// getPasswordInput - Prompt for password.
func getPasswordInput(prompt string) (string, error) {
	psw, err := getPassword(prompt, true, os.Stdin, os.Stdout)
	if err != nil {
		return "", err
	}
	return string(psw), nil
}

// getPassword returns the input read from terminal.
// If prompt is not empty, it will be output as a prompt to the user
// If masked is true, typing will be matched by asterisks on the screen.
// Otherwise, typing will echo nothing.
func getPassword(prompt string, masked bool, r FdReader, w io.Writer) ([]byte, error) {
	var err error
	var password, backspace, mask []byte
	if masked {
		backspace = []byte("\b \b")
		mask = []byte("*")
	}

	if isTerminal(r.Fd()) {
		if oldState, err := makeRaw(r.Fd()); err != nil {
			return password, err
		} else {
			defer func() {
				restore(r.Fd(), oldState)
				fmt.Fprintln(w)
			}()
		}
	}

	if prompt != "" {
		fmt.Fprint(w, prompt)
	}

	// Track total bytes read, not just bytes in the password.  This ensures any
	// errors that might flood the console with nil or -1 bytes infinitely are
	// capped.
	var counter int
	for counter = 0; counter <= maxLength; counter++ {
		if v, e := getch(r); e != nil {
			err = e
			break
		} else if v == 127 || v == 8 {
			if l := len(password); l > 0 {
				password = password[:l-1]
				fmt.Fprint(w, string(backspace))
			}
		} else if v == 13 || v == 10 {
			break
		} else if v == 3 {
			err = ErrInterrupted
			break
		} else if v != 0 {
			password = append(password, v)
			fmt.Fprint(w, string(mask))
		}
	}

	if counter > maxLength {
		err = ErrMaxLengthExceeded
	}

	return password, err
}

func restore(fd uintptr, oldState *terminalState) error {
	return terminal.Restore(int(fd), oldState.state)
}

func isTerminal(fd uintptr) bool {
	return terminal.IsTerminal(int(fd))
}

func makeRaw(fd uintptr) (*terminalState, error) {
	state, err := terminal.MakeRaw(int(fd))

	return &terminalState{
		state: state,
	}, err
}
