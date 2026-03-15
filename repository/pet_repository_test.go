package repository

import (
	"api-petstore-service-layer/models"
	"reflect"
	"strings"
	"testing"
)

func TestAddPet(t *testing.T) {
	repo := NewPetRepository()
	inputPet := models.Pet{
		Name:   "fido",
		Status: "available",
		Tags: []models.Tag{
			{Id: 1, Name: "small"},
			{Id: 2, Name: "friendly"},
		},
	}

	want := models.Pet{
		Id:     1,
		Name:   "fido",
		Status: "available",
		Tags: []models.Tag{
			{Id: 1, Name: "small"},
			{Id: 2, Name: "friendly"},
		},
	}

	got, err := repo.AddPet(inputPet)
	if err != nil {
		t.Fatalf("AddPet() unexpected error: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("AddPet() returned %#v, want %#v", got, want)
	}

	storedPet, exists := repo.pets[int64(want.Id)]
	if !exists {
		t.Fatalf("AddPet() did not store pet with Id %d", want.Id)
	}

	if !reflect.DeepEqual(*storedPet, want) {
		t.Fatalf("stored pet = %#v, want %#v", *storedPet, want)
	}

	if repo.nextId != 2 {
		t.Fatalf("nextId = %d, want 2", repo.nextId)
	}
}

func TestUpdatePet(t *testing.T) {
	repo := NewPetRepository()
	existingPet := models.Pet{
		Id:     1,
		Name:   "fido",
		Status: "available",
		Tags: []models.Tag{
			{Id: 1, Name: "small"},
		},
	}
	repo.pets[int64(existingPet.Id)] = &existingPet
	repo.nextId = 2

	inputPet := models.Pet{
		Id:     1,
		Name:   "fido-updated",
		Status: "sold",
		Tags: []models.Tag{
			{Id: 1, Name: "small"},
		},
	}

	want := models.Pet{
		Id:     1,
		Name:   "fido-updated",
		Status: "sold",
		Tags: []models.Tag{
			{Id: 1, Name: "small"},
		},
	}

	got, err := repo.UpdatePet(inputPet)
	if err != nil {
		t.Fatalf("UpdatePet() unexpected error: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("UpdatePet() returned %#v, want %#v", got, want)
	}

	storedPet, exists := repo.pets[int64(want.Id)]
	if !exists {
		t.Fatalf("UpdatePet() did not store pet with Id %d", want.Id)
	}

	if !reflect.DeepEqual(*storedPet, want) {
		t.Fatalf("stored pet = %#v, want %#v", *storedPet, want)
	}

	if repo.nextId != 2 {
		t.Fatalf("nextId = %d, want 2", repo.nextId)
	}
}

func TestUpdatePet_NotFound(t *testing.T) {
	repo := NewPetRepository()

	_, err := repo.UpdatePet(models.Pet{Id: 777, Name: "ghost"})
	if err == nil {
		t.Fatal("UpdatePet() error = nil, want not found error")
	}
	if !strings.Contains(err.Error(), "not found") || !strings.Contains(err.Error(), "Id") {
		t.Fatalf("UpdatePet() error = %q, want substring match for not found + Id", err.Error())
	}
}
