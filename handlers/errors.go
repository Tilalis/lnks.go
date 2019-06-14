package handlers

import "errors"

// ErrWrongToken when token is wrong
var ErrWrongToken = errors.New("lnks: wrong jwt token")
