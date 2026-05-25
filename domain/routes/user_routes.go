package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/username/project-name/config"
	"github.com/username/project-name/internal/user/handler"
	"github.com/username/project-name/internal/user/repository"
	"github.com/username/project-name/internal/user/route"
	"github.com/username/project-name/internal/user/service"
	"gorm.io/gorm"
)

func RegisterUserRoutes(api fiber.Router, cfg *config.Config, db *gorm.DB) {
	usrHandler := handler.NewUserHandler(service.NewUserService(repository.NewUserRepository(db), cfg))
	route.RegisterRoute(api, usrHandler, cfg)
}
