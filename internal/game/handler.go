package game

import (
	"miners_game/internal/game/shop"
	"miners_game/pkg/middleware"
	"miners_game/pkg/tadapter"
	"miners_game/views"
	"miners_game/views/components"
	"miners_game/views/widgets"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Handler struct {
	router      fiber.Router
	logger      *zerolog.Logger
	gameService *Service
	store       *session.Store
}

type HandlerDeps struct {
	Router      fiber.Router
	Logger      *zerolog.Logger
	GameService *Service
	Store       *session.Store
}

func NewHandler(deps HandlerDeps) {
	h := &Handler{
		router:      deps.Router,
		logger:      deps.Logger,
		gameService: deps.GameService,
		store:       deps.Store,
	}
	g := h.router.Group("/game")
	g.Use(middleware.GameMiddleware(h.store))
	g.Get("/", h.game)
	g.Get("/hud", h.hud)
	g.Get("/new", h.newGame)
	g.Post("/buy", h.buy)
	g.Get("/panel/:tab", h.shopTab)
}

func (h *Handler) game(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	gameID := c.Locals("game_id").(string)
	sess := c.Locals("sess").(*session.Session)
	if gameID == "" {
		gameID = uuid.NewString()
	}
	if _, err := h.gameService.enterGame(userID, gameID); err != nil {
		h.logger.Error().Msg("EnterGame failed")
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	sess.Set("game_id", gameID)
	if err := sess.Save(); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	component := views.Game()
	return tadapter.Render(c, component, fiber.StatusOK)
}

func (h *Handler) hud(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	gameID, ok := c.Locals("game_id").(string)
	if !ok || gameID == "" {
		c.Set("HX-Redirect", "/")
		return c.SendStatus(fiber.StatusNoContent)
	}
	balance, income, err := h.gameService.getHud(userID, gameID)
	if err != nil {
		h.logger.Warn().Msg("Hud bad request")
		return c.SendStatus(fiber.StatusNoContent)
	}
	component := widgets.HUD(balance, income)
	return tadapter.Render(c, component, fiber.StatusOK)
}

func (h *Handler) newGame(c *fiber.Ctx) error {
	h.logger.Info().Msg("new game")
	userID := c.Locals("user_id").(string)
	sess := c.Locals("sess").(*session.Session)

	gameID := uuid.NewString()

	_, err := h.gameService.enterGame(userID, gameID)
	if err != nil {
		return c.SendStatus(fiber.StatusNoContent)
	}

	sess.Set("game_id", gameID)
	if err := sess.Save(); err != nil {
		h.logger.Error().Msg("Ошибка сохраниния сессии")
		return c.SendStatus(500)
	}
	c.Set("HX-Redirect", "/game")
	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) buy(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	gameID, ok := c.Locals("game_id").(string)
	if !ok {
		return c.SendStatus(fiber.StatusNoContent)
	}
	name := c.FormValue("name")
	kind := c.FormValue("kind")

	cases := map[string]func(string, string, string, string) (shop.ShopCard, error){
		"miner":     h.gameService.buyMiner,
		"equipment": h.gameService.buyEquipment,
		"upgrade":   h.gameService.buyUpgrade,
	}
	if cs, ok := cases[kind]; ok {
		if card, err := cs(userID, gameID, name, kind); err != nil {
			component := components.ShopCard(card)
			h.logger.Error().Msg("buy error")
			return tadapter.Render(c, component, fiber.StatusOK)
		}
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) shopTab(c *fiber.Ctx) error {
	tab := c.Params("tab")
	cards := h.gameService.getShopState(tab)
	component := widgets.BottomPanel(tab, cards)
	return tadapter.Render(c, component, fiber.StatusOK)
}
