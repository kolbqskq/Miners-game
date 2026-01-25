package game

import (
	"context"
	"encoding/json"
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

func (r *Repository) Save(gameState *domain.GameState) error {
	minersJSON, err := json.Marshal(gameState.Miners)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to marshal miners")
		return errs.ErrServer
	}
	equipmentsJSON, err := json.Marshal(gameState.Equipments)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to marshal equipments")
		return errs.ErrServer
	}
	upgradesJSON, err := json.Marshal(gameState.Upgrades)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to marshal upgrades")
		return errs.ErrServer
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
	if _, err = r.dbPool.Exec(context.Background(), query, args); err != nil {
		r.logger.Error().Err(err).Str("user_id", gameState.UserID).Str("game_id", gameState.GameID).Msg("failed to save game state")
		return errs.ErrServer
	}
	return nil
}

func (r *Repository) Load(userID, gameID string) (*domain.GameState, error) {
	query := `
		SELECT balance, income, last_update_at, miners, equipments, upgrades
		FROM games
		WHERE user_id = @user_id AND game_id = @game_id
	`
	rows := r.dbPool.QueryRow(context.Background(), query, pgx.NamedArgs{
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
			return nil, errs.ErrGameNotFound
		}
		r.logger.Error().Err(err).Str("user_id", userID).Str("game_id", gameID).Msg("failed to load game state")
		return nil, errs.ErrServer
	}

	var miners map[string]*miners.Miner
	var equipments []equipments.Equipment
	var upgrades []upgrades.Upgrade
	if err := json.Unmarshal(minersJSON, &miners); err != nil {
		r.logger.Error().Err(err).Msg("failed to unmarshal miners")
		return nil, errs.ErrServer
	}
	if err := json.Unmarshal(equipmentsJSON, &equipments); err != nil {
		r.logger.Error().Err(err).Msg("failed to unmarshal equipments")
		return nil, errs.ErrServer
	}
	if err := json.Unmarshal(upgradesJSON, &upgrades); err != nil {
		r.logger.Error().Err(err).Msg("failed to unmarshal upgrades")
		return nil, errs.ErrServer
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
