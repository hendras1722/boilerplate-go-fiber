package route

import (
	"github.com/gofiber/fiber/v3"
	"github.com/username/project-name/internal/example-module/handler"
)

func RegisterRoute(router fiber.Router, h handler.ExampleHandler) {
	group := router.Group("/example")
	group.Get("/", h.Example)
}
