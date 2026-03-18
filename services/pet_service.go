package services

import (
	"log/slog"

	"github.com/maakh3/api-petstore-service-layer/models"
	"github.com/maakh3/api-petstore-service-layer/repository"
)

type PetService struct {
	repo   *repository.PetRepository
	logger *slog.Logger
}

func NewPetService(repo *repository.PetRepository, logger ...*slog.Logger) *PetService {
	selectedLogger := slog.Default()
	if len(logger) > 0 && logger[0] != nil {
		selectedLogger = logger[0]
	}

	return &PetService{repo: repo, logger: selectedLogger}
}

func (s *PetService) AddPet(pet models.Pet) (models.Pet, error) {
	s.logger.Debug("service add pet", "name", pet.Name, "status", pet.Status)
	createdPet, err := s.repo.AddPet(pet)
	if err != nil {
		s.logger.Error("service failed to add pet", "name", pet.Name, "status", pet.Status, "error", err)
		return models.Pet{}, err
	}

	s.logger.Info("service added pet", "pet_id", createdPet.Id)
	return createdPet, nil
}

func (s *PetService) UpdatePet(pet models.Pet) (models.Pet, error) {
	s.logger.Debug("service update pet", "pet_id", pet.Id)
	updated, err := s.repo.UpdatePet(pet)
	if err != nil {
		s.logger.Info("service pet not found during update", "pet_id", pet.Id)
		return models.Pet{}, ErrPetNotFound
	}

	s.logger.Info("service updated pet", "pet_id", updated.Id, "status", updated.Status)
	return updated, nil
}

func (s *PetService) FindPetsByStatus(status string) ([]models.Pet, error) {
	s.logger.Debug("service find pets by status", "status", status)
	pets, err := s.repo.FindPetsByStatus(status)
	if err != nil {
		s.logger.Error("service failed to find pets by status", "status", status, "error", err)
		return nil, err
	}

	s.logger.Info("service found pets by status", "status", status, "count", len(pets))
	return pets, nil
}
