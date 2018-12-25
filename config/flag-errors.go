package config

import flags "github.com/jessevdk/go-flags"

// IsFlagErrorType determines whether a given error is of a given flags.ErrorType.
// It is safe to call IsFlagErrorType with err = nil.
func IsFlagErrorType(err error, errorType flags.ErrorType) bool {
	if err == nil {
		return false
	}
	if flagErr, ok := err.(*flags.Error); ok && flagErr.Type == errorType {
		return true
	}
	return false
}
