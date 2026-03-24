package repository

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/maakh3/api-petstore-service-layer/models"
)

type PetRepository struct {
	mu     sync.RWMutex // to handle concurrent access to the pets map
	pets   map[int64]*models.Pet
	nextId int64
	logger *slog.Logger
}

func NewPetRepository(logger ...*slog.Logger) *PetRepository {
	selectedLogger := slog.Default()
	if len(logger) > 0 && logger[0] != nil {
		selectedLogger = logger[0]
	}

	return &PetRepository{
		pets:   make(map[int64]*models.Pet),
		nextId: 1,
		logger: selectedLogger,
	}
}

func (r *PetRepository) AddPet(pet models.Pet) (models.Pet, error) {
	r.logger.Debug("repository add pet", "name", pet.Name, "status", pet.Status)
	r.mu.Lock()         // pessimistically lock the mutex to ensure thread safety when modifying the pets map
	defer r.mu.Unlock() // ensure the mutex is unlocked after the function returns

	pet.Id = int(r.nextId)       // assign a new Id to the pet
	r.pets[int64(pet.Id)] = &pet // store the pet in the map
	r.nextId++                   // increment the nextId for the next pet
	r.logger.Info("repository added pet", "pet_id", pet.Id)

	return pet, nil
}

func (r *PetRepository) UpdatePet(pet models.Pet) (models.Pet, error) {
	r.logger.Debug("repository update pet", "pet_id", pet.Id)
	r.mu.Lock()
	defer r.mu.Unlock()

	existingPet, exists := r.pets[int64(pet.Id)]
	if !exists {
		r.logger.Info("repository pet not found during update", "pet_id", pet.Id)
		return models.Pet{}, fmt.Errorf("pet with Id %d not found", pet.Id)
	}

	*existingPet = pet // update the existing pet with the new data
	r.logger.Info("repository updated pet", "pet_id", pet.Id, "status", pet.Status)

	return *existingPet, nil
}

func (r *PetRepository) FindPetsByStatus(status string) ([]models.Pet, error) {
	r.logger.Debug("repository find pets by status", "status", status)
	r.mu.RLock() // use read lock for concurrent reads
	defer r.mu.RUnlock()

	pets := make([]models.Pet, 0)
	for _, pet := range r.pets {
		if pet.Status == status {
			pets = append(pets, *pet)
		}
	}
	r.logger.Info("repository found pets by status", "status", status, "count", len(pets))
	return pets, nil
}

func (r *PetRepository) FindPetsByTags(tags []models.Tag) ([]models.Pet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var pets = make([]models.Pet, 0)
	for _, pet := range r.pets {
		if containsAllTags(pet.Tags, tags) {
			pets = append(pets, *pet)
		}
	}
	return pets, nil
}

// helper function to check if one pet's tags contain all the specified tags
func containsAllTags(petTags []models.Tag, searchTags []models.Tag) bool {
	tagSet := make(map[string]struct{}) // create a set of the pet's tags for efficient lookup
	for _, tag := range petTags {       // iterate over the pet's tags and add them to the set
		tagSet[tag.Name] = struct{}{} // the value doesn't matter, we just care about the keys for existence checks
	}

	for _, searchTag := range searchTags { // iterate over the search tags and check if each one exists in the pet's tag set
		if _, exists := tagSet[searchTag.Name]; !exists { // if any search tag is not found in the pet's tags, end the search
			return false
		}
	}
	return true
}
