package repository

import (
	"api-petstore-service-layer/models"
	"reflect"
	"strings"
	"testing"
)

func TestRepositoryAddPet(t *testing.T) {
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

func TestRepositoryUpdatePet(t *testing.T) {
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

func TestRepositoryUpdatePet_NotFound(t *testing.T) {
	repo := NewPetRepository()

	_, err := repo.UpdatePet(models.Pet{Id: 777, Name: "ghost"})
	if err == nil {
		t.Fatal("UpdatePet() error = nil, want not found error")
	}
	if !strings.Contains(err.Error(), "not found") || !strings.Contains(err.Error(), "Id") {
		t.Fatalf("UpdatePet() error = %q, want substring match for not found + Id", err.Error())
	}
}

func TestRepositoryFindPetsByStatus(t *testing.T) {
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
}

func TestRepositoryFindPetsByStatus_NoMatch(t *testing.T) {
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
}

//func TestGetById(t *testing.T) {
//	repo := NewPetRepository()
//	created, _ := repo.AddPet(testPet("fido", "available"))
//
//	t.Run("found", func(t *testing.T) {
//		got, err := repo.GetById(int64(created.Id))
//		if err != nil {
//			t.Fatalf("GetById() unexpected error: %v", err)
//		}
//		if got.Id != created.Id {
//			t.Fatalf("GetById() Id = %d, want %d", got.Id, created.Id)
//		}
//	})
//
//	t.Run("not found", func(t *testing.T) {
//		_, err := repo.GetById(999)
//		if err == nil {
//			t.Fatal("GetById() error = nil, want not found error")
//		}
//		if !strings.Contains(err.Error(), "not found") || !strings.Contains(err.Error(), "Id") {
//			t.Fatalf("GetById() error = %q, want substring match for not found + Id", err.Error())
//		}
//	})
//}

//	func TestDelete(t *testing.T) {
//		repo := NewPetRepository()
//		created, _ := repo.AddPet(testPet("fido", "available"))
//
//		t.Run("deletes existing pet", func(t *testing.T) {
//			err := repo.Delete(int64(created.Id))
//			if err != nil {
//				t.Fatalf("Delete() unexpected error: %v", err)
//			}
//
//			_, err = repo.GetById(int64(created.Id))
//			if err == nil {
//				t.Fatal("expected deleted pet to be missing")
//			}
//		})
//
//		t.Run("not found", func(t *testing.T) {
//			err := repo.Delete(1234)
//			if err == nil {
//				t.Fatal("Delete() error = nil, want not found error")
//			}
//			if !strings.Contains(err.Error(), "not found") || !strings.Contains(err.Error(), "Id") {
//				t.Fatalf("Delete() error = %q, want substring match for not found + Id", err.Error())
//			}
//		})
//	}
//
//
//	func TestRepositoryFindPetsByTags(t *testing.T) {
//		repo := NewPetRepository()
//		a, _ := repo.AddPet(testPet("a", "available", "friendly", "small"))
//		b, _ := repo.AddPet(testPet("b", "available", "friendly"))
//		c, _ := repo.AddPet(testPet("c", "available", "aggressive"))
//
//		t.Run("full match", func(t *testing.T) {
//			got, err := repo.FindPetsByTags([]models.Tag{{Name: "friendly"}})
//			if err != nil {
//				t.Fatalf("FindPetsByTags() unexpected error: %v", err)
//			}
//
//			want := map[int]struct{}{a.Id: {}, b.Id: {}}
//			gotSet := toIdSet(got)
//			if len(gotSet) != len(want) {
//				t.Fatalf("FindPetsByTags() len = %d, want %d", len(gotSet), len(want))
//			}
//			for id := range want {
//				if _, ok := gotSet[id]; !ok {
//					t.Fatalf("FindPetsByTags() missing Id %d", id)
//				}
//			}
//		})
//
//		t.Run("partial miss", func(t *testing.T) {
//			got, err := repo.FindPetsByTags([]models.Tag{{Name: "friendly"}, {Name: "small"}, {Name: "missing"}})
//			if err != nil {
//				t.Fatalf("FindPetsByTags() unexpected error: %v", err)
//			}
//			if len(got) != 0 {
//				t.Fatalf("FindPetsByTags() len = %d, want 0", len(got))
//			}
//		})
//
//		t.Run("empty search tags is match all", func(t *testing.T) {
//			got, err := repo.FindPetsByTags(nil)
//			if err != nil {
//				t.Fatalf("FindPetsByTags() unexpected error: %v", err)
//			}
//
//			want := map[int]struct{}{a.Id: {}, b.Id: {}, c.Id: {}}
//			gotSet := toIdSet(got)
//			if len(gotSet) != len(want) {
//				t.Fatalf("FindPetsByTags() len = %d, want %d", len(gotSet), len(want))
//			}
//			for id := range want {
//				if _, ok := gotSet[id]; !ok {
//					t.Fatalf("FindPetsByTags() missing Id %d", id)
//				}
//			}
//		})
//	}
//
//	func TestUploadImage(t *testing.T) {
//		repo := NewPetRepository()
//		created, _ := repo.AddPet(testPet("fido", "available"))
//
//		t.Run("appends image url", func(t *testing.T) {
//			updated, err := repo.UploadImage(int64(created.Id), "https://img/1.png")
//			if err != nil {
//				t.Fatalf("UploadImage() unexpected error: %v", err)
//			}
//			if len(updated.PhotoUrls) != 1 || updated.PhotoUrls[0] != "https://img/1.png" {
//				t.Fatalf("UploadImage() got PhotoUrls=%v, want [https://img/1.png]", updated.PhotoUrls)
//			}
//		})
//
//		t.Run("not found", func(t *testing.T) {
//			_, err := repo.UploadImage(9999, "https://img/missing.png")
//			if err == nil {
//				t.Fatal("UploadImage() error = nil, want not found error")
//			}
//			if !strings.Contains(err.Error(), "not found") || !strings.Contains(err.Error(), "Id") {
//				t.Fatalf("UploadImage() error = %q, want substring match for not found + Id", err.Error())
//			}
//		})
//	}
//
//	func TestContainsAllTags(t *testing.T) {
//		tests := []struct {
//			name       string
//			petTags    []models.Tag
//			searchTags []models.Tag
//			want       bool
//		}{
//			{
//				name:       "empty search tags matches all",
//				petTags:    []models.Tag{{Name: "a"}},
//				searchTags: nil,
//				want:       true,
//			},
//			{
//				name:       "empty pet tags with non-empty search fails",
//				petTags:    nil,
//				searchTags: []models.Tag{{Name: "a"}},
//				want:       false,
//			},
//			{
//				name:       "all search tags present",
//				petTags:    []models.Tag{{Name: "a"}, {Name: "b"}},
//				searchTags: []models.Tag{{Name: "a"}},
//				want:       true,
//			},
//			{
//				name:       "missing search tag",
//				petTags:    []models.Tag{{Name: "a"}},
//				searchTags: []models.Tag{{Name: "a"}, {Name: "b"}},
//				want:       false,
//			},
//			{
//				name:       "duplicate search tags still true when present",
//				petTags:    []models.Tag{{Name: "a"}},
//				searchTags: []models.Tag{{Name: "a"}, {Name: "a"}},
//				want:       true,
//			},
//		}
//
//		for _, tc := range tests {
//			t.Run(tc.name, func(t *testing.T) {
//				got := containsAllTags(tc.petTags, tc.searchTags)
//				if got != tc.want {
//					t.Fatalf("containsAllTags() = %v, want %v", got, tc.want)
//				}
//			})
//		}
//	}
//func testPet(name, status string, tags ...string) models.Pet {
//	petTags := make([]models.Tag, 0, len(tags))
//	for i, tag := range tags {
//		petTags = append(petTags, models.Tag{Id: i + 1, Name: tag})
//	}
//
//	return models.Pet{
//		Name:   name,
//		Status: status,
//		Tags:   petTags,
//	}
//}

//
//func toIdSet(pets []models.Pet) map[int]struct{} {
//	// this function converts a slice of Pet objects into a set of their Ids for easier comparison in tests
//	out := make(map[int]struct{}, len(pets))
//	for _, pet := range pets {
//		out[pet.Id] = struct{}{}
//	}
//	return out
//}
