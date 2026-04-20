package services

import (
	"github.com/maakh3/api-petstore-service-layer/models"
)

//go:generate mockgen -source=pet_service_interface.go -destination=../mocks/mock_pet_service.go -package=mocks github.com/maakh3/api-petstore-service-layer/services PetServiceInterface
type PetServiceInterface interface {
	AddPet(pet models.Pet) (models.Pet, error)
	UpdatePet(pet models.Pet) (models.Pet, error)
	FindPetsByStatus(status string) ([]models.Pet, error)
	FindPetsByTags(tags []models.Tag) ([]models.Pet, error)
	GetPetById(id int) (models.Pet, error)
	DeletePet(id int) error
	UpdatePetByForm(id int, name string, status string) (models.Pet, error)
	UploadImage(id int, imageData []byte) (models.Pet, error)
}
