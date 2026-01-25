package user

import (
	"context"
	"miners_game/pkg/errs"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type Repository struct {
	dbPool *pgxpool.Pool
	logger zerolog.Logger
}

type RepositoryDeps struct {
	DbPool *pgxpool.Pool
	Logger zerolog.Logger
}

func NewRepository(deps RepositoryDeps) *Repository {
	return &Repository{
		dbPool: deps.DbPool,
		logger: deps.Logger,
	}
}

func (r *Repository) SaveUser(user *User) error {
	query := `
		INSERT INTO users (user_id, email, password, username)
		VALUES (@user_id, @email, @password, @username)
	`
	args := pgx.NamedArgs{
		"user_id":  user.ID,
		"email":    user.Email,
		"password": user.Password,
		"username": user.UserName,
	}

	if _, err := r.dbPool.Exec(context.Background(), query, args); err != nil {
		r.logger.Error().Err(err).Str("user_id", user.ID).Msg("failed to save user")
		return err
	}

	return nil
}

func (r *Repository) FindByEmail(email string) (*User, error) {
	query := `
		SELECT user_id, username, password
		FROM users
		WHERE email = @email
	`
	rows := r.dbPool.QueryRow(context.Background(), query, pgx.NamedArgs{
		"email": email,
	})

	var userID string
	var userName string
	var password string

	if err := rows.Scan(&userID, &userName, &password); err != nil {
		if err == pgx.ErrNoRows {
			return nil, errs.ErrUserNotFound
		}
		r.logger.Error().Err(err).Str("user_id", userID).Msg("failed to find user")
		return nil, errs.ErrServer
	}
	return &User{
		ID:       userID,
		Email:    email,
		Password: password,
		UserName: userName,
	}, nil

}
