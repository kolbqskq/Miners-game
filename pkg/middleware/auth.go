package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func AuthMiddleware(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			return c.SendStatus(500)
		}
		userID := ""
		userName := ""
		id, ok := sess.Get("user_id").(string)
		if ok {
			userID = id
			name, _ := sess.Get("username").(string)
			userName = name
		}
		c.Locals("user_id", userID)
		c.Locals("username", userName)
		return c.Next()
	}

}
