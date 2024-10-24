package service

import "errors"

var (
	ErrInvalidData        = errors.New("Invalid data")
	ErrInvalidCredentials = errors.New("Invalid credentials")
)
