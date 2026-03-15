package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"services"
	"models"
)

type PetHandler struct {
	service *services.PetService
}

func NewPetHandler(service *services.PetService) *PetHandler {
	return &PetHandler{service: service}
}

// the handler functions
func (p *PetHandler) AddPet(w http.ResponseWriter, r *http.Request) {
	var pet models.Pet
	err := json.NewDecoder(r.Body).Decode(&pet)
	if err != nil {
		http.Error(w, "Invalid request payload", 400)
		return
	}

	createdPet, err := p.service.AddPet(pet)
	if err != nil {
		http.Error(w, "Failed to add pet", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201) // Status Created
	json.NewEncoder(w).Encode(createdPet)
}

func (p *PetHandler) UpdatePet(w http.ResponseWriter, r *http.Request) {
	var pet models.Pet
	err := json.NewDecoder(r.Body).Decode(&pet)
	if err != nil {
		http.Error(w, "Invalid request payload", 400) // Status Bad Request
		return
	}

	updatedPet, err := p.service.UpdatePet(pet)
	if err != nil {
		if errors.Is(err, services.ErrPetNotFound) {
			http.Error(w, "Pet not found", 404) // Status Not Found
			return
		}
		http.Error(w, "Failed to update pet", 500) // Status Internal Server Error
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200) // Status OK
	json.NewEncoder(w).Encode(updatedPet)
}

func (p *PetHandler) FindPetByStatus(w http.ResponseWriter, r *http.Request) {
	petStatus := r.URL.Query().Get("status")
	// petStatus is a comma-separated list of statuses, e.g. "available,pending"

	if petStatus == "" {
		http.Error(w, "Status query parameter is required", 400)
		return
	}

	pets, err := p.service.FindPetByStatus(petStatus)
	if err != nil {
		http.Error(w, "Failed to retrieve pets", 500)
		// todo - maybe all these 500s become 422s once i implement the service layer?
		// because if the service layer validates the input and returns an error for invalid status values,
		// then that should be a 422 Unprocessable Entity rather than a 500 Internal Server Error.
		// todo - also should i return an error message in the body here?
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(pets)
}

func (p *PetHandler) FindPetByTags(w http.ResponseWriter, r *http.Request) {
	tags := r.URL.Query().Get("tags")
	// tags is also a comma-separated list of tags, e.g. "tag1,tag2"

	if tags == "" {
		http.Error(w, "Tags query parameter is required", 400)
		return
	}

	pets, err := p.service.FindPetByTags(tags)
	if err != nil {
		http.Error(w, "Failed to retrieve pets", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(pets)
}

func (p *PetHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("petId"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid pet ID", 400)
		return
	}

	pet, err := p.service.GetPetById(id)
	if err != nil {
		if errors.Is(err, services.ErrPetNotFound) {
			http.Error(w, "Pet not found", 404)
			return
		}
		http.Error(w, "Failed to retrieve pet", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(pet)
}

func (p *PetHandler) UpdatePetByForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("petId"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid pet ID", 400)
		return
	}

	pet, err := p.service.GetPetById(id)
	if err != nil {
		if errors.Is(err, services.ErrPetNotFound) {
			http.Error(w, "Pet not found", 404)
			return
		}
		http.Error(w, "Failed to retrieve pet", 500)
		return
	}

	if name := r.FormValue("name"); name != "" {
		pet.Name = name
	}
	if status := r.FormValue("status"); status != "" {
		pet.Status = status
	}

	updatedPet, err := p.service.UpdatePet(pet)
	if err != nil {
		http.Error(w, "Failed to update pet", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(updatedPet)

}

func (p *PetHandler) DeletePet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("petId"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid pet ID", 400)
		return
	}

	err = p.service.DeletePet(id)
	if err != nil {
		if errors.Is(err, services.ErrPetNotFound) {
			http.Error(w, "Pet not found", 404)
			return
		}
		http.Error(w, "Failed to delete pet", 500)
		return
	}

	w.WriteHeader(204) // Status No Content
	//w.Write([]byte("Pet deleted successfully"))
}

func (p *PetHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("petId"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid pet ID", 400)
		return
	}
	if _, err := p.service.GetPetById(id); err != nil {
		if errors.Is(err, services.ErrPetNotFound) {
			http.Error(w, "Pet not found", 404)
			return
		}
		http.Error(w, "Failed to retrieve pet", 500)
		return
	}

	err = r.ParseMultipartForm(10 << 20) // limits upload size to 10MB
	if err != nil {
		http.Error(w, "Failed to parse multipart form", 400)
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to retrieve file from form data", 400)
		return
	}
	defer file.Close()

	err = p.service.UploadPetImage(id, file)
	if err != nil {
		if errors.Is(err, services.ErrPetNotFound) {
			http.Error(w, "Pet not found", 404)
			return
		}
		http.Error(w, "Failed to upload image", 500)
		return
	}

	w.WriteHeader(200)
	//w.Write([]byte("Image uploaded successfully"))
}
