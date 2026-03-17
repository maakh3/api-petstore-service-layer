package services

import (
	"api-petstore-service-layer/models"
	"api-petstore-service-layer/repository"
)

type PetService struct {
	repo *repository.PetRepository
}

func NewPetService(repo *repository.PetRepository) *PetService {
	return &PetService{repo: repo}
}

func (s *PetService) AddPet(pet models.Pet) (models.Pet, error) {
	return s.repo.AddPet(pet)
}

func (s *PetService) UpdatePet(pet models.Pet) (models.Pet, error) {
	updated, err := s.repo.UpdatePet(pet)
	if err != nil {
		return models.Pet{}, ErrPetNotFound
	}
	return updated, nil
}

func (s *PetService) FindPetsByStatus(status string) ([]models.Pet, error) {
	return s.repo.FindPetsByStatus(status)
}

