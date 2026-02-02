package robots

import (
	"github.com/gofiber/fiber/v2"
)

type RobotsHandler struct {
	router fiber.Router
	data   []byte
}

type RobotsHandlerDeps struct {
	Router fiber.Router
	Data   []byte
}

func NewHandler(deps RobotsHandlerDeps) {
	h := RobotsHandler{
		router: deps.Router,
		data:   deps.Data,
	}
	h.router.Get("/robots.txt", h.robots)
}

func (h *RobotsHandler) robots(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/plain")
	return c.Send(h.data)
}
