package main

// @title           Person API
// @version         1.0
// @description     Person managing API

import (
	"log"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ivanjabrony/personApi/cmd/app"
	"github.com/ivanjabrony/personApi/cmd/config"
	"github.com/ivanjabrony/personApi/cmd/initDB"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.New()

	db, err := initDB.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err := initDB.RunMigrations(db, cfg.Database.Name, "file://migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	application := app.New(db)
	if err := application.Run(cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	if err := initDB.DownMigrations(db, cfg.Database.Name, "file://migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
}
