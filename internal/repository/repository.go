package repository

import "errors"

var (
	ErrUsernameTaken  = errors.New("username is already taken")
	ErrSecretNotFound = errors.New("secret not found")
	ErrUserNotFound   = errors.New("user not found")
)
