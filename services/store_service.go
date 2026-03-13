package services

import (
	"repository"
)

type StoreService struct {
	repo *repository.StoreRepository
}

func NewStoreService(repo *repository.StoreRepository) *StoreService {
	return &StoreService{repo: repo}
}
