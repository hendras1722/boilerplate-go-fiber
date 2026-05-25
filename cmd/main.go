package main

import (
	"log"

	"github.com/username/project-name/cmd/engine"
	"github.com/username/project-name/config"
	"github.com/username/project-name/database"
)

func main() {
	// 1. Load config
	cfg := config.LoadConfig()

	// 2. Setup database
	db := database.ConnectDB(cfg)

	// 3. Setup engine (Fiber app & DI)
	app := engine.SetupApp(cfg, db)

	// 4. Run application
	log.Printf("Starting server on port %s...", cfg.AppPort)
	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
