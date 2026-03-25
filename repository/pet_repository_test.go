package repository

import (
	"reflect"
	"strings"
	"testing"

	"github.com/maakh3/api-petstore-service-layer/models"
)

func TestRepositoryAddPet(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
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
	})
}

func TestRepositoryUpdatePet(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
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
	})
	t.Run("not found", func(t *testing.T) {
		repo := NewPetRepository()

		_, err := repo.UpdatePet(models.Pet{Id: 777, Name: "ghost"})
		if err == nil {
			t.Fatal("UpdatePet() error = nil, want not found error")
		}
		if !strings.Contains(err.Error(), "not found") || !strings.Contains(err.Error(), "Id") {
			t.Fatalf("UpdatePet() error = %q, want substring match for not found + Id", err.Error())
		}
	})
}

func TestRepositoryFindPetsByStatus(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		repo := NewPetRepository()

		a, err := repo.AddPet(models.Pet{Name: "a", Status: "available"})
		if err != nil {
			t.Fatalf("AddPet() setup unexpected error: %v", err)
		}
		_, err = repo.AddPet(models.Pet{Name: "b", Status: "pending"})
		if err != nil {
			t.Fatalf("AddPet() setup unexpected error: %v", err)
		}
		c, err := repo.AddPet(models.Pet{Name: "c", Status: "available"})
		if err != nil {
			t.Fatalf("AddPet() setup unexpected error: %v", err)
		}

		got, err := repo.FindPetsByStatus("available")
		if err != nil {
			t.Fatalf("FindPetsByStatus() unexpected error: %v", err)
		}

		wantIDs := map[int]struct{}{a.Id: {}, c.Id: {}}
		if len(got) != len(wantIDs) {
			t.Fatalf("FindPetsByStatus() len = %d, want %d", len(got), len(wantIDs))
		}

		gotIDs := make(map[int]struct{}, len(got))
		for _, pet := range got {
			if pet.Status != "available" {
				t.Fatalf("FindPetsByStatus() returned status %q, want %q", pet.Status, "available")
			}
			gotIDs[pet.Id] = struct{}{}
		}

		for id := range wantIDs {
			if _, ok := gotIDs[id]; !ok {
				t.Fatalf("FindPetsByStatus() missing Id %d", id)
			}
		}
	})
	t.Run("no match", func(t *testing.T) {
		repo := NewPetRepository()

		_, err := repo.AddPet(models.Pet{Name: "a", Status: "sold"})
		if err != nil {
			t.Fatalf("AddPet() setup unexpected error: %v", err)
		}

		got, err := repo.FindPetsByStatus("pending")
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

func TestRepositoryFindPetsByTags(t *testing.T) {
	t.Run("full match", func(t *testing.T) {
		repo := NewPetRepository()
		a, _ := repo.AddPet(testPet("a", "available", "friendly", "small"))
		b, _ := repo.AddPet(testPet("b", "available", "friendly"))
		_, _ = repo.AddPet(testPet("c", "available", "aggressive"))

		got, err := repo.FindPetsByTags([]models.Tag{{Name: "friendly"}})
		if err != nil {
			t.Fatalf("FindPetsByTags() unexpected error: %v", err)
		}

		want := map[int]struct{}{a.Id: {}, b.Id: {}}
		gotSet := toIdSet(got)
		if len(gotSet) != len(want) {
			t.Fatalf("FindPetsByTags() len = %d, want %d", len(gotSet), len(want))
		}
		for id := range want {
			if _, ok := gotSet[id]; !ok {
				t.Fatalf("FindPetsByTags() missing Id %d", id)
			}
		}
	})

	t.Run("partial miss", func(t *testing.T) {
		repo := NewPetRepository()
		repo.AddPet(testPet("a", "available", "friendly", "small"))

		got, err := repo.FindPetsByTags([]models.Tag{{Name: "friendly"}, {Name: "small"}, {Name: "missing"}})
		if err != nil {
			t.Fatalf("FindPetsByTags() unexpected error: %v", err)
		}
		if len(got) != 0 {
			t.Fatalf("FindPetsByTags() len = %d, want 0", len(got))
		}
	})

	t.Run("empty search tags is match all", func(t *testing.T) {
		repo := NewPetRepository()
		a, _ := repo.AddPet(testPet("a", "available", "friendly", "small"))
		b, _ := repo.AddPet(testPet("b", "available", "friendly"))
		c, _ := repo.AddPet(testPet("c", "available", "aggressive"))

		got, err := repo.FindPetsByTags(nil)
		if err != nil {
			t.Fatalf("FindPetsByTags() unexpected error: %v", err)
		}

		want := map[int]struct{}{a.Id: {}, b.Id: {}, c.Id: {}}
		gotSet := toIdSet(got)
		if len(gotSet) != len(want) {
			t.Fatalf("FindPetsByTags() len = %d, want %d", len(gotSet), len(want))
		}
		for id := range want {
			if _, ok := gotSet[id]; !ok {
				t.Fatalf("FindPetsByTags() missing Id %d", id)
			}
		}
	})
}

func testPet(name, status string, tags ...string) models.Pet {
	petTags := make([]models.Tag, 0, len(tags))
	for i, tag := range tags {
		petTags = append(petTags, models.Tag{Id: i + 1, Name: tag})
	}

	return models.Pet{
		Name:   name,
		Status: status,
		Tags:   petTags,
	}
}

func toIdSet(pets []models.Pet) map[int]struct{} {
	out := make(map[int]struct{}, len(pets))
	for _, pet := range pets {
		out[pet.Id] = struct{}{}
	}
	return out
}
