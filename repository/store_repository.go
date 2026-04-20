package repository

import (
	"database/sql"
	"log/slog"

	"github.com/maakh3/api-petstore-service-layer/models"
)

type StoreRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewStoreRepository(db *sql.DB, logger ...*slog.Logger) *StoreRepository {
	selectedLogger := slog.Default()
	if len(logger) > 0 && logger[0] != nil {
		selectedLogger = logger[0]
	}

	return &StoreRepository{
		db:     db,
		logger: selectedLogger,
	}
}

func (r *StoreRepository) GetInventory() {
	//TODO implement me
	panic("implement me")
}

func (r *StoreRepository) PlaceOrder() {
	//TODO implement me
	panic("implement me")
}

func (r *StoreRepository) FindOrderById(orderId int64) (models.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (r *StoreRepository) DeleteOrder(orderId int64) error {
	//TODO implement med
	panic("implement me")
}
