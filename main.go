package main

import (
	"fmt"
	"net/http"

	"handlers"
	"repository"
	"services"
)

func main() {

	petRepo := repository.NewPetRepository()
	petService := services.NewPetService(petRepo)
	petHandler := handlers.NewPetHandler(petService)

	storeRepo := repository.NewStoreRepository()
	storeService := services.NewStoreService(storeRepo)
	storeHandler := handlers.NewStoreHandler(storeService)

	mux := http.NewServeMux()

	// pet endpoints
	mux.HandleFunc("POST /pet", petHandler.AddPet)
	mux.HandleFunc("PUT /pet", petHandler.UpdatePet)
	mux.HandleFunc("GET /pet/findByStatus", petHandler.FindPetsByStatus)
	mux.HandleFunc("GET /pet/findByTags", petHandler.FindPetsByTags)
	mux.HandleFunc("GET /pet/{petId}", petHandler.GetById)
	mux.HandleFunc("POST /pet/{petId}", petHandler.UpdatePetByForm)
	mux.HandleFunc("DELETE /pet/{petId}", petHandler.DeletePet)
	mux.HandleFunc("POST /pet/{petId}/uploadImage", petHandler.UploadImage)

	// store endpoints
	mux.HandleFunc("GET /store/inventory", storeHandler.GetStoreInventory)
	mux.HandleFunc("POST /store/order", storeHandler.CreateNewOrder)
	mux.HandleFunc("GET /store/order/{storeId}", storeHandler.GetOrderById)
	mux.HandleFunc("DELETE /store/order/{storeId}", storeHandler.DeleteOrder)

	http.ListenAndServe(":8080", mux)
}
