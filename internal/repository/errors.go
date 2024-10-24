package repository

import "errors"

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrUnknown      = errors.New("unknown error")
	ErrAppExists    = errors.New("app already exists")
	ErrAppNotFound  = errors.New("app not found")
)
