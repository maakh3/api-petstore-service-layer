package repository

import (
	"fmt"
	"reflect"
	"strings"

	//"reflect"
	//"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/maakh3/api-petstore-service-layer/models"
)

func TestRepositoryAddPet(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create sqlmock: %v", err)
		}
		defer db.Close()

		repo := NewPetRepository(db)

		input := models.Pet{
			Name:   "fido",
			Status: "available",
			Tags: []models.Tag{
				{Id: 1, Name: "small"},
				{Id: 2, Name: "friendly"},
			},
		}

		// expect the INSERT and return id = 1
		mock.ExpectQuery(`INSERT INTO pets`).
			WithArgs(
				input.Id,
				input.Name,
				sqlmock.AnyArg(), // category JSON
				sqlmock.AnyArg(), // photo_urls array
				sqlmock.AnyArg(), // tags JSON
				input.Status,
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		got, err := repo.AddPet(input)
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

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unfulfilled sql expectations: %v", err)
		}
	})

	t.Run("db error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create sqlmock: %v", err)
		}
		defer db.Close()

		repo := NewPetRepository(db)

		mock.ExpectQuery(`INSERT INTO pets`).
			WillReturnError(fmt.Errorf("db connection error"))

		_, err = repo.AddPet(models.Pet{Name: "fido", Status: "available"})
		if err == nil {
			t.Fatal("AddPet() error = nil, want db error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unfulfilled sql expectations: %v", err)
		}
	})
}

func TestRepositoryUpdatePet(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create sqlmock: %v", err)
		}
		defer db.Close()

		repo := NewPetRepository(db)

		input := models.Pet{
			Id:   1,
			Name: "fido-updated",
			Tags: []models.Tag{
				{Id: 1, Name: "small"},
			},
			Status: "sold",
		}

		// expect the UPDATE to succeed and affect 1 row
		mock.ExpectExec(`UPDATE pets`).
			WithArgs(
				input.Id,
				input.Name,
				input.Status,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnResult(sqlmock.NewResult(0, 1)) // last insert id, 1 row affected

		got, err := repo.UpdatePet(input)
		if err != nil {
			t.Fatalf("UpdatePet() unexpected error: %v", err)
		}

		// UpdatePet returns the input pet object, so we can check that it was returned unchanged
		if !reflect.DeepEqual(got, input) {
			t.Fatalf("UpdatePet() returned %#v, want %#v", got, input)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unfulfilled sql expectations: %v", err)
		}
	})
	t.Run("not found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create sqlmock: %v", err)
		}
		defer db.Close()

		repo := NewPetRepository(db)

		// Expect the UPDATE but affect 0 rows (no pet with ID 777)
		mock.ExpectExec(`UPDATE pets`).
			WithArgs(
				777,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected

		_, err = repo.UpdatePet(models.Pet{Id: 777, Name: "ghost", Status: "available"})
		if err == nil {
			t.Fatal("UpdatePet() error = nil, want not found error")
		}
		if !strings.Contains(err.Error(), "not found") || !strings.Contains(err.Error(), "Id") {
			t.Fatalf("UpdatePet() error = %q, want substring match for 'not found' and 'Id'", err.Error())
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unfulfilled sql expectations: %v", err)
		}
	})

	t.Run("db exec error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create sqlmock: %v", err)
		}
		defer db.Close()

		repo := NewPetRepository(db)

		// Simulate a database error
		mock.ExpectExec(`UPDATE pets`).
			WillReturnError(fmt.Errorf("connection lost"))

		_, err = repo.UpdatePet(models.Pet{Id: 1, Name: "test", Status: "available"})
		if err == nil {
			t.Fatal("UpdatePet() error = nil, want db error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unfulfilled sql expectations: %v", err)
		}
	})
}
