package authservice

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmptyPassword = errors.New("password is empty")
)

type AuthService struct {
	repo Repository
}

type Repository interface {
	SaveUser(username, passwordHash string) (int, error)
}

func New(repo Repository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (a *AuthService) Register(username, password string) (int, error) {
	passwordHash, err := hashPassword(password)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	userID, err := a.repo.SaveUser(username, passwordHash)
	if err != nil {
		return 0, fmt.Errorf("failed to save user: %w", err)
	}

	return userID, nil
}

func hashPassword(password string) (string, error) {
	if password == "" {
		return "", ErrEmptyPassword
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}
