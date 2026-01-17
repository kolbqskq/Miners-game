package game

import (
	"context"
	"encoding/json"
	"fmt"
	"miners_game/internal/game/domain"
	"miners_game/internal/game/equipments"
	"miners_game/internal/game/upgrades"
	"miners_game/internal/miners"
	"miners_game/pkg/errs"

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

func (r *Repository) Save(gameState *domain.GameState) error {
	minersJSON, err := json.Marshal(gameState.Miners)
	if err != nil {
		r.Logger.Error().Err(err).Msg("error Marshal minersJSON")
		return fmt.Errorf("Ошибка преобразования minersJSON: %w", err)
	}
	equipmentsJSON, err := json.Marshal(gameState.Equipments)
	if err != nil {
		r.Logger.Error().Msg("error Marshal equipmentsJSON")
		return fmt.Errorf("Ошибка преобразования equipmentsJSON: %w", err)
	}
	upgradesJSON, err := json.Marshal(gameState.Upgrades)
	if err != nil {
		r.Logger.Error().Msg("error Marshal upgradesJSON")
		return fmt.Errorf("Ошибка преобразования upgradesJSON: %w", err)
	}
	query := `
			INSERT INTO games (user_id, game_id, balance, income, last_update_at, miners, equipments, upgrades)
			VAlUES (@user_id, @game_id, @balance, @income, @last_update_at, @miners, @equipments, @upgrades)
			ON CONFLICT (user_id, game_id) DO UPDATE SET balance = EXCLUDED.balance, income = EXCLUDED.income, last_update_at = EXCLUDED.last_update_at, miners = EXCLUDED.miners, equipments = EXCLUDED.equipments, upgrades = EXCLUDED.upgrades`
	args := pgx.NamedArgs{
		"user_id":        gameState.UserID,
		"game_id":        gameState.GameID,
		"balance":        gameState.Balance,
		"income":         gameState.IncomePerSec,
		"last_update_at": gameState.LastUpdateAt,
		"miners":         minersJSON,
		"equipments":     equipmentsJSON,
		"upgrades":       upgradesJSON,
	}
	if _, err = r.DbPool.Exec(context.Background(), query, args); err != nil {
		r.Logger.Error().Err(err).Msg("Ошибка сохранения игры в БД")
		return fmt.Errorf("Невозможно сохронить игру: %w", err)
	}
	return nil
}

func (r *Repository) Load(userID, gameID string) (*domain.GameState, error) {
	query := `
		SELECT balance, income, last_update_at, miners, equipments, upgrades
		FROM games
		WHERE user_id = @user_id AND game_id = @game_id
	`
	rows := r.DbPool.QueryRow(context.Background(), query, pgx.NamedArgs{
		"user_id": userID,
		"game_id": gameID,
	})
	var balance int64
	var income int64
	var lastUpdateAt int64
	var minersJSON []byte
	var equipmentsJSON []byte
	var upgradesJSON []byte

	if err := rows.Scan(&balance, &income, &lastUpdateAt, &minersJSON, &equipmentsJSON, &upgradesJSON); err != nil {
		if err == pgx.ErrNoRows {
			r.Logger.Error().Msg("Запрос на не существующее сохранение в БД")
			return nil, errs.ErrGameNotFound
		}
		r.Logger.Error().Err(err).Msg("Ошибка загрузки сохранения")
		return nil, err
	}

	var miners map[string]*miners.Miner
	var equipments []equipments.Equipment
	var upgrades []upgrades.Upgrade
	if err := json.Unmarshal(minersJSON, &miners); err != nil {
		r.Logger.Error().Err(err).Msg("error Unmarshal GetGameState")
		return nil, err
	}
	if err := json.Unmarshal(equipmentsJSON, &equipments); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(upgradesJSON, &upgrades); err != nil {
		return nil, err
	}

	gs := &domain.GameState{
		UserID:       userID,
		GameID:       gameID,
		Balance:      balance,
		IncomePerSec: income,
		LastUpdateAt: lastUpdateAt,
		Miners:       miners,
		Equipments:   equipments,
		Upgrades:     upgrades,
	}

	return gs, nil
}
