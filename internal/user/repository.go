package user

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type Repository struct {
	DbPool *pgxpool.Pool
	Logger *zerolog.Logger
}

type RepositoryDeps struct {
	DbPool *pgxpool.Pool
	Logger *zerolog.Logger
}

func NewRepository(deps RepositoryDeps) *Repository {
	return &Repository{
		DbPool: deps.DbPool,
		Logger: deps.Logger,
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

	if _, err := r.DbPool.Exec(context.Background(), query, args); err != nil {
		r.Logger.Error().Err(err).Msg("Ошибка сохранения пользователя") //log
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
	rows := r.DbPool.QueryRow(context.Background(), query, pgx.NamedArgs{
		"email": email,
	})

	var userID string
	var userName string
	var password string

	if err := rows.Scan(&userID, &userName, &password); err != nil {
		if err == pgx.ErrNoRows {
			r.Logger.Error().Msg("Пользователь не найден в БД") //log
			return nil, fmt.Errorf("Пользователь не найден") //errs
		}
		r.Logger.Error().Err(err).Msg("Ошибка загрузки пользователя") //log
		return nil, err
	}
	return &User{
		ID:       userID,
		Email:    email,
		Password: password,
		UserName: userName,
	}, nil

}
