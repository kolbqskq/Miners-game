package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func GameMiddleware(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := c.Locals("logger").(zerolog.Logger)
		sess, err := store.Get(c)
		if err != nil {
			logger.Error().Err(err).Msg("failed to get session")
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		userID, ok := sess.Get("user_id").(string)
		if !ok || userID == "" {
			logger.Warn().Msg("failed to find user_id")
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		gameID, ok := sess.Get("game_id").(string)
		if !ok {
			gameID = uuid.NewString()
			sess.Set("game_id", gameID)
			if err := sess.Save(); err != nil {
				logger.Error().Err(err).Msg("failed to save session")
				return c.SendStatus(fiber.StatusInternalServerError)
			}
		}

		gameLogger := logger.With().Str("game_id", gameID).Logger()

		c.Locals("logger", gameLogger)
		c.Locals("game_id", gameID)
		c.Locals("sess", sess)
		c.Locals("user_id", userID)

		return c.Next()
	}

}
