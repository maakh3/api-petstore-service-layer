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

func (r *PetRepository) FindPetsByStatus(status string) ([]models.Pet, error) {
	r.mu.RLock() // use read lock for concurrent reads
	defer r.mu.RUnlock()

	pets := []models.Pet{}
	for _, pet := range r.pets {
		if pet.Status == status {
			pets = append(pets, *pet)
		}
	}
	return pets, nil
}
//
//func (r *PetRepository) FindPetsByTags(tags []models.Tag) ([]models.Pet, error) {
//	r.mu.RLock()
//	defer r.mu.RUnlock()
//
//	var pets []models.Pet
//	for _, pet := range r.pets {
//		if containsAllTags(pet.Tags, tags) {
//			pets = append(pets, *pet)
//		}
//	}
//	return pets, nil
//}
//
//func (r *PetRepository) GetById(id int64) (models.Pet, error) {
//	r.mu.RLock()
//	defer r.mu.RUnlock()
//
//	pet, exists := r.pets[id]
//	if !exists {
//		return models.Pet{}, fmt.Errorf("pet with Id %d not found", id)
//	}
//	return *pet, nil
//}
//
//func (r *PetRepository) UploadImage(id int64, imageUrl string) (models.Pet, error) {
//	r.mu.Lock()
//	defer r.mu.Unlock()
//
//	pet, exists := r.pets[id]
//	if !exists {
//		return models.Pet{}, fmt.Errorf("pet with Id %d not found", id)
//	}
//
//	pet.PhotoUrls = append(pet.PhotoUrls, imageUrl) // add the new image URL to the pet's photo URLs
//	return *pet, nil
//}
//
//func (r *PetRepository) Delete(id int64) error {
//	r.mu.Lock()
//	defer r.mu.Unlock()
//
//	if _, exists := r.pets[id]; !exists {
//		return fmt.Errorf("pet with Id %d not found", id)
//	}
//
//	delete(r.pets, id)
//	return nil
//}
//
//// helper function to check if one pet's tags contain all the specified tags
//func containsAllTags(petTags []models.Tag, searchTags []models.Tag) bool {
//	tagSet := make(map[string]struct{}) // create a set of the pet's tags for efficient lookup
//	for _, tag := range petTags {       // iterate over the pet's tags and add them to the set
//		tagSet[tag.Name] = struct{}{} // the value doesn't matter, we just care about the keys for existence checks
//	}
//
//	for _, searchTag := range searchTags { // iterate over the search tags and check if each one exists in the pet's tag set
//		if _, exists := tagSet[searchTag.Name]; !exists { // if any search tag is not found in the pet's tags, end the search
//			return false
//		}
//	}
//	return true
//}
