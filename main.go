package main

import (
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

	petRepo := repository.NewPetRepository(logger)
	petService := services.NewPetService(petRepo, logger)
	petHandler := handlers.NewPetHandler(petService, logger)

	storeRepo := repository.NewStoreRepository()
	storeService := services.NewStoreService(storeRepo)
	storeHandler := handlers.NewStoreHandler(storeService)

	mux := http.NewServeMux()

	// pet endpoints
	mux.HandleFunc("POST /pet", petHandler.AddPet)
	mux.HandleFunc("PUT /pet", petHandler.UpdatePet)
	mux.HandleFunc("GET /pet/findByStatus", petHandler.FindPetsByStatus)
	//mux.HandleFunc("GET /pet/findByTags", petHandler.FindPetsByTags)
	//mux.HandleFunc("GET /pet/{petId}", petHandler.GetById)
	//mux.HandleFunc("POST /pet/{petId}", petHandler.UpdatePetByForm)
	//mux.HandleFunc("DELETE /pet/{petId}", petHandler.DeletePet)
	//mux.HandleFunc("POST /pet/{petId}/uploadImage", petHandler.UploadImage)

	// store endpoints
	_ = storeHandler
	// mux.HandleFunc("GET /store/inventory", storeHandler.GetStoreInventory)
	// mux.HandleFunc("POST /store/order", storeHandler.CreateNewOrder)
	// mux.HandleFunc("GET /store/order/{storeId}", storeHandler.GetOrderById)
	// mux.HandleFunc("DELETE /store/order/{storeId}", storeHandler.DeleteOrder)

	logger.Info("api-petstore-service-layer is up and running", "port", 8080)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		logger.Error("server failed to start", "error", err)
		os.Exit(1)
	}
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
