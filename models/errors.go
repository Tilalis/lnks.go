package models

import (
	"errors"
)

// ErrNoConnection no connection errors
var ErrNoConnection = errors.New("lnks: you need to connect to database")

// ErrWrongUserPassword wrong password
var ErrWrongUserPassword = errors.New("lnks: wrong user password")

// ErrWrongAlias wron alias
var ErrWrongAlias = errors.New("lnks: wrong alias name")

// ErrWrongURL wrong url
var ErrWrongURL = errors.New("lnks: wrong url")
