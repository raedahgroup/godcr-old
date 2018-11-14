package terminalprompt

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

type ValidatorFunction func(string) error

func skipEOFError(value string, err error) (string, error) {
	switch err {
	case io.EOF:
		return "", nil
	case nil:
		return value, nil
	default:
		return "", err
	}
}

// RequestInput requests input from the user.
// If an error other than EOF occurs while requesting input, the error is returned.
// It calls `validate` on the received input. If `validate` returns an error, the user is prompted
// again for a correct input.
func RequestInput(message string, validate ValidatorFunction) (string, error) {
	for {
		value, err := skipEOFError(readInput(fmt.Sprintf("%s: ", message), false))
		if err != nil {
			return "", err
		}
		if err = validate(value); err != nil {
			fmt.Println(err.Error())
			continue
		}
		return value, nil
	}
}

// RequestInputSecure requests input from the user, disabling terminal echo.
// If an error other than EOF occurs while requesting input, the error is returned.
// It calls `validate` on the received input. If `validate` returns an error, the user is prompted
// again for a correct input.
func RequestInputSecure(message string, validate ValidatorFunction) (string, error) {
	for {
		value, err := skipEOFError(readInput(fmt.Sprintf("%s: ", message), true))
		if err != nil {
			return "", err
		}
		if err = validate(value); err != nil {
			fmt.Println(err.Error())
			continue
		}
		return value, nil
	}
}

// RequestSelection prompts the user to select from a list of options. The user is expected to enter
// a number that corresponds to an item in the list.
// If an error other than EOF occurs while requesting input, the error is returned.
// It calls `validate` on the received input. If `validate` returns an error, the user is prompted
// again for a correct input.
func RequestSelection(message string, options []string, validate ValidatorFunction) (string, error) {
	var label = message + strings.Repeat("\n", 2)
	for idx, opt := range options {
		label += fmt.Sprintf("%d. %s\n", idx, opt)
	}
	label += "\n=> "
	for {
		value, err := skipEOFError(readInput(label, false))
		if err != nil {
			return "", err
		}
		if err = validate(value); err != nil {
			fmt.Println(err.Error())
			continue
		}
		return value, nil
	}
}

func readInput(prompt string, secure bool) (string, error) {
	// Get the initial state of the terminal.
	initialTermState, e1 := terminal.GetState(syscall.Stdin)
	if e1 != nil {
		return "", e1
	}

	// Restore it in the event of an interrupt.
	// CITATION: Konstantin Shaposhnikov - https://groups.google.com/forum/#!topic/golang-nuts/kTVAbtee9UA
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		<-c
		_ = terminal.Restore(syscall.Stdin, initialTermState)
		os.Exit(1)
	}()

	// Now get the password.
	fmt.Print(prompt)
	var (
		b   []byte
		err error
	)
	if secure {
		b, err = terminal.ReadPassword(syscall.Stdin)
	} else {
		buf := bufio.NewReader(os.Stdin)
		b, _, err = buf.ReadLine()
	}
	fmt.Println("")
	if err != nil {
		return "", err
	}

	// Stop looking for ^C on the channel.
	signal.Stop(c)

	// Return the password as a string.
	return string(b), nil

}
