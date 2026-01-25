package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/rs/zerolog"
)

func AuthMiddleware(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := c.Locals("logger").(zerolog.Logger)
		sess, err := store.Get(c)
		if err != nil {
			logger.Error().Err(err).Msg("failed to get session")
			return c.SendStatus(500)
		}
		userID := ""
		userName := ""
		id, ok := sess.Get("user_id").(string)
		if ok {
			userID = id
			userName, _ = sess.Get("username").(string)
		}

		userLogger := logger.With().Str("user_id", userID).Logger()

		c.Locals("logger", userLogger)
		c.Locals("user_id", userID)
		c.Locals("username", userName)
		return c.Next()
	}

}
