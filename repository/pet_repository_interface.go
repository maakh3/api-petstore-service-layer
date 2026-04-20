package repository

import (
	"github.com/maakh3/api-petstore-service-layer/models"
)

//go:generate mockgen -source=pet_repository_interface.go -destination=../mocks/mock_pet_repository.go -package=mocks github.com/maakh3/api-petstore-service-layer/repository PetRepositoryInterface
type PetRepositoryInterface interface {
	AddPet(pet models.Pet) (models.Pet, error)
	UpdatePet(pet models.Pet) (models.Pet, error)
	FindPetsByStatus(status string) ([]models.Pet, error)
	FindPetsByTags(tags []models.Tag) ([]models.Pet, error)
	GetById(id int64) (models.Pet, error)
	UpdatePetByForm(id int64, name *string, status *string) (models.Pet, error)
	UploadImage(id int64, imageUrl string) (models.Pet, error)
	DeletePet(id int64) error
}
