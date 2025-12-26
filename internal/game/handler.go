package game

import (
	"miners_game/internal/miners"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Handler struct {
	router      fiber.Router
	logger      *zerolog.Logger
	gameService *Service
}

type HandlerDeps struct {
	Router      fiber.Router
	Logger      *zerolog.Logger
	GameService *Service
}

func NewHandler(deps HandlerDeps) {
	h := &Handler{
		router:      deps.Router,
		logger:      deps.Logger,
		gameService: deps.GameService,
	}
	h.router.Post("/save", h.save)
	h.router.Post("/balance", h.balance)
	h.router.Post("/new", h.newGame)
}

func (h *Handler) save(c *fiber.Ctx) error {
	h.logger.Error().Msg("Save")
	if err := h.gameService.repo.SaveGameState(&GameState{
		UserID:       uuid.NewString(),
		SaveID:       uuid.NewString(),
		Balance:      100,
		LastUpdateAt: time.Now().Unix(),
		Miners:       map[string]*miners.Miner{},
	}); err != nil {
		h.logger.Error().Msg("unluck")
		return c.SendStatus(fiber.StatusBadRequest)
	}
	h.logger.Info().Msg("good")
	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) balance(c *fiber.Ctx) error {

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) newGame(c *fiber.Ctx) error {

	return c.SendStatus(fiber.StatusOK)
}
