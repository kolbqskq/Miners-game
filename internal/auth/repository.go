package auth

import (
	"context"

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

func (r *Repository) SaveSession(session *Session) error {
	query := `
		INSERT INTO users (id, user_id, expires_at)
		VALUES (@id, @user_id, @expires_at)
	`
	args := pgx.NamedArgs{
		"id":         session.ID,
		"user_id":    session.UserID,
		"expires_at": session.ExpiresAt,
	}

	if _, err := r.DbPool.Exec(context.Background(), query, args); err != nil {
		r.Logger.Error().Err(err).Msg("Ошибка сохранения пользователя")
		return err
	}

	return nil
}
