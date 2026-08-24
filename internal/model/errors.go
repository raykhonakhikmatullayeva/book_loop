package model

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrInvalid  = errors.New("invalid input")
	ErrForbidden = errors.New("forbidden")
	ErrPasswordLength = errors.New("password can't be less than 8 symbols")
	ErrLoginLength = errors.New("login can't be less than 3 symbols")
)

