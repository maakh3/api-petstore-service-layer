package repository

import (
	"api-petstore-service-layer/models"
	"fmt"
	"sync"
)

type PetRepository struct {
	mu     sync.RWMutex // to handle concurrent access to the pets map
	pets   map[int64]*models.Pet
	nextId int64
}

func NewPetRepository() *PetRepository {
	return &PetRepository{
		pets:   make(map[int64]*models.Pet),
		nextId: 1,
	}
}

func (r *PetRepository) AddPet(pet models.Pet) (models.Pet, error) {
	r.mu.Lock()         // lock the mutex to ensure thread safety when modifying the pets map
	defer r.mu.Unlock() // ensure the mutex is unlocked after the function returns

	pet.Id = int(r.nextId)       // assign a new Id to the pet
	r.pets[int64(pet.Id)] = &pet // store the pet in the map
	r.nextId++                   // increment the nextId for the next pet

	return pet, nil
}

func (r *PetRepository) UpdatePet(pet models.Pet) (models.Pet, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	existingPet, exists := r.pets[int64(pet.Id)]
	if !exists {
		return models.Pet{}, fmt.Errorf("pet with Id %d not found", pet.Id)
	}

	*existingPet = pet // update the existing pet with the new data

	return *existingPet, nil
}
