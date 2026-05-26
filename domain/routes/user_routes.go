package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/username/msa-boilerplate-go/config"
	"github.com/username/msa-boilerplate-go/internal/user/handler"
	"github.com/username/msa-boilerplate-go/internal/user/repository"
	"github.com/username/msa-boilerplate-go/internal/user/route"
	"github.com/username/msa-boilerplate-go/internal/user/service"
	"gorm.io/gorm"
)

func RegisterUserRoutes(api fiber.Router, cfg *config.Config, db *gorm.DB) {
	usrHandler := handler.NewUserHandler(service.NewUserService(repository.NewUserRepository(db), cfg))
	route.RegisterRoute(api, usrHandler, cfg)
}
