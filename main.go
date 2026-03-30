package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/maakh3/api-petstore-service-layer/handlers"
	"github.com/maakh3/api-petstore-service-layer/repository"
	"github.com/maakh3/api-petstore-service-layer/services"
)

func main() {
	logger := setupLogging()
	db, err := openDB()
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	petRepo := repository.NewPetRepository(db, logger)
	petService := services.NewPetService(petRepo, logger)
	petHandler := handlers.NewPetHandler(petService, logger)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /pet", petHandler.AddPet)
	mux.HandleFunc("PUT /pet", petHandler.UpdatePet)
	mux.HandleFunc("GET /pet/findByStatus", petHandler.FindPetsByStatus)
	mux.HandleFunc("GET /pet/findByTags", petHandler.FindPetsByTags)
	mux.HandleFunc("GET /pet/{petId}", petHandler.GetById)

	logger.Info("api-petstore-service-layer is up and running", "port", 8080)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		logger.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}

func openDB() (*sql.DB, error) {
	host := getenv("DB_HOST", "localhost")
	port := getenv("DB_PORT", "5432")
	user := getenv("DB_USER", "petstore")
	pass := getenv("DB_PASSWORD", "petstore")
	name := getenv("DB_NAME", "petstore")
	ssl := getenv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, pass, name, ssl)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func setupLogging() *slog.Logger {
	logLevel := new(slog.LevelVar)
	logLevel.Set(slog.LevelInfo)

	if strings.EqualFold(os.Getenv("LOG_LEVEL"), "debug") {
		logLevel.Set(slog.LevelDebug)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(logger)

	return logger
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
