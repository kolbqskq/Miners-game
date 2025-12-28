package auth

import (
	"miners_game/pkg/tadapter"
	"miners_game/views/components"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gookit/validate"
	"github.com/rs/zerolog"
)

type Handler struct {
	router      fiber.Router
	logger      *zerolog.Logger
	authService *Service
	store       *session.Store
}

type HandlerDeps struct {
	Router      fiber.Router
	Logger      *zerolog.Logger
	AuthService *Service
	Store       *session.Store
}

func NewHandler(deps HandlerDeps) {
	h := &Handler{
		router:      deps.Router,
		logger:      deps.Logger,
		authService: deps.AuthService,
		store:       deps.Store,
	}
	auth := h.router.Group("/auth")
	auth.Post("/login", h.login)
	auth.Post("/logout", h.logout)
	auth.Post("/register", h.register)
}

func (h *Handler) logout(c *fiber.Ctx) error {

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) login(c *fiber.Ctx) error {
	form := LoginForm{
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
	}
	v := validate.Struct(&form)
	if !v.Validate() {
		h.logger.Error().Msg(v.Errors.String())
		component := components.Notification(v.Errors.OneError().Error(), components.NotificationFail)
		return tadapter.Render(c, component, fiber.StatusBadRequest)
	}

	userID, err := h.authService.Login(form.Email, form.Password)
	if err != nil {
		h.logger.Error().Msg(err.Error())
		component := components.Notification(err.Error(), components.NotificationFail)
		return tadapter.Render(c, component, fiber.StatusBadRequest)
	}

	sess, err := h.store.Get(c)
	if err != nil {
		h.logger.Error().Msg(err.Error())
		component := components.Notification("Server error", components.NotificationFail)
		return tadapter.Render(c, component, fiber.StatusInternalServerError)
	}
	sess.Set("user_id", userID)
	if err := sess.Save(); err != nil {
		h.logger.Error().Msg(err.Error()) //err
		component := components.Notification("Server error", components.NotificationFail)
		return tadapter.Render(c, component, fiber.StatusInternalServerError)
	}
	c.Set("HX-Redirect", "/")
	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) register(c *fiber.Ctx) error {
	form := RegisterForm{
		Email:           c.FormValue("email"),
		UserName:        c.FormValue("userName"),
		Password:        c.FormValue("password"),
		PasswordConfirm: c.FormValue("passwordConfirm"),
	}
	v := validate.Struct(&form)
	if !v.Validate() {
		h.logger.Error().Msg(v.Errors.String())
		component := components.Notification(v.Errors.OneError().Error(), components.NotificationFail)
		return tadapter.Render(c, component, fiber.StatusBadRequest)
	}
	userID, err := h.authService.Register(form.Email, form.Password, form.UserName)
	if err != nil {
		h.logger.Error().Msg(err.Error())
		component := components.Notification(err.Error(), components.NotificationFail)
		return tadapter.Render(c, component, fiber.StatusBadRequest)
	}
	sess, err := h.store.Get(c)
	if err != nil {
		h.logger.Error().Msg(err.Error()) //err
		component := components.Notification("Server error", components.NotificationFail)
		return tadapter.Render(c, component, fiber.StatusInternalServerError)
	}
	sess.Set("user_id", userID)
	if err := sess.Save(); err != nil {
		h.logger.Error().Msg(err.Error()) //err
		component := components.Notification("Server error", components.NotificationFail)
		return tadapter.Render(c, component, fiber.StatusInternalServerError)
	}
	c.Set("HX-Redirect", "/")
	return c.SendStatus(fiber.StatusOK)
}
