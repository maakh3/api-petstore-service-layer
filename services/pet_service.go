package services

import (
	"repository"
)

type PetService struct {
	repo *repository.PetRepository
}

func NewPetService(repo *repository.PetRepository) *PetService {
	return &PetService{repo: repo}
}
