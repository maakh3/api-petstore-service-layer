package services

import (
	"errors"
	"fmt"
	"testing"

	"github.com/maakh3/api-petstore-service-layer/mocks"
	"github.com/maakh3/api-petstore-service-layer/models"
	"go.uber.org/mock/gomock"
)

func TestPetService_AddPet(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockRepo := mocks.NewMockPetRepositoryInterface(ctrl)

		input := models.Pet{Id: 1, Name: "fido", Status: "available", Tags: []models.Tag{{Id: 1, Name: "small"}}}
		want := models.Pet{Id: 1, Name: "fido", Status: "available", Tags: []models.Tag{{Id: 1, Name: "small"}}}
		mockRepo.EXPECT().AddPet(input).Return(want, nil) // repository add pet call, mocked

		service := NewPetService(mockRepo)
		got, err := service.AddPet(input) // service add pet call, the unit under test
		if err != nil {
			t.Fatalf("AddPet() unexpected error: %v", err)
		}
		if got.Id != input.Id {
			t.Fatalf("AddPet() Id = %d, want 1", got.Id)
		}
		if got.Name != input.Name {
			t.Fatalf("AddPet() Name = %q, want %q", got.Name, input.Name)
		}
		if got.Status != input.Status {
			t.Fatalf("AddPet() Status = %q, want %q", got.Status, input.Status)
		}
		if len(got.Tags) != len(input.Tags) {
			t.Fatalf("AddPet() Tags len = %d, want %d", len(got.Tags), len(input.Tags))
		}
	})
}

func TestPetService_UpdatePet(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockRepo := mocks.NewMockPetRepositoryInterface(ctrl)
		input := models.Pet{Name: "fido-updated", Status: "sold"}
		want := models.Pet{Id: 0, Name: "fido-updated", Status: "sold"}
		mockRepo.EXPECT().UpdatePet(input).Return(want, nil)

		service := NewPetService(mockRepo)

		got, err := service.UpdatePet(input)
		if err != nil {
			t.Fatalf("UpdatePet() unexpected error: %v", err)
		}

		if got.Id != want.Id {
			t.Fatalf("UpdatePet() Id = %d, want %d", got.Id, want.Id)
		}
		if got.Name != "fido-updated" {
			t.Fatalf("UpdatePet() Name = %q, want %q", got.Name, "fido-updated")
		}
		if got.Status != "sold" {
			t.Fatalf("UpdatePet() Status = %q, want %q", got.Status, "sold")
		}
	})
	t.Run("not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockRepo := mocks.NewMockPetRepositoryInterface(ctrl)
		input := models.Pet{Id: 999, Name: "ghost", Status: "available"}
		mockRepo.EXPECT().UpdatePet(input).Return(models.Pet{}, fmt.Errorf("pet with Id 999 not found"))

		service := NewPetService(mockRepo)
		_, err := service.UpdatePet(models.Pet{Id: 999, Name: "ghost", Status: "available"})
		if err == nil {
			t.Fatal("UpdatePet() error = nil, want ErrPetNotFound")
		}
		if !errors.Is(err, ErrPetNotFound) {
			t.Fatalf("UpdatePet() error = %v, want %v", err, ErrPetNotFound)
		}
	})
}

func TestPetService_FindPetsByStatus(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockRepo := mocks.NewMockPetRepositoryInterface(ctrl)

		want := []models.Pet{
			{Id: 1, Name: "fido", Status: "available"},
			{Id: 3, Name: "rex", Status: "available"},
		}
		mockRepo.EXPECT().FindPetsByStatus("available").Return(want, nil)

		service := NewPetService(mockRepo)
		got, err := service.FindPetsByStatus("available")
		if err != nil {
			t.Fatalf("FindPetsByStatus() unexpected error: %v", err)
		}

		if len(got) != len(want) {
			t.Fatalf("FindPetsByStatus() len = %d, want %d", len(got), len(want))
		}
		for _, pet := range got {
			if pet.Status != "available" {
				t.Fatalf("FindPetsByStatus() returned status %q, want %q", pet.Status, "available")
			}
		}
	})
	t.Run("no match", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockRepo := mocks.NewMockPetRepositoryInterface(ctrl)

		mockRepo.EXPECT().FindPetsByStatus("pending").Return([]models.Pet{}, nil)

		service := NewPetService(mockRepo)
		got, err := service.FindPetsByStatus("pending")
		if err != nil {
			t.Fatalf("FindPetsByStatus() unexpected error: %v", err)
		}

		if got == nil {
			t.Fatal("FindPetsByStatus() returned nil slice, want empty slice")
		}
		if len(got) != 0 {
			t.Fatalf("FindPetsByStatus() len = %d, want 0", len(got))
		}
	})
	t.Run("repository error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockRepo := mocks.NewMockPetRepositoryInterface(ctrl)

		wantErr := errors.New("db unavailable")
		mockRepo.EXPECT().FindPetsByStatus("available").Return(nil, wantErr)

		service := NewPetService(mockRepo)
		_, err := service.FindPetsByStatus("available")
		if err == nil {
			t.Fatal("FindPetsByStatus() error = nil, want propagated error")
		}
		if !errors.Is(err, wantErr) {
			t.Fatalf("FindPetsByStatus() error = %v, want %v", err, wantErr)
		}
	})
}

func TestPetService_FindPetsByTags(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockRepo := mocks.NewMockPetRepositoryInterface(ctrl)

		searchTags := []models.Tag{{Id: 1, Name: "small"}}
		want := []models.Pet{
			{Id: 1, Name: "fido", Status: "available", Tags: []models.Tag{{Id: 1, Name: "small"}, {Id: 2, Name: "brown"}}},
			{Id: 2, Name: "mittens", Status: "sold", Tags: []models.Tag{{Id: 3, Name: "small"}, {Id: 4, Name: "black"}}},
		}
		mockRepo.EXPECT().FindPetsByTags(searchTags).Return(want, nil)

		service := NewPetService(mockRepo)
		got, err := service.FindPetsByTags(searchTags)
		if err != nil {
			t.Fatalf("FindPetsByTags() unexpected error: %v", err)
		}

		if len(got) != len(want) {
			t.Fatalf("FindPetsByTags() len = %d, want %d", len(got), len(want))
		}
		for _, pet := range got {
			hasSmallTag := false
			for _, tag := range pet.Tags {
				if tag.Name == "small" {
					hasSmallTag = true
					break
				}
			}
			if !hasSmallTag {
				t.Fatalf("FindPetsByTags() returned pet without 'small' tag")
			}
		}
	})
	t.Run("no match", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockRepo := mocks.NewMockPetRepositoryInterface(ctrl)

		searchTags := []models.Tag{{Id: 999, Name: "nonexistent"}}
		mockRepo.EXPECT().FindPetsByTags(searchTags).Return([]models.Pet{}, nil)

		service := NewPetService(mockRepo)
		got, err := service.FindPetsByTags(searchTags)
		if err != nil {
			t.Fatalf("FindPetsByTags() unexpected error: %v", err)
		}

		if got == nil {
			t.Fatal("FindPetsByTags() returned nil slice, want empty slice")
		}
		if len(got) != 0 {
			t.Fatalf("FindPetsByTags() len = %d, want 0", len(got))
		}
	})
	t.Run("repository error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockRepo := mocks.NewMockPetRepositoryInterface(ctrl)

		searchTags := []models.Tag{{Id: 1, Name: "small"}}
		wantErr := errors.New("db error")
		mockRepo.EXPECT().FindPetsByTags(searchTags).Return(nil, wantErr)

		service := NewPetService(mockRepo)
		_, err := service.FindPetsByTags(searchTags)
		if err == nil {
			t.Fatal("FindPetsByTags() error = nil, want propagated error")
		}
		if !errors.Is(err, wantErr) {
			t.Fatalf("FindPetsByTags() error = %v, want %v", err, wantErr)
		}
	})
}
