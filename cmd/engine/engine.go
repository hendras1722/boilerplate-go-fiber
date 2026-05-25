package engine

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/username/project-name/config"
	"gorm.io/gorm"

	exampleHandler "github.com/username/project-name/internal/example-module/handler"
	exampleRepo "github.com/username/project-name/internal/example-module/repository"
	exampleRoute "github.com/username/project-name/internal/example-module/route"
	exampleSvc "github.com/username/project-name/internal/example-module/service"

	userHandler "github.com/username/project-name/internal/user/handler"
	userRepo "github.com/username/project-name/internal/user/repository"
	userRoute "github.com/username/project-name/internal/user/route"
	userSvc "github.com/username/project-name/internal/user/service"
)

func SetupApp(cfg *config.Config, db *gorm.DB) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "Boilerplate App",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())

	// Example Module Dependency Injection
	exHandler := exampleHandler.NewExampleHandler(exampleSvc.NewExampleService(exampleRepo.NewExampleRepository(db)))

	// User Module Dependency Injection
	usrHandler := userHandler.NewUserHandler(userSvc.NewUserService(userRepo.NewUserRepository(db), cfg))

	// Routing
	api := app.Group("/api")

	// Register Routes
	exampleRoute.RegisterRoute(api, exHandler)
	userRoute.RegisterRoute(api, usrHandler, cfg)

	return app
}
