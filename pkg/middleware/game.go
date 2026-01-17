package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func GameMiddleware(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		userID, ok := sess.Get("user_id").(string)
		if !ok || userID == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		gameID, ok := sess.Get("game_id").(string)

		c.Locals("game_id", gameID)
		c.Locals("sess", sess)
		c.Locals("user_id", userID)

		return c.Next()
	}

}
