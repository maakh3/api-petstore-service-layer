package handlers

import (
	"api-petstore-service-layer/models"
	"api-petstore-service-layer/repository"
	"api-petstore-service-layer/services"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func newTestPetHandler() *PetHandler {
	repo := repository.NewPetRepository()
	service := services.NewPetService(repo)
	return NewPetHandler(service)
}

func TestPetHandlerAddPet(t *testing.T) {
	h := newTestPetHandler()

	body := `{"name":"fido","status":"available","tags":[{"id":1,"name":"small"}]}`
	req := httptest.NewRequest(http.MethodPost, "/pet", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.AddPet(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("AddPet() status = %d, want %d", res.StatusCode, http.StatusCreated)
	}

	var got models.Pet
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode AddPet() response: %v", err)
	}

	if got.Id != 1 {
		t.Fatalf("AddPet() Id = %d, want 1", got.Id)
	}
	if got.Name != "fido" {
		t.Fatalf("AddPet() Name = %q, want %q", got.Name, "fido")
	}
	if got.Status != "available" {
		t.Fatalf("AddPet() Status = %q, want %q", got.Status, "available")
	}
}

func TestPetHandlerAddPet_InvalidPayload(t *testing.T) {
	h := newTestPetHandler()

	req := httptest.NewRequest(http.MethodPost, "/pet", strings.NewReader("{"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.AddPet(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("AddPet() status = %d, want %d", res.StatusCode, http.StatusBadRequest)
	}
}

func TestPetHandlerUpdatePet(t *testing.T) {
	repo := repository.NewPetRepository()
	service := services.NewPetService(repo)
	h := NewPetHandler(service)

	created, err := service.AddPet(models.Pet{Name: "fido", Status: "available"})
	if err != nil {
		t.Fatalf("setup AddPet() unexpected error: %v", err)
	}

	body := `{"id":` + strconv.Itoa(created.Id) + `,"name":"fido-updated","status":"sold"}`
	req := httptest.NewRequest(http.MethodPut, "/pet", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.UpdatePet(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("UpdatePet() status = %d, want %d", res.StatusCode, http.StatusOK)
	}

	var got models.Pet
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode UpdatePet() response: %v", err)
	}

	if got.Id != created.Id {
		t.Fatalf("UpdatePet() Id = %d, want %d", got.Id, created.Id)
	}
	if got.Name != "fido-updated" {
		t.Fatalf("UpdatePet() Name = %q, want %q", got.Name, "fido-updated")
	}
	if got.Status != "sold" {
		t.Fatalf("UpdatePet() Status = %q, want %q", got.Status, "sold")
	}
}

func TestPetHandlerUpdatePet_InvalidPayload(t *testing.T) {
	h := newTestPetHandler()

	req := httptest.NewRequest(http.MethodPut, "/pet", strings.NewReader("{"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.UpdatePet(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("UpdatePet() status = %d, want %d", res.StatusCode, http.StatusBadRequest)
	}
}

func TestPetHandlerUpdatePet_MissingId(t *testing.T) {
	h := newTestPetHandler()

	body := `{"name":"fido-updated","status":"sold"}`
	req := httptest.NewRequest(http.MethodPut, "/pet", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.UpdatePet(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("UpdatePet() status = %d, want %d", res.StatusCode, http.StatusBadRequest)
	}
}

func TestPetHandlerUpdatePet_NotFound(t *testing.T) {
	h := newTestPetHandler()

	body := `{"id":999,"name":"ghost","status":"available"}`
	req := httptest.NewRequest(http.MethodPut, "/pet", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.UpdatePet(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("UpdatePet() status = %d, want %d", res.StatusCode, http.StatusNotFound)
	}
}


