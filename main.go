package main

import (
	"log"
	"net/http"

	"api-petstore-service-layer/handlers"
	"api-petstore-service-layer/repository"
	"api-petstore-service-layer/services"
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
	//mux.HandleFunc("GET /pet/findByStatus", petHandler.FindPetsByStatus)
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

	log.Println("\napi-petstore-service-layer IS UP & RUNNING ON PORT 8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
