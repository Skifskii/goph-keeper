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

// ErrEmptyPassword is returned when an empty password is provided.
var ErrEmptyPassword = errors.New("password is empty")

// ErrInvalidCredentials indicates the provided username/password pair
// does not match any known user.
var ErrInvalidCredentials = errors.New("invalid credentials")

// ErrInvalidJWTToken indicates a parsed JWT token was invalid.
var ErrInvalidJWTToken = errors.New("invalid jwt token")

// AuthService provides authentication-related operations such as
// user registration, credential verification and JWT handling.
type AuthService struct {
	repo      Repository
	secretKey string
	tokenTTL  time.Duration
}

// Repository defines the storage operations required by AuthService.
type Repository interface {
	SaveUser(username, passwordHash string) (int, error)
	GetUser(username string) (user.User, error)
}

// Claims represents JWT claims used by the service, extending the
// standard registered claims with a UserID field.
type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

// New constructs a new AuthService with a repository, secret key for
// signing tokens and token TTL duration.
func New(repo Repository, secretKey string, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		repo:      repo,
		secretKey: secretKey,
		tokenTTL:  tokenTTL,
	}
}

// Register creates a new user with the provided credentials and
// returns the created user ID.
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

// Login verifies credentials and returns a signed JWT token on success.
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
	token, err := a.buildJWTToken(user.ID, a.tokenTTL)
	if err != nil {
		return "", fmt.Errorf("failed to build jwt token: %w", err)
	}

	return token, nil
}

func (a *AuthService) buildJWTToken(userID int, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(a.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign string: %w", err)
	}

	return tokenString, nil
}

// AuthorizeWithJWT validates a JWT token string and returns the
// associated user ID if the token is valid.
func (a *AuthService) AuthorizeWithJWT(jwtTokenString string) (int, error) {
	claims := Claims{}
	token, err := jwt.ParseWithClaims(jwtTokenString, &claims,
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(a.secretKey), nil
		})
	if err != nil {
		return 0, fmt.Errorf("failed to parse jwt token: %w", err)
	}

	if !token.Valid {
		return 0, ErrInvalidJWTToken
	}

	return claims.UserID, nil
}
