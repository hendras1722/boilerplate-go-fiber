package route

import (
	"github.com/gofiber/fiber/v3"
	"github.com/username/project-name/config"
	"github.com/username/project-name/internal/middleware"
	"github.com/username/project-name/internal/user/handler"
)

func RegisterRoute(router fiber.Router, h handler.UserHandler, cfg *config.Config) {
	group := router.Group("/user")
	
	group.Post("/register", h.Register)
	group.Post("/login", h.Login)
	group.Post("/refresh-token", h.RefreshToken)

	// Protected routes
	group.Use(middleware.AuthMiddleware(cfg))
	group.Get("/", h.List)
	group.Get("/:id", h.Detail)
}
