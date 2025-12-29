package auth

import (
	"fmt"
	"math/rand"
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
	sess, err := h.store.Get(c)
	if err != nil {
		h.logger.Error().Msg(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if err := sess.Destroy(); err != nil {
		h.logger.Error().Msg(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	c.Locals("user_id", "")
	c.Locals("username", "")

	component := layout.Menu()
	return tadapter.Render(c, component, fiber.StatusOK)
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

	userID, userName, err := h.authService.Login(form.Email, form.Password)
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
	sess.Set("username", userName)
	if err := sess.Save(); err != nil {
		h.logger.Error().Msg(err.Error()) //err
		component := components.Notification("Server error", components.NotificationFail)
		return tadapter.Render(c, component, fiber.StatusInternalServerError)
	}

	c.Set("HX-Redirect", "/")
	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) register(c *fiber.Ctx) error {
	sess, err := h.store.Get(c)
	if err != nil {
		h.logger.Error().Msg(err.Error()) //err
		return c.SendStatus(500)
	}
	step := c.FormValue("step")
	switch step {

	case "email":
		data := sess.Get("register")
		if data == nil {
			component := components.Notification("Сессия истекла", components.NotificationFail)
			return tadapter.Render(c, component, fiber.StatusBadRequest)

		}
		reg := data.(RegisterSession)
		if time.Now().Unix() > reg.ExpiresAt {
			sess.Delete("register")
			sess.Save()
			if err := sess.Save(); err != nil {
				h.logger.Error().Msg("Ошибка сохраниния сессии")
				return c.SendStatus(500)
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
		userID, err := h.authService.Register(reg.Email, reg.Username, reg.HashedPassword)
		if err != nil {
			h.logger.Error().Msg(err.Error())
			component := components.Notification(err.Error(), components.NotificationFail)
			return tadapter.Render(c, component, fiber.StatusBadRequest)
		}
		sess.Delete("register")
		sess.Set("user_id", userID)
		sess.Set("username", reg.Username)
		sess.Save()
		component2:=components.Notification("Регистрация успешна!", components.NotificationSuccess)
		tadapter.Render(c,component2,fiber.StatusOK)
		c.Set("HX-Redirect", "/")
		return c.SendStatus(fiber.StatusOK)
	default:
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
		code := fmt.Sprintf("%06d", rand.Intn(100000))
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
		if err != nil {
			h.logger.Error().Msg("Ошибка хеширования пароля")
			return c.SendStatus(500)
		}
		sess.Set("register", RegisterSession{
			Email:          form.Email,
			Code:           code,
			Username:       form.UserName,
			HashedPassword: string(hashedPassword),
			ExpiresAt:      time.Now().Add(10 * time.Minute).Unix(),
		})
		if err := sess.Save(); err != nil {
			h.logger.Error().Msg("Ошибка сохраниния сессии")
			return c.SendStatus(500)
		}
		go h.authService.SendEmail(form.Email, code)
		return tadapter.Render(c, views.RegisterVerification(), fiber.StatusOK)
	}

}
