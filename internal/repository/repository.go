package repository

import "errors"

var (
	ErrUsernameTaken = errors.New("username is already taken")
)
