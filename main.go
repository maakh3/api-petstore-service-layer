package main

import (
	"fmt"
	"net/http"

	"repository"
	"services"
	"handlers"
)

func main() {

	petRepo := repository.NewPetRepository()
	petService := services.NewPetService(petRepo)
	petHandler := handlers.NewPetHandler(petService)

	mux := http.NewServeMux()

	// pet endpoints
	mux.HandleFunc("PUT /pet", _)
	mux.HandleFunc("POST /pet", _)
	mux.HandleFunc("GET /pet/findByStatus", _)
	mux.HandleFunc("GET /pet/findByTags", _)
	mux.HandleFunc("GET /pet/{petId}", _)
	mux.HandleFunc("POST /pet/{petId}", _)
	mux.HandleFunc("DELETE /pet/{petId}", _)
	mux.HandleFunc("POST /pet/{petId}/uploadImage", _)

	// store endpoints
	mux.HandleFunc("GET /store/inventory", _)
	mux.HandleFunc("POST /store/order", _)
	mux.HandleFunc("GET /store/order/{storeId}", _)
	mux.HandleFunc("DELETE /store/order/{storeId}", _)

	http.ListenAndServe(":8080", mux)
}
