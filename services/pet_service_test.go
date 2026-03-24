package services

import (
	"testing"

	"github.com/maakh3/api-petstore-service-layer/models"
	"github.com/maakh3/api-petstore-service-layer/repository"
)

func TestPetService_AddPet(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		repo := repository.NewPetRepository()
		service := NewPetService(repo)

		input := models.Pet{
			Name:   "fido",
			Status: "available",
			Tags: []models.Tag{
				{Id: 1, Name: "small"},
			},
		}

		got, err := service.AddPet(input)
		if err != nil {
			t.Fatalf("AddPet() unexpected error: %v", err)
		}

		if got.Id != 1 {
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
		repo := repository.NewPetRepository()
		service := NewPetService(repo)

		created, err := repo.AddPet(models.Pet{Name: "fido", Status: "available"})
		if err != nil {
			t.Fatalf("AddPet() setup unexpected error: %v", err)
		}

		input := models.Pet{Id: created.Id, Name: "fido-updated", Status: "sold"}
		got, err := service.UpdatePet(input)
		if err != nil {
			t.Fatalf("UpdatePet() unexpected error: %v", err)
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
	})
	t.Run("not found", func(t *testing.T) {
		repo := repository.NewPetRepository()
		service := NewPetService(repo)

		_, err := service.UpdatePet(models.Pet{Id: 999, Name: "ghost", Status: "available"})
		if err == nil {
			t.Fatal("UpdatePet() error = nil, want ErrPetNotFound")
		}
		if err != ErrPetNotFound {
			t.Fatalf("UpdatePet() error = %v, want %v", err, ErrPetNotFound)
		}
	})
}

func TestPetService_FindPetsByStatus(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		repo := repository.NewPetRepository()
		service := NewPetService(repo)

		_, err := repo.AddPet(models.Pet{Name: "fido", Status: "available"})
		if err != nil {
			t.Fatalf("AddPet() setup unexpected error: %v", err)
		}
		_, err = repo.AddPet(models.Pet{Name: "mittens", Status: "sold"})
		if err != nil {
			t.Fatalf("AddPet() setup unexpected error: %v", err)
		}
		_, err = repo.AddPet(models.Pet{Name: "rex", Status: "available"})
		if err != nil {
			t.Fatalf("AddPet() setup unexpected error: %v", err)
		}

		got, err := service.FindPetsByStatus("available")
		if err != nil {
			t.Fatalf("FindPetsByStatus() unexpected error: %v", err)
		}

		if len(got) != 2 {
			t.Fatalf("FindPetsByStatus() len = %d, want 2", len(got))
		}
		for _, pet := range got {
			if pet.Status != "available" {
				t.Fatalf("FindPetsByStatus() returned status %q, want %q", pet.Status, "available")
			}
		}
	})
	t.Run("No match", func(t *testing.T) {
		repo := repository.NewPetRepository()
		service := NewPetService(repo)

		_, err := repo.AddPet(models.Pet{Name: "fido", Status: "sold"})
		if err != nil {
			t.Fatalf("AddPet() setup unexpected error: %v", err)
		}

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
}

func TestPetService_FindPetsByTags(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		repo := repository.NewPetRepository()
		service := NewPetService(repo)

		_, err := repo.AddPet(models.Pet{Name: "fido", Status: "available", Tags: []models.Tag{{Id: 1, Name: "small"}, {Id: 2, Name: "brown"}}})
		if err != nil {
			t.Fatalf("AddPet() setup unexpected error: %v", err)
		}
		_, err = repo.AddPet(models.Pet{Name: "mittens", Status: "sold", Tags: []models.Tag{{Id: 3, Name: "small"}, {Id: 4, Name: "black"}}})
		if err != nil {
			t.Fatalf("AddPet() setup unexpected error: %v", err)
		}
		_, err = repo.AddPet(models.Pet{Name: "rex", Status: "available", Tags: []models.Tag{{Id: 5, Name: "large"}, {Id: 6, Name: "green"}}})
		if err != nil {
			t.Fatalf("AddPet() setup unexpected error: %v", err)
		}

		got, err := service.FindPetsByTags([]models.Tag{{Id: 1, Name: "small"}})
		if err != nil {
			t.Fatalf("FindPetsByTags() unexpected error: %v", err)
		}

		if len(got) != 2 {
			t.Fatalf("FindPetsByTags() len = %d, want 2", len(got))
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
	t.Run("No match", func(t *testing.T) {
		repo := repository.NewPetRepository()
		service := NewPetService(repo)

		_, err := repo.AddPet(models.Pet{Name: "fido", Status: "available", Tags: []models.Tag{{Id: 1, Name: "small"}, {Id: 2, Name: "brown"}}})
		if err != nil {
			t.Fatalf("AddPet() setup unexpected error: %v", err)
		}

		got, err := service.FindPetsByTags([]models.Tag{{Id: 999, Name: "nonexistent"}})
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
}
