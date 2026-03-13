package main

import (
	"fmt"
	"net/http"
)

func main() {

	// pet endpoints
	http.HandleFunc("PUT /pet", _)
	http.HandleFunc("POST /pet", _)
	http.HandleFunc("GET /pet/findByStatus", _)
	http.HandleFunc("GET /pet/findByTags", _)
	http.HandleFunc("GET /pet/{petId}", _)
	http.HandleFunc("POST /pet/{petId}", _)
	http.HandleFunc("DELETE /pet/{petId}", _)
	http.HandleFunc("POST /pet/{petId}/uploadImage", _)

	// store endpoints
	http.HandleFunc("GET /store/inventory", _)
	http.HandleFunc("POST /store/order", _)
	http.HandleFunc("GET /store/order/{storeId}", _)
	http.HandleFunc("DELETE /store/order/{storeId}", _)

	http.ListenAndServe(":8080", nil)
}
