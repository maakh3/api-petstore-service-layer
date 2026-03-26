package services

import (
	"github.com/maakh3/api-petstore-service-layer/repository"
)

type StoreService struct {
	repo *repository.StoreRepository
}

func NewStoreService(repo *repository.StoreRepository) *StoreService {
	return &StoreService{repo: repo}
}
