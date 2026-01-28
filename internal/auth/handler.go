package auth

import (
	"miners_game/pkg/code"
	"miners_game/pkg/errs"
	"miners_game/pkg/tadapter"
	"miners_game/views"
	"miners_game/views/components"
	"miners_game/views/layout"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gookit/validate"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
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
		logger.Error().Err(err).Str("handler", "logout").Msg("failed to destroy session")
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
	v := validate.Struct(&form)
	if !v.Validate() {
		logger.Warn().Err(v.Errors).Msg("failed to validate login form")
		component := components.Notification(v.Errors.OneError().Error(), components.NotificationFail)
		return tadapter.Render(c, component, fiber.StatusBadRequest)
	}

	userID, userName, err := h.authService.login(form.Email, form.Password)
	if err != nil {
		component := components.Notification(err.Error(), components.NotificationFail)
		return tadapter.Render(c, component, fiber.StatusBadRequest)
	}

	sess := c.Locals("sess").(*session.Session)

	sess.Set("user_id", userID)
	sess.Set("username", userName)
	if err := sess.Save(); err != nil {
		logger.Error().Err(err).Str("handler", "login").Msg("failed to save session")
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
			component := components.Notification("Сессия истекла", components.NotificationFail)
			return tadapter.Render(c, component, fiber.StatusBadRequest)

		}
		reg := data.(RegisterSession)
		if time.Now().Unix() > reg.ExpiresAt {
			sess.Delete("register")
			if err := sess.Save(); err != nil {
				logger.Error().Err(err).Msg("failed to save session")
				return c.SendStatus(fiber.StatusInternalServerError)
			}
			component := components.Notification("Сессия истекла", components.NotificationFail)
			return tadapter.Render(c, component, fiber.StatusBadRequest)

		}
		if c.FormValue("code") == "" {
			component := components.Notification("Введите код", components.NotificationFail)
			return tadapter.Render(c, component, fiber.StatusBadRequest)
		}
		if c.FormValue("code") != reg.Code {
			component := components.Notification("Неверный код", components.NotificationFail)
			return tadapter.Render(c, component, fiber.StatusBadRequest)

		}
		userID, err := h.authService.register(reg.Email, reg.Username, reg.HashedPassword)
		if err != nil {
			component := components.Notification(err.Error(), components.NotificationFail)
			return tadapter.Render(c, component, fiber.StatusBadRequest)
		}
		sess.Delete("register")
		sess.Set("user_id", userID)
		sess.Set("username", reg.Username)
		if err := sess.Save(); err != nil {
			logger.Error().Err(err).Msg("failed to save session")
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
		v := validate.Struct(&form)
		if !v.Validate() {
			logger.Warn().Err(v.Errors).Msg("failed to validate register form")
			component := components.Notification(v.Errors.OneError().Error(), components.NotificationFail)
			return tadapter.Render(c, component, fiber.StatusBadRequest)
		}
		if h.authService.emailExist(form.Email) {
			logger.Warn().Msg("failed email already exist")
			component := components.Notification(errs.ErrEmailAlreadyExist.Error(), components.NotificationFail)
			return tadapter.Render(c, component, fiber.StatusBadRequest)
		}
		code := code.Generate()
		logger.Debug().Msg("code: " + code)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
		if err != nil {
			logger.Error().Err(err).Msg("failed to hash password")
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		sess.Set("register", RegisterSession{
			Email:          form.Email,
			Code:           code,
			Username:       form.UserName,
			HashedPassword: string(hashedPassword),
			ExpiresAt:      time.Now().Add(10 * time.Minute).Unix(),
		})
		if err := sess.Save(); err != nil {
			logger.Error().Err(err).Msg("failed to save session")
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		go h.authService.sendEmail(form.Email, code)
		return tadapter.Render(c, views.RegisterVerification(), fiber.StatusOK)
	}

}
