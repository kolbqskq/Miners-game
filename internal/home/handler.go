package home

import (
	"miners_game/pkg/tadapter"
	"miners_game/views"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type Handler struct {
	router fiber.Router
	logger *zerolog.Logger
}

type HandlerDeps struct {
	Router fiber.Router
	Logger *zerolog.Logger
}

func NewHandler(deps HandlerDeps) {
	h := &Handler{
		router: deps.Router,
		logger: deps.Logger,
	}
	h.router.Get("/", h.error)
}

func (h *Handler) error(c *fiber.Ctx) error {
	h.logger.Error().Msg("Error")
	component := views.Main()
	return tadapter.Render(c, component)
}
