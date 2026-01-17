package pages

import (
	"miners_game/pkg/tadapter"
	"miners_game/views"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/rs/zerolog"
)

type Handler struct {
	router fiber.Router
	logger *zerolog.Logger
	store  *session.Store
}

type HandlerDeps struct {
	Router fiber.Router
	Logger *zerolog.Logger
	Store  *session.Store
}

func NewHandler(deps HandlerDeps) {
	h := &Handler{
		router: deps.Router,
		logger: deps.Logger,
		store:  deps.Store,
	}
	h.router.Get("/", h.home)
	h.router.Get("/login", h.login)
	h.router.Get("/register", h.register)
}

func (h *Handler) home(c *fiber.Ctx) error {
	component := views.Main()
	return tadapter.Render(c, component, fiber.StatusOK)
}

func (h *Handler) login(c *fiber.Ctx) error {
	component := views.Login()
	return tadapter.Render(c, component, fiber.StatusOK)
}

func (h *Handler) register(c *fiber.Ctx) error {
	component := views.Register()
	return tadapter.Render(c, component, fiber.StatusOK)
}
