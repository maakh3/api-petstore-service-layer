package handlers

import (
	"net/http"
)

//go:generate mockgen -source=pet_handler_interface.go -destination=../mocks/mock_pet_handler.go -package=mocks github.com/maakh3/api-petstore-service-layer/handler PetHandlerInterface
type PetHandlerInterface interface {
	AddPet(w http.ResponseWriter, r *http.Request)
	UpdatePet(w http.ResponseWriter, r *http.Request)
	FindPetsByStatus(w http.ResponseWriter, r *http.Request)
	FindPetsByTags(w http.ResponseWriter, r *http.Request)
	GetById(w http.ResponseWriter, r *http.Request)
}
