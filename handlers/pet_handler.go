package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	//"strconv"

	"github.com/maakh3/api-petstore-service-layer/models"
	"github.com/maakh3/api-petstore-service-layer/services"
)

type PetHandler struct {
	service *services.PetService
	logger  *slog.Logger
}

func NewPetHandler(service *services.PetService, logger ...*slog.Logger) *PetHandler {
	selectedLogger := slog.Default()
	if len(logger) > 0 && logger[0] != nil {
		selectedLogger = logger[0]
	}

	return &PetHandler{service: service, logger: selectedLogger}
}

func (p *PetHandler) AddPet(w http.ResponseWriter, r *http.Request) {
	p.logger.Debug("add pet request received", "method", r.Method, "path", r.URL.Path)

	var pet models.Pet
	err := json.NewDecoder(r.Body).Decode(&pet)
	if err != nil {
		p.logger.Error("failed to decode add pet payload", "error", err)
		http.Error(w, "Invalid request payload", 400)
		return
	}

	createdPet, err := p.service.AddPet(pet)
	if err != nil {
		p.logger.Error("failed to add pet", "name", pet.Name, "status", pet.Status, "error", err)
		http.Error(w, "Failed to add pet", 500)
		return
	}

	p.logger.Info("pet added", "pet_id", createdPet.Id, "status", createdPet.Status)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201) // Status Created
	err = json.NewEncoder(w).Encode(createdPet)
	if err != nil {
		p.logger.Error("failed to encode add pet response", "pet_id", createdPet.Id, "error", err)
		http.Error(w, "Failed to encode response", 500)
		return
	}
}

func (p *PetHandler) UpdatePet(w http.ResponseWriter, r *http.Request) {
	p.logger.Debug("update pet request received", "method", r.Method, "path", r.URL.Path)

	var pet models.Pet
	err := json.NewDecoder(r.Body).Decode(&pet)
	if err != nil {
		p.logger.Error("failed to decode update pet payload", "error", err)
		http.Error(w, "Invalid request payload", 400) // Status Bad Request
		return
	}

	if pet.Id == 0 {
		p.logger.Debug("update pet request missing id")
		http.Error(w, "Pet ID is required for update", 400) // Status Bad Request
		return
	}

	updatedPet, err := p.service.UpdatePet(pet)
	if err != nil {
		if errors.Is(err, services.ErrPetNotFound) {
			p.logger.Info("pet not found during update", "pet_id", pet.Id)
			http.Error(w, "Pet not found", 404) // Status Not Found
			return
		}
		p.logger.Error("failed to update pet", "pet_id", pet.Id, "error", err)
		http.Error(w, "Failed to update pet", 500) // Status Internal Server Error
		return
	}

	p.logger.Info("pet updated", "pet_id", updatedPet.Id, "status", updatedPet.Status)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200) // Status OK
	json.NewEncoder(w).Encode(updatedPet)
}

func (p *PetHandler) FindPetsByStatus(w http.ResponseWriter, r *http.Request) {
	petStatus := r.URL.Query().Get("status")
	p.logger.Debug("find pets by status request received", "method", r.Method, "path", r.URL.Path, "status", petStatus)
	// petStatus is a comma-separated list of statuses, e.g. "available,pending"

	if petStatus == "" {
		p.logger.Debug("find pets by status missing status query parameter")
		http.Error(w, "Status query parameter is required", 400)
		return
	}

	pets, err := p.service.FindPetsByStatus(petStatus)
	if err != nil {
		p.logger.Error("failed to find pets by status", "status", petStatus, "error", err)
		http.Error(w, "Failed to retrieve pets", 500)
		return
	}

	p.logger.Info("pets retrieved by status", "status", petStatus, "count", len(pets))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(pets)
}

func (p *PetHandler) FindPetsByTags(w http.ResponseWriter, r *http.Request) {
	tags := r.URL.Query().Get("tags")
	// tags is also a comma-separated list of tags, e.g. "tag1,tag2"

	if tags == "" {
		http.Error(w, "Tags query parameter is required", 400)
		return
	}

	tagList := strings.Split(tags, ",")
	modelTags := make([]models.Tag, 0, len(tagList))
	for _, tagList := range tagList {
		name := strings.TrimSpace(tagList)
		if name == "" {
			continue
		}
		modelTags = append(modelTags, models.Tag{Name: name})
	}

	pets, err := p.service.FindPetsByTags(modelTags)
	if err != nil {
		http.Error(w, "Failed to retrieve pets", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(pets)
}
