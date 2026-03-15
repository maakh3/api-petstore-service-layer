package repository

import (
	"models"
)

// todo populate with the appropriate fields
type PetRepository struct {
	mu     sync.RWMutex          // to handle concurrent access to the pets map
	pets   map[int64]*models.Pet // in-memory storage for pets, keyed by their ID
	nextId int                   // to simulate auto-incrementing primary key
}

func NewPetRepository() *PetRepository {
	return &PetRepository{
		pets:   make(map[int64]*models.Pet),
		nextId: 1,
	}
}
