package authservice

import (
	"errors"
	"fmt"
	"time"

	"github.com/Skifskii/goph-keeper/internal/domain/user"
	"github.com/Skifskii/goph-keeper/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmptyPassword      = errors.New("password is empty")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService struct {
	repo      Repository
	secretKey string
	tokenTTL  time.Duration
}

type Repository interface {
	SaveUser(username, passwordHash string) (int, error)
	GetUser(username string) (user.User, error)
}

func New(repo Repository, secretKey string, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		repo:      repo,
		secretKey: secretKey,
		tokenTTL:  tokenTTL,
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

func (a *AuthService) Login(username, password string) (string, error) {
	// get user from repo
	user, err := a.repo.GetUser(username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", fmt.Errorf("failed to get user from repo: %w", err)
	}

	// check password
	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	// build jwt token
	token, err := a.buildJWTToken(user.ID, user.Username, a.tokenTTL)
	if err != nil {
		return "", fmt.Errorf("failed to build jwt token: %w", err)
	}

	return token, nil
}

func (a *AuthService) buildJWTToken(userID int, username string, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = userID
	claims["username"] = username
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(a.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign string: %w", err)
	}

	return tokenString, nil
}
