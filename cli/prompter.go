package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// prompter fetches input from a prompt, returning the value received, and an error, if any.
type prompter interface {
	Prompt() (int, string, error)
}

type validatorFunc func(value string) error

type promptOption struct {
	Label    interface{}
	Options  []interface{}
	Default  string
	Validate validatorFunc
	Secure   bool
}

func (p promptOption) Prompt() (string, error) {
	label := ""
	if p.Default != "" {
		label += fmt.Sprintf("[default: %s]", p.Default)
	}
	label += "=> "
	fmt.Println(p.Label)
	for _, opt := range p.Options {
		fmt.Println(opt)
	}
	input, err := readInput(label, p.Secure)
	if err != nil {
		return "", err
	}
	if p.Validate != nil {
		if vErr := p.Validate(input); vErr != nil {
			return "", vErr
		}
	}
	return strings.TrimSpace(input), nil
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
