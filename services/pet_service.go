package services

import (
	"log/slog"

	"github.com/maakh3/api-petstore-service-layer/models"
	"github.com/maakh3/api-petstore-service-layer/repository"
)

type PetService struct {
	repo   repository.PetRepositoryInterface
	logger *slog.Logger
}

func NewPetService(repo repository.PetRepositoryInterface, logger ...*slog.Logger) *PetService {
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

func (s *PetService) FindPetsByTags(tags []models.Tag) ([]models.Pet, error) {
	stringTags := make([]string, len(tags))
	for i, tag := range tags {
		stringTags[i] = tag.Name
	}
	s.logger.Debug("service find pets by tags", "tags", stringTags)

	pets, err := s.repo.FindPetsByTags(tags)
	if err != nil {
		s.logger.Error("service failed to find pets by tags", "tags", tags, "error", err)
		return nil, err
	}

	petsCount := len(pets)
	s.logger.Info("service found pets by tags", "tags", stringTags, "count", petsCount)
	return pets, nil
}

func (s *PetService) GetPetById(id int) (models.Pet, error) {
	s.logger.Debug("service get by id", "pet_id", id)
	updated, err := s.repo.GetById(int64(id))
	if err != nil {
		s.logger.Info("service pet not found during get by id", "pet_id", id)
		return models.Pet{}, ErrPetNotFound
	}

	s.logger.Info("service get by id", "pet_id", updated.Id, "status", updated.Status)
	return updated, nil
}

func (s *PetService) DeletePet(id int) error {
	s.logger.Debug("service delete pet", "pet_id", id)
	err := s.repo.DeletePet(int64(id))
	if err != nil {
		s.logger.Info("service pet not found during delete", "pet_id", id)
		return ErrPetNotFound
	}

	s.logger.Info("service deleted pet", "pet_id", id)
	return nil
}

func (s *PetService) UpdatePetByForm(id int, name string, status string) (models.Pet, error) {
	s.logger.Debug("service update pet by form", "pet_id", id)
	updated, err := s.repo.UpdatePetByForm(int64(id), &name, &status)
	if err != nil {
		s.logger.Info("service pet not found during update by form", "pet_id", id)
		return models.Pet{}, ErrPetNotFound
	}

	s.logger.Info("service updated pet by form", "pet_id", updated.Id, "status", updated.Status)
	return updated, nil
}

func (s *PetService) UploadImage(petId int, imageUrl string) error {
	s.logger.Debug("service upload image", "pet_id", petId)
	_, err := s.repo.UploadImage(int64(petId), imageUrl)
	if err != nil {
		s.logger.Info("service pet not found during upload image", "pet_id", petId)
		return ErrPetNotFound
	}

	s.logger.Info("service uploaded image for pet", "pet_id", petId)
	return nil
}
