package game

import (
	"context"
	"encoding/json"
	"fmt"
	"miners_game/internal/miners"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type GameRepository struct {
	DbPool *pgxpool.Pool
	Logger *zerolog.Logger
}

type GameRepositoryDeps struct {
	DbPool *pgxpool.Pool
	Logger *zerolog.Logger
}

func NewRepository(deps GameRepositoryDeps) *GameRepository {
	return &GameRepository{
		DbPool: deps.DbPool,
		Logger: deps.Logger,
	}
}

func (r *GameRepository) SaveGameState(gameState *GameState) error {
	minersJSON, err := json.Marshal(gameState.Miners)
	if err != nil {
		r.Logger.Error().Err(err).Msg("error Marshal SaveGameState")
		return fmt.Errorf("Ошибка преобразования MinersJSON: %w", err)
	}
	query := `
			INSERT INTO game_saves (user_id, save_id, balance, last_update_at, miners)
			VAlUES (@user_id, @save_id, @balance, @last_update_at, @miners)
			ON CONFLICT (user_id, save_id) DO UPDATE SET balance = EXCLUDED.balance, last_update_at = EXCLUDED.last_update_at, miner = EXCLUDED.miners`
	args := pgx.NamedArgs{
		"user_id":        gameState.UserID,
		"save_id":        gameState.SaveID,
		"balance":        gameState.Balance,
		"last_update_at": gameState.LastUpdateAt,
		"miners":         minersJSON,
	}
	_, err = r.DbPool.Exec(context.Background(), query, args)
	if err != nil {
		r.Logger.Error().Err(err).Msg("Ошибка сохранения игры в БД")
		return fmt.Errorf("Невозможно сохронить игру: %w", err)
	}
	return nil
}

func (r *GameRepository) GetGameState(userID, saveID string) (*GameState, error) {
	query := `
		SELECT balance, last_update_at, miners
		FROM game_saves
		WHERE user_id = @user_id AND save_id = @save_id
	`
	rows := r.DbPool.QueryRow(context.Background(), query, pgx.NamedArgs{
		"user_id": userID,
		"save_id": saveID,
	})
	var balance int64
	var lastUpdateAt int64
	var minersJSON []byte

	if err := rows.Scan(&balance, &lastUpdateAt, &minersJSON); err != nil {
		if err == pgx.ErrNoRows {
			r.Logger.Error().Msg("Запрос на не существующее сохранение в БД")
			return nil, fmt.Errorf("сохранение не найдено")
		}
		r.Logger.Error().Err(err).Msg("Ошибка загрузки сохранения")
		return nil, err
	}

	var miners []miners.Miner
	if err := json.Unmarshal(minersJSON, &miners); err != nil {
		r.Logger.Error().Err(err).Msg("error Unmarshal GetGameState")
		return nil, err
	}

	gs := &GameState{
		UserID:       userID,
		SaveID:       saveID,
		Balance:      balance,
		LastUpdateAt: lastUpdateAt,
		Miners:       miners,
	}

	return gs, nil
}
