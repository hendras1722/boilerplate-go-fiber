package main

import (
	"log"

	"github.com/username/msa-boilerplate-go/cmd/engine"
	"github.com/username/msa-boilerplate-go/config"
	"github.com/username/msa-boilerplate-go/database"
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
