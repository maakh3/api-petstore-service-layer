package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/lib/pq"
	"github.com/maakh3/api-petstore-service-layer/models"
)

type PetRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewPetRepository(db *sql.DB, logger ...*slog.Logger) *PetRepository {
	selectedLogger := slog.Default()
	if len(logger) > 0 && logger[0] != nil {
		selectedLogger = logger[0]
	}

	return &PetRepository{
		db:     db,
		logger: selectedLogger,
	}
}

func (r *PetRepository) AddPet(pet models.Pet) (models.Pet, error) {
	r.logger.Debug("repository add pet", "pet_id", pet.Id)

	categoryJson, err := json.Marshal(pet.Category)
	if err != nil {
		r.logger.Error("repository failed to marshal category", "error", err, "pet_id", pet.Id)
		return models.Pet{}, fmt.Errorf("failed to marshal category: %w", err)
	}

	tagsJson, err := json.Marshal(pet.Tags)
	if err != nil {
		r.logger.Error("repository failed to marshal tags", "error", err, "pet_id", pet.Id)
		return models.Pet{}, fmt.Errorf("failed to marshal tags: %w", err)
	}

	const q = `
		INSERT INTO pets (id, name, category, photo_urls, tags, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
`
	if err := r.db.QueryRow(
		q,
		pet.Id,
		pet.Name,
		categoryJson,
		pq.Array(pet.PhotoUrls),
		tagsJson,
		pet.Status,
	).Scan(&pet.Id); err != nil {
		return models.Pet{}, fmt.Errorf("error adding pet: %w", err)
	}

	r.logger.Info("repository added pet", "pet_id", pet.Id, "status", pet.Status)
	return pet, nil
}

func (r *PetRepository) UpdatePet(pet models.Pet) (models.Pet, error) {
	categoryJSON, err := json.Marshal(pet.Category)
	if err != nil {
		return models.Pet{}, fmt.Errorf("marshal category: %w", err)
	}
	tagsJSON, err := json.Marshal(pet.Tags)
	if err != nil {
		return models.Pet{}, fmt.Errorf("marshal tags: %w", err)
	}

	const q = `
		UPDATE pets
		SET name = $2, status = $3, category = $4::jsonb, tags = $5::jsonb, photo_urls = $6
		WHERE id = $1
	`
	res, err := r.db.Exec(q, pet.Id, pet.Name, pet.Status, categoryJSON, tagsJSON, pq.Array(pet.PhotoUrls))
	if err != nil {
		return models.Pet{}, fmt.Errorf("update pet: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return models.Pet{}, fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return models.Pet{}, fmt.Errorf("pet with Id %d not found", pet.Id)
	}
	return pet, nil
}

func (r *PetRepository) FindPetsByStatus(status string) ([]models.Pet, error) {
	//TODO implement me
	panic("implement me")
}

func (r *PetRepository) FindPetsByTags(tags []models.Tag) ([]models.Pet, error) {
	//TODO implement me
	panic("implement me")
}

func (r *PetRepository) GetById(id int64) (models.Pet, error) {
	//TODO implement me
	panic("implement me")
}

func (r *PetRepository) UpdatePetByForm(id int64, name *string, status *string) (models.Pet, error) {
	//TODO implement me
	panic("implement me")
}

func (r *PetRepository) UploadImage(id int64, imageData []byte) (models.Pet, error) {
	//TODO implement me
	panic("implement me")
}

func (r *PetRepository) DeletePet(id int64) error {
	//TODO implement me
	panic("implement me")
}
