package handlers

import (
	//"net/http"

	"github.com/maakh3/api-petstore-service-layer/services"
)

type StoreHandler struct {
	service *services.StoreService
}

func NewStoreHandler(service *services.StoreService) *StoreHandler {

	return &StoreHandler{service: service}
}
