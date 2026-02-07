package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Skifskii/goph-keeper/internal/domain/secret"
	"github.com/Skifskii/goph-keeper/internal/domain/user"
	"github.com/Skifskii/goph-keeper/internal/repository"
	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// ErrEmptyDSN is returned when a Postgres constructor is invoked with
// an empty Data Source Name.
var ErrEmptyDSN = errors.New("DSN is empty")

// Postgres is a repository implementation that persists data in a
// PostgreSQL database.
type Postgres struct {
	db *sql.DB
}

// New creates a Postgres repository, runs migrations and opens a DB
// connection using the provided DSN.
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

// Stop closes the underlying database connection.
func (p *Postgres) Stop() error {
	return p.db.Close()
}

// SaveSecret inserts an encrypted secret into the database and returns
// the newly created secret ID.
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

// GetSecret retrieves an encrypted secret by its ID. It returns
// repository.ErrSecretNotFound when no row exists for the given id.
func (p *Postgres) GetSecret(secretID int) (enc secret.EncryptedSecret, err error) {
	row := p.db.QueryRow(
		`SELECT
			id,
			user_id,
			encrypted_secret,
			secret_type,
			metadata
		FROM secrets
		WHERE id = $1
		LIMIT 1;`,
		secretID,
	)

	err = row.Scan(&enc.ID, &enc.UserID, &enc.EncPayload, &enc.SecretType, &enc.Metadata)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return secret.EncryptedSecret{}, repository.ErrSecretNotFound
		}
		return secret.EncryptedSecret{}, fmt.Errorf("failed to scan row: %w", err)
	}
	return enc, nil
}

// SaveUser creates a new user record and returns its database ID.
// It maps database uniqueness violations to repository.ErrUsernameTaken.
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

// GetUser fetches a user by username. It returns repository.ErrUserNotFound
// when the username does not exist.
func (p *Postgres) GetUser(username string) (u user.User, err error) {
	row := p.db.QueryRow(
		`SELECT
			id,
			username,
			password_hash
		FROM users
		WHERE username = $1
		LIMIT 1;`,
		username,
	)

	err = row.Scan(&u.ID, &u.Username, &u.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user.User{}, repository.ErrUserNotFound
		}
		return user.User{}, fmt.Errorf("failed to scan row: %w", err)
	}
	return u, nil
}

// GetUserSecrets returns a paginated list of BaseSecret records for a
// given user.
func (p *Postgres) GetUserSecrets(userID, limit, offset int) ([]secret.BaseSecret, error) {
	rows, err := p.db.Query(
		`SELECT
			id,
			user_id,
			metadata,
			secret_type
		FROM secrets
		WHERE user_id = $1
		ORDER BY id
		LIMIT $2 OFFSET $3;`,
		userID,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query secrets: %w", err)
	}
	defer rows.Close()

	secrets := make([]secret.BaseSecret, 0)

	for rows.Next() {
		var s secret.BaseSecret
		if err := rows.Scan(&s.ID, &s.UserID, &s.Metadata, &s.SecretType); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		secrets = append(secrets, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return secrets, nil
}

// UpdateSecret updates an existing secret's encrypted payload and
// metadata. It returns repository.ErrSecretNotFound when the target
// row does not exist.
func (p *Postgres) UpdateSecret(enc secret.EncryptedSecret) error {
	res, err := p.db.Exec(
		`UPDATE secrets
		SET encrypted_secret = $1,
		    metadata = $2
		WHERE id = $3;`,
		enc.EncPayload,
		enc.Metadata,
		enc.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to run update query: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return repository.ErrSecretNotFound
	}

	return nil
}

// DeleteSecret removes a secret row by id. If no row is deleted, it
// returns sql.ErrNoRows.
func (p *Postgres) DeleteSecret(secretID int) error {
	res, err := p.db.Exec(
		`DELETE FROM secrets
		 WHERE id = $1;`,
		secretID,
	)
	if err != nil {
		return fmt.Errorf("failed to execute delete query: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
