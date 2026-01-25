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
	"github.com/rs/zerolog"
)

type Handler struct {
	router      fiber.Router
	gameService *Service
	store       *session.Store
}

type HandlerDeps struct {
	Router      fiber.Router
	GameService *Service
	Store       *session.Store
}

func NewHandler(deps HandlerDeps) {
	h := &Handler{
		router:      deps.Router,
		gameService: deps.GameService,
		store:       deps.Store,
	}
	g := h.router.Group("/game")
	g.Use(middleware.GameMiddleware(h.store))
	g.Get("/", h.game)
	g.Get("/hud", h.hud)
	g.Post("/buy", h.buy)
	g.Get("/panel/:tab", h.shopTab)
	g.Get("/upgrade", h.refreshUpgrade)
	g.Get("/shop/card/:kind/:name", h.shopCard)
}

func (h *Handler) game(c *fiber.Ctx) error {
	logger := c.Locals("logger").(zerolog.Logger)
	userID := c.Locals("user_id").(string)
	gameID := c.Locals("game_id").(string)

	if _, err := h.gameService.enterGame(userID, gameID); err != nil {
		logger.Error().Err(err).Msg("failed enterGame service")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	component := views.Game()
	return tadapter.Render(c, component, fiber.StatusOK)
}

func (h *Handler) hud(c *fiber.Ctx) error {
	logger := c.Locals("logger").(zerolog.Logger)
	userID := c.Locals("user_id").(string)
	gameID := c.Locals("game_id").(string)
	balance, income, err := h.gameService.getHud(userID, gameID)
	if err != nil {
		logger.Error().Err(err).Msg("failed getHud service")
		return c.SendStatus(fiber.StatusNoContent)
	}
	component := widgets.HUD(balance, income)
	return tadapter.Render(c, component, fiber.StatusOK)
}

func (h *Handler) buy(c *fiber.Ctx) error {
	logger := c.Locals("logger").(zerolog.Logger)
	userID := c.Locals("user_id").(string)
	gameID := c.Locals("game_id").(string)

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
			logger.Warn().Err(err).Msg("failed buy service")
			return tadapter.Render(c, component, fiber.StatusOK)
		}
	}
	if kind == "upgrade" {
		c.Set("HX-Trigger", "refresh-upgrade")
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) shopTab(c *fiber.Ctx) error {
	tab := c.Params("tab")
	cards := h.gameService.getShopState(tab)
	component := widgets.BottomPanel(tab, cards)
	return tadapter.Render(c, component, fiber.StatusOK)
}

func (h *Handler) refreshUpgrade(c *fiber.Ctx) error {
	logger := c.Locals("logger").(zerolog.Logger)
	userID := c.Locals("user_id").(string)
	gameID := c.Locals("game_id").(string)

	currUpgrade, err := h.gameService.getCurrUpgrade(userID, gameID)
	if err != nil {
		logger.Error().Err(err).Msg("failed getCurrUpgrade service")
		return c.SendStatus(fiber.StatusNoContent)
	}
	component := widgets.Scene(currUpgrade)
	return tadapter.Render(c, component, fiber.StatusOK)
}

func (h *Handler) shopCard(c *fiber.Ctx) error {
	kind := c.Params("kind")
	name := c.Params("name")
	card := GetShopCardByName(name, kind)
	component := components.ShopCard(card)
	return tadapter.Render(c, component, fiber.StatusOK)
}
