package auth

import (
	"errors"
	"miners_game/pkg/errs"
	"miners_game/pkg/tadapter"
	"miners_game/views"
	"miners_game/views/components"
	"miners_game/views/layout"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/rs/zerolog"
)

const (
	RegisterStepEmail = "email"
)

type Handler struct {
	router      fiber.Router
	authService *Service
	store       *session.Store
}

type HandlerDeps struct {
	Router      fiber.Router
	AuthService *Service
	Store       *session.Store
}

func NewHandler(deps HandlerDeps) {
	h := &Handler{
		router:      deps.Router,
		authService: deps.AuthService,
		store:       deps.Store,
	}
	auth := h.router.Group("/auth")
	auth.Post("/login", h.login)
	auth.Post("/logout", h.logout)
	auth.Post("/register", h.register)
}

func (h *Handler) logout(c *fiber.Ctx) error {
	logger := c.Locals("logger").(zerolog.Logger)

	sess := c.Locals("sess").(*session.Session)
	if err := sess.Destroy(); err != nil {
		logger.Error().Err(err).Msg("failed to destroy session")
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	c.Locals("user_id", "")
	c.Locals("username", "")
	component := layout.Menu()
	return tadapter.Render(c, component, fiber.StatusOK)
}

func (h *Handler) login(c *fiber.Ctx) error {
	logger := c.Locals("logger").(zerolog.Logger)

	form := LoginForm{
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
	}

	userID, userName, err := h.authService.login(form)
	if err != nil {
		component := components.Notification(err.Error(), components.NotificationFail)
		return tadapter.Render(c, component, fiber.StatusBadRequest)
	}

	sess := c.Locals("sess").(*session.Session)

	sess.Set("user_id", userID)
	sess.Set("username", userName)
	if err := sess.Save(); err != nil {
		logger.Error().Err(err).Msg("failed to save session")
		component := components.Notification("Server error", components.NotificationFail)
		return tadapter.Render(c, component, fiber.StatusInternalServerError)
	}

	c.Set("HX-Redirect", "/")
	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) register(c *fiber.Ctx) error {
	logger := c.Locals("logger").(zerolog.Logger)
	sess := c.Locals("sess").(*session.Session)

	step := c.FormValue("step")

	switch step {

	case RegisterStepEmail:
		data := sess.Get("register")
		if data == nil {
			component := components.Notification(errs.ErrServer.Error(), components.NotificationFail)
			return tadapter.Render(c, component, fiber.StatusInternalServerError)
		}
		regSess := data.(RegisterSession)
		code := c.FormValue("code")
		userID, err := h.authService.completeRegistration(regSess, code)
		if err != nil {
			if errors.Is(err, errs.ErrExpireSession) {
				sess.Delete("register")
				sess.Save()
			}
			component := components.Notification(err.Error(), components.NotificationFail)
			return tadapter.Render(c, component, fiber.StatusBadRequest)
		}
		sess.Delete("register")
		sess.Set("user_id", userID)
		sess.Set("username", regSess.Username)
		if err := sess.Save(); err != nil {
			logger.Error().Err(err).Msg("failed to save session")
			component := components.Notification(errs.ErrServer.Error(), components.NotificationFail)
			return tadapter.Render(c, component, fiber.StatusInternalServerError)
		}
		component2 := components.Notification("Регистрация успешна!", components.NotificationSuccess)

		c.Set("HX-Redirect", "/")
		return tadapter.Render(c, component2, fiber.StatusOK)
	default:
		form := RegisterForm{
			Email:           c.FormValue("email"),
			UserName:        c.FormValue("userName"),
			Password:        c.FormValue("password"),
			PasswordConfirm: c.FormValue("passwordConfirm"),
		}
		regSess, err := h.authService.startRegistration(form)
		if err != nil {
			logger.Warn().Err(err).Msg("failed register step default")
			component := components.Notification(err.Error(), components.NotificationFail)
			return tadapter.Render(c, component, fiber.StatusBadRequest)
		}
		sess.Set("register", regSess)
		if err := sess.Save(); err != nil {
			logger.Error().Err(err).Msg("failed save session")
			component := components.Notification(errs.ErrServer.Error(), components.NotificationFail)
			return tadapter.Render(c, component, fiber.StatusInternalServerError)
		}
		return tadapter.Render(c, views.RegisterVerification(), fiber.StatusOK)
	}

}
