package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func LoggerContextMiddleware(logger *zerolog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		reqID := c.Locals("requestid").(string)
		reqLogger := logger.With().Str("request_id", reqID).Logger()
		c.Locals("logger", reqLogger)
		c.Set("X-Request-ID", reqID)
		return c.Next()
	}
}
