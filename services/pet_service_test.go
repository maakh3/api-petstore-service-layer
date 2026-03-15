package services

import (
	"api-petstore-service-layer/models"
	"api-petstore-service-layer/repository"
	"testing"
)

func TestPetServiceAddPet(t *testing.T) {
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
}

func TestPetServiceUpdatePet(t *testing.T) {
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
}

func TestPetServiceUpdatePet_NotFound(t *testing.T) {
	repo := repository.NewPetRepository()
	service := NewPetService(repo)

	_, err := service.UpdatePet(models.Pet{Id: 999, Name: "ghost", Status: "available"})
	if err == nil {
		t.Fatal("UpdatePet() error = nil, want ErrPetNotFound")
	}
	if err != ErrPetNotFound {
		t.Fatalf("UpdatePet() error = %v, want %v", err, ErrPetNotFound)
	}
}
