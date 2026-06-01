package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"cms-builder/api/internal/config"
	"cms-builder/api/internal/database"
	"cms-builder/api/internal/handlers"
	"cms-builder/api/internal/services"
)

func main() {
	cfg := config.Load()
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database open failed: %v", err)
	}

	startupCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := waitForDatabase(startupCtx, db); err != nil {
		log.Fatalf("database unavailable: %v", err)
	}

	migrationsCtx, migrationsCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer migrationsCancel()

	if err := database.RunMigrations(migrationsCtx, db, "migrations"); err != nil {
		log.Fatalf("database migrations failed: %v", err)
	}

	svc := services.New(db, cfg)
	mux := handlers.NewRouter(svc, cfg)

	addr := ":" + cfg.APIPort
	log.Printf("api listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func waitForDatabase(ctx context.Context, db *sql.DB) error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		if err := db.PingContext(ctx); err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
