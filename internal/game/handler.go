package game

import (
	"miners_game/internal/miners"
	"miners_game/pkg/middleware"
	"miners_game/pkg/tadapter"
	"miners_game/views/widgets"
	"strconv"
	"time"

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
	g.Post("/save", h.save)
	g.Get("/hud", h.hud)
	g.Get("/new", h.newGame)
	g.Post("buy", h.buy)
	g.Get("/panel/:tab", h.shopPanel)

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

func (h *Handler) hud(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	saveID, ok := c.Locals("save_id").(string)
	if !ok {
		c.Set("HX-Redirect", "game/new")
		return c.SendStatus(fiber.StatusNoContent)
	}
	gameState := h.gameService.GetGameStateFromMemory(userID, saveID)
	if gameState == nil {
		c.Set("HX-Redirect", "/")
		return c.SendStatus(fiber.StatusNoContent)
	}
	incomeInt := gameState.RecalculateBalance()
	balance := strconv.Itoa(int(gameState.Balance))
	income := strconv.Itoa(int(incomeInt))
	component := widgets.HUD(balance, income)
	return tadapter.Render(c, component, fiber.StatusOK)
}

func (h *Handler) newGame(c *fiber.Ctx) error {
	h.logger.Info().Msg("new game")
	userID := c.Locals("user_id").(string)
	sess := c.Locals("sess").(*session.Session)

	saveID := uuid.NewString()
	h.gameService.StartNewGame(userID, saveID)
	sess.Set("save_id", saveID)
	if err := sess.Save(); err != nil {
		h.logger.Error().Msg("Ошибка сохраниния сессии")
		return c.SendStatus(500)
	}
	c.Set("HX-Redirect", "/game")
	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) buy(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	saveID, ok := c.Locals("save_id").(string)
	if !ok {
		c.Set("HX-Redirect", "game/new")
		c.SendStatus(fiber.StatusNoContent)
	}
	class := c.FormValue("class")
	gameState, incomeInt, err := h.gameService.BuyMiner(class, userID, saveID)
	if err != nil {
		h.logger.Error().Msg(err.Error())
		return c.SendStatus(500)
	}
	balance := strconv.Itoa(int(gameState.Balance))
	income := strconv.Itoa(int(incomeInt))
	component := widgets.HUD(balance, income)
	return tadapter.Render(c, component, fiber.StatusOK)
}

func (h *Handler) shopPanel(c *fiber.Ctx) error {
	tab := c.Params("tab")
	component:=widgets.BottomPanel(tab)
	return tadapter.Render(c, component, fiber.StatusOK)
}
