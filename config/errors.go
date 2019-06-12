package config

import (
	"fmt"
)

type configError struct {
	message  string
	filename string
}

func (e *configError) Error() string {
	if e.filename == "" {
		return e.message
	}

	return fmt.Sprintf("%s; filename: '%s'", e.message, e.filename)
}

func (e *configError) SetFile(filename string) *configError {
	e.filename = filename
	return e
}

// ErrConfigFileNotFound config file not found
var ErrConfigFileNotFound = &configError{message: "lnks: config file not found"}

// ErrReadingConfigFile error while reading config files
var ErrReadingConfigFile = &configError{message: "lnks: error while reading config"}

// ErrMalformedConfigFile config file is corrupred or malformed
var ErrMalformedConfigFile = &configError{message: "lnks: config file is corrupred or malformed"}
