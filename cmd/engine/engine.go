package engine

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/username/msa-boilerplate-go/config"
	"github.com/username/msa-boilerplate-go/domain/routes"
	"gorm.io/gorm"

	"github.com/username/msa-boilerplate-go/internal/cronjob"
	exampleHandler "github.com/username/msa-boilerplate-go/internal/example-module/handler"
	exampleRepo "github.com/username/msa-boilerplate-go/internal/example-module/repository"
	exampleRoute "github.com/username/msa-boilerplate-go/internal/example-module/route"
	exampleSvc "github.com/username/msa-boilerplate-go/internal/example-module/service"
)

func SetupApp(cfg *config.Config, db *gorm.DB) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "Boilerplate App",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())

	// Setup Cron Jobs
	cronJob := cronjob.SetupCronJobs(db)
	cronJob.Start()

	exHandler := exampleHandler.NewExampleHandler(exampleSvc.NewExampleService(exampleRepo.NewExampleRepository(db)))

	// Routing
	api := app.Group("/api")

	exampleRoute.RegisterRoute(api, exHandler)

	// Register Routes
	routes.RegisterUserRoutes(api, cfg, db)
	return app
}
