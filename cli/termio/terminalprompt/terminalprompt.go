package terminalprompt

import (
	"fmt"
	"io"
	"strings"
)

// ValidatorFunction  validates the input string according to its custom logic.
type ValidatorFunction func(string) error

// EmptyValidator is a noop validator that can be used if no validation is needed.
var EmptyValidator = func(v string) error {
	return nil
}

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
		value, err := skipEOFError(getTextInput(fmt.Sprintf("%s: ", message)))
		if err != nil {
			return "", err
		}
		if err = validate(value); err != nil {
			fmt.Printf("%s\n", err.Error())
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
		value, err := skipEOFError(getPasswordInput(fmt.Sprintf("%s: ", message)))
		if err != nil {
			return "", err
		}
		if err = validate(value); err != nil {
			fmt.Printf("%s\n", err.Error())
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
	var promptMessage = message + "\n"
	for idx, opt := range options {
		promptMessage += fmt.Sprintf(" [%d]: %s\n", idx+1, opt)
	}
	promptMessage += "=> "
	for {
		value, err := skipEOFError(getTextInput(promptMessage))
		if err != nil {
			return "", err
		}
		if err = validate(value); err != nil {
			fmt.Printf("%s\n", err.Error())
			continue
		}
		return value, nil
	}
}

func RequestYesNoOption(message, defaultOption string) (string, error) {
	yesNoOptions := []string{ "y", "yes", "n", "no" }

	validateUserResponse := func(userResponse string) error {
		userResponse = strings.TrimSpace(userResponse)
		if defaultOption != "" && userResponse == "" {
			return nil
		}

		for _, option := range yesNoOptions {
			if strings.EqualFold(option, userResponse) {
				return nil
			}
		}

		return fmt.Errorf("Invalid option, try again")
	}

	var options string
	if strings.EqualFold(defaultOption, "Y") || strings.EqualFold(defaultOption, "Yes") {
		options = "Y/n"
	} else if strings.EqualFold(defaultOption, "N") || strings.EqualFold(defaultOption, "No") {
		options = "y/N"
	} else {
		options = "y/n"
	}

	// append options to message for display
	message = fmt.Sprintf("%s (%s)", message, options)
	userResponse, err := RequestInput(message, validateUserResponse)
	if err != nil {
		return "", err
	}

	userResponse = strings.TrimSpace(userResponse)
	if defaultOption != "" && userResponse == "" {
		return defaultOption, nil
	}

	return userResponse, nil
}
