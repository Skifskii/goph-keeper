package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Skifskii/goph-keeper/internal/domain/secret"
	"github.com/Skifskii/goph-keeper/internal/repository"
	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	ErrEmptyDSN = errors.New("DSN is empty")
)

type Postgres struct {
	db *sql.DB
}

func New(log *slog.Logger, dsn string) (*Postgres, error) {
	if dsn == "" {
		return nil, ErrEmptyDSN
	}

	// run migrations
	if err := runMigrations(log, dsn); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}

	// connect to DB
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return &Postgres{db: db}, nil
}

func runMigrations(log *slog.Logger, dsn string) error {
	m, err := migrate.New(
		"file://./migrations",
		dsn,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Debug("migrations applied successfully!")
	return nil
}

func (p *Postgres) SaveSecret(enc secret.EncryptedSecret) (secretID int, err error) {
	err = p.db.QueryRow(
		`INSERT INTO secrets (user_id, encrypted_secret, secret_type, metadata)
		VALUES ($1, $2, $3, $4)
		RETURNING id;`,
		enc.UserID, enc.EncPayload, enc.SecretType, enc.Metadata,
	).Scan(&secretID)
	if err != nil {
		return 0, fmt.Errorf("failed to run query: %w", err)
	}

	return secretID, nil
}

func (p *Postgres) GetSecret(secretID int) (enc secret.EncryptedSecret, err error) {
	row := p.db.QueryRow(
		`SELECT
			user_id,
			encrypted_secret,
			secret_type,
			metadata
		FROM secrets
		WHERE id = $1
		LIMIT 1;`,
		secretID,
	)

	err = row.Scan(&enc.UserID, &enc.EncPayload, &enc.SecretType, &enc.Metadata)
	if err != nil {
		return secret.EncryptedSecret{}, fmt.Errorf("failed to scan row: %w", err)
	}
	return enc, nil
}

func (p *Postgres) SaveUser(username, passwordHash string) (userID int, err error) {
	err = p.db.QueryRow(
		`INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING id`,
		username, passwordHash,
	).Scan(&userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return 0, repository.ErrUsernameTaken
			}
		}
		return 0, fmt.Errorf("failed to run query: %w", err)
	}

	return userID, nil
}
