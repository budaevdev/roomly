package storage

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

type User struct {
	ID           int64
	Username     string
	PasswordHash string
}

var ErrUsernameTaken = errors.New("username already taken")

func CreateUser(db *sql.DB, u *User) error {
	query := `INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id`

	err := db.QueryRow(query, u.Username, u.PasswordHash).Scan(&u.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrUsernameTaken
		}
		return err
	}

	return nil
}

var ErrUserNotFound = errors.New("user not found")

func GetUserByUsername(db *sql.DB, username string) (*User, error) {
	query := `SELECT id, username, password_hash FROM users WHERE username = $1`

	var u User
	err := db.QueryRow(query, username).Scan(&u.ID, &u.Username, &u.PasswordHash)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}
