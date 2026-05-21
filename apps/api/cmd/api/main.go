package main

import (
	"log"
	"net/http"
	"os"

	"cms-builder/api/internal/config"
	"cms-builder/api/internal/database"
	"cms-builder/api/internal/handlers"
	"cms-builder/api/internal/services"
)

func main() {
	cfg := config.Load()
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Printf("database unavailable: %v", err)
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

