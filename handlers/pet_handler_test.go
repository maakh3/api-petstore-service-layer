package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/maakh3/api-petstore-service-layer/mocks"
	"github.com/maakh3/api-petstore-service-layer/models"
	"github.com/maakh3/api-petstore-service-layer/services"
	"go.uber.org/mock/gomock"
)

func TestPetHandler_AddPet(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := mocks.NewMockPetServiceInterface(ctrl)
		h := NewPetHandler(mockService)

		input := models.Pet{Name: "fido", Status: "available", Tags: []models.Tag{{Id: 1, Name: "small"}}}
		created := models.Pet{Id: 1, Name: "fido", Status: "available", Tags: []models.Tag{{Id: 1, Name: "small"}}}
		mockService.EXPECT().AddPet(input).Return(created, nil)

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
	})
	t.Run("invalid payload", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := mocks.NewMockPetServiceInterface(ctrl)
		h := NewPetHandler(mockService)

		req := httptest.NewRequest(http.MethodPost, "/pet", strings.NewReader("{"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.AddPet(w, req)

		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Fatalf("AddPet() status = %d, want %d", res.StatusCode, http.StatusBadRequest)
		}
	})
	t.Run("service error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := mocks.NewMockPetServiceInterface(ctrl)
		h := NewPetHandler(mockService)

		input := models.Pet{Name: "fido", Status: "available"}
		mockService.EXPECT().AddPet(input).Return(models.Pet{}, errors.New("service error"))

		body := `{"name":"fido","status":"available"}`
		req := httptest.NewRequest(http.MethodPost, "/pet", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.AddPet(w, req)

		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusInternalServerError {
			t.Fatalf("AddPet() status = %d, want %d", res.StatusCode, http.StatusInternalServerError)
		}
	})
}

func TestPetHandler_UpdatePet(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := mocks.NewMockPetServiceInterface(ctrl)
		h := NewPetHandler(mockService)

		input := models.Pet{Id: 1, Name: "fido-updated", Status: "sold"}
		want := models.Pet{Id: 1, Name: "fido-updated", Status: "sold"}
		mockService.EXPECT().UpdatePet(input).Return(want, nil)

		body := `{"id":1,"name":"fido-updated","status":"sold"}`
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

		if got.Id != 1 {
			t.Fatalf("UpdatePet() Id = %d, want 1", got.Id)
		}
		if got.Name != "fido-updated" {
			t.Fatalf("UpdatePet() Name = %q, want %q", got.Name, "fido-updated")
		}
		if got.Status != "sold" {
			t.Fatalf("UpdatePet() Status = %q, want %q", got.Status, "sold")
		}
	})
	t.Run("invalid payload", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := mocks.NewMockPetServiceInterface(ctrl)
		h := NewPetHandler(mockService)

		req := httptest.NewRequest(http.MethodPut, "/pet", strings.NewReader("{"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.UpdatePet(w, req)

		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Fatalf("UpdatePet() status = %d, want %d", res.StatusCode, http.StatusBadRequest)
		}
	})
	t.Run("missing id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := mocks.NewMockPetServiceInterface(ctrl)
		h := NewPetHandler(mockService)

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
	})
	t.Run("not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := mocks.NewMockPetServiceInterface(ctrl)
		h := NewPetHandler(mockService)

		input := models.Pet{Id: 999, Name: "ghost", Status: "available"}
		mockService.EXPECT().UpdatePet(input).Return(models.Pet{}, services.ErrPetNotFound)

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
	})
}

func TestPetHandler_FindPetsByStatus(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := mocks.NewMockPetServiceInterface(ctrl)
		h := NewPetHandler(mockService)

		want := []models.Pet{
			{Id: 1, Name: "fido", Status: "available"},
			{Id: 2, Name: "rex", Status: "available"},
		}
		mockService.EXPECT().FindPetsByStatus("available").Return(want, nil)

		req := httptest.NewRequest(http.MethodGet, "/pet/findByStatus?status=available", nil)
		w := httptest.NewRecorder()

		h.FindPetsByStatus(w, req)

		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Fatalf("FindPetsByStatus() status = %d, want %d", res.StatusCode, http.StatusOK)
		}

		var got []models.Pet
		if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode FindPetsByStatus() response: %v", err)
		}

		if len(got) != 2 {
			t.Fatalf("FindPetsByStatus() returned %d pets, want 2", len(got))
		}
		for _, pet := range got {
			if pet.Status != "available" {
				t.Fatalf("FindPetsByStatus() returned pet with status %q, want %q", pet.Status, "available")
			}
		}
	})
	t.Run("no match", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := mocks.NewMockPetServiceInterface(ctrl)
		h := NewPetHandler(mockService)

		mockService.EXPECT().FindPetsByStatus("pending").Return([]models.Pet{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/pet/findByStatus?status=pending", nil)
		w := httptest.NewRecorder()

		h.FindPetsByStatus(w, req)

		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Fatalf("FindPetsByStatus() status = %d, want %d", res.StatusCode, http.StatusOK)
		}

		var got []models.Pet
		if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode FindPetsByStatus() response: %v", err)
		}

		if len(got) != 0 {
			t.Fatalf("FindPetsByStatus() returned %d pets, want 0", len(got))
		}
	})
	t.Run("missing status", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := mocks.NewMockPetServiceInterface(ctrl)
		h := NewPetHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/pet/findByStatus", nil)
		w := httptest.NewRecorder()

		h.FindPetsByStatus(w, req)

		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Fatalf("FindPetsByStatus() status = %d, want %d", res.StatusCode, http.StatusBadRequest)
		}
	})
}

func TestPetHandler_FindPetsByTags(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := mocks.NewMockPetServiceInterface(ctrl)
		h := NewPetHandler(mockService)

		searchTags := []models.Tag{{Name: "small"}}
		want := []models.Pet{
			{Id: 1, Name: "fido", Status: "available", Tags: []models.Tag{{Id: 1, Name: "small"}}},
			{Id: 3, Name: "mittens", Status: "sold", Tags: []models.Tag{{Id: 1, Name: "small"}}},
		}
		mockService.EXPECT().FindPetsByTags(searchTags).Return(want, nil)

		req := httptest.NewRequest(http.MethodGet, "/pet/findByTags?tags=small", nil)
		w := httptest.NewRecorder()

		h.FindPetsByTags(w, req)

		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Fatalf("FindPetsByTags() status = %d, want %d", res.StatusCode, http.StatusOK)
		}

		var got []models.Pet
		if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode FindPetsByTags() response: %v", err)
		}

		if len(got) != 2 {
			t.Fatalf("FindPetsByTags() returned %d pets, want 2", len(got))
		}
		for _, pet := range got {
			foundSmall := false
			for _, tag := range pet.Tags {
				if tag.Name == "small" {
					foundSmall = true
					break
				}
			}
			if !foundSmall {
				t.Fatalf("FindPetsByTags() returned pet without 'small' tag")
			}
		}
	})
	t.Run("missing tags", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := mocks.NewMockPetServiceInterface(ctrl)
		h := NewPetHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/pet/findByTags", nil)
		w := httptest.NewRecorder()

		h.FindPetsByTags(w, req)

		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Fatalf("FindPetsByTags() status = %d, want %d", res.StatusCode, http.StatusBadRequest)
		}
	})
	t.Run("no match", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := mocks.NewMockPetServiceInterface(ctrl)
		h := NewPetHandler(mockService)

		searchTags := []models.Tag{{Name: "nonexistent"}}
		mockService.EXPECT().FindPetsByTags(searchTags).Return([]models.Pet{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/pet/findByTags?tags=nonexistent", nil)
		w := httptest.NewRecorder()

		h.FindPetsByTags(w, req)

		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Fatalf("FindPetsByTags() status = %d, want %d", res.StatusCode, http.StatusOK)
		}

		var got []models.Pet
		if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode FindPetsByTags() response: %v", err)
		}

		if len(got) != 0 {
			t.Fatalf("FindPetsByTags() returned %d pets, want 0", len(got))
		}
	})
}
