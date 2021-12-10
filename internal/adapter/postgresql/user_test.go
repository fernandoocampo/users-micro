package postgresql_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fernandoocampo/users-micro/internal/adapter/postgresql"
	"github.com/fernandoocampo/users-micro/internal/adapter/repository"
	"github.com/stretchr/testify/assert"
)

func TestSaveUser(t *testing.T) {
	ctx := context.TODO()
	givenUser := repository.User{
		ID:        "123",
		City:      "Cali",
		FirstName: "Alonso",
		LastName:  "Ojeda",
		Skills:    []string{"painter"},
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO jobseeker").ExpectExec().
		WithArgs(
			givenUser.ID,
			givenUser.FirstName,
			givenUser.LastName,
			givenUser.City,
			givenUser.Skills,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	userRepository := postgresql.NewUserRepository(db)

	// WHEN
	saveError := userRepository.Save(ctx, givenUser)

	assert.NoError(t, saveError)
}

func TestSaveUserButUnexpectedError(t *testing.T) {
	ctx := context.TODO()
	givenUser := repository.User{
		ID:        "123",
		City:      "Cali",
		FirstName: "Alonso",
		LastName:  "Ojeda",
		Skills:    []string{"painter"},
	}
	expectedError := errors.New("user cannot be stored")
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO jobseeker").ExpectExec().
		WithArgs(
			givenUser.ID,
			givenUser.FirstName,
			givenUser.LastName,
			givenUser.City,
			givenUser.Skills,
		).
		WillReturnError(errors.New("unexpected error"))

	userRepository := postgresql.NewUserRepository(db)

	// WHEN
	saveError := userRepository.Save(ctx, givenUser)

	assert.Error(t, saveError)
	assert.Equal(t, expectedError, saveError)
}

func TestUpdateUser(t *testing.T) {
	ctx := context.TODO()
	givenUser := repository.User{
		ID:        "123",
		City:      "Cali",
		FirstName: "Alonso",
		LastName:  "Ojeda",
		Skills:    []string{"painter"},
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare("UPDATE jobseeker").ExpectExec().
		WithArgs(
			givenUser.FirstName,
			givenUser.LastName,
			givenUser.City,
			givenUser.Skills,
			givenUser.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	userRepository := postgresql.NewUserRepository(db)

	// WHEN
	saveError := userRepository.Update(ctx, givenUser)

	assert.NoError(t, saveError)
}

func TestFindUserByID(t *testing.T) {
	ctx := context.TODO()
	givenUserID := "123"
	expectedUser := repository.User{
		ID:        "123",
		City:      "Cali",
		FirstName: "Alonso",
		LastName:  "Ojeda",
		Skills:    []string{"painter"},
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "firstname", "lastname", "city", "skills"}).
		AddRow("123", "Alonso", "Ojeda", "Cali", []byte(`["painter"]`))

	mock.ExpectQuery("SELECT (.+) FROM jobseeker").
		WillReturnRows(rows)

	userRepository := postgresql.NewUserRepository(db)

	// WHEN
	got, saveError := userRepository.FindByID(ctx, givenUserID)

	assert.NoError(t, saveError)
	assert.Equal(t, &expectedUser, got)
}

func TestFindUserByIDButError(t *testing.T) {
	ctx := context.TODO()
	givenUserID := "123"
	expectedError := errors.New("user cannot be read in the database")
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT (.+) FROM jobseeker").
		WillReturnError(errors.New("error"))

	userRepository := postgresql.NewUserRepository(db)

	// WHEN
	got, saveError := userRepository.FindByID(ctx, givenUserID)

	assert.Error(t, saveError)
	assert.Nil(t, got)
	assert.Equal(t, expectedError, saveError)
}

func TestFindUsersByCity(t *testing.T) {
	ctx := context.TODO()
	givenFilter := repository.UserFilter{
		City:        "Cali",
		Page:        1,
		RowsPerPage: 10,
	}
	expectedResult := repository.FindUsersResult{
		Users: []repository.User{
			{
				ID:        "123",
				City:      "Cali",
				FirstName: "Alonso",
				LastName:  "Ojeda",
				Skills:    []string{"painter"},
			},
			{
				ID:        "124",
				City:      "Cali",
				FirstName: "Alicia",
				LastName:  "Cifuentes",
				Skills:    []string{"sculptor"},
			},
		},
		Total:       2,
		Page:        1,
		RowsPerPage: 10,
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	countRow := sqlmock.NewRows([]string{"COUNT(*)"}).
		AddRow("2")

	mock.ExpectPrepare("SELECT (.+) FROM jobseeker WHERE city").
		ExpectQuery().WithArgs("Cali").
		WillReturnRows(countRow)

	rows := sqlmock.NewRows([]string{"id", "firstname", "lastname", "city", "skills"}).
		AddRow("123", "Alonso", "Ojeda", "Cali", []byte(`["painter"]`)).
		AddRow("124", "Alicia", "Cifuentes", "Cali", []byte(`["sculptor"]`))

	mock.ExpectQuery("SELECT (.+) FROM jobseeker WHERE city").
		WithArgs("Cali", 10, 0).
		WillReturnRows(rows)

	userRepository := postgresql.NewUserRepository(db)

	// WHEN
	got, findError := userRepository.SearchWithFilters(ctx, givenFilter)

	assert.NoError(t, findError)
	assert.Equal(t, expectedResult, got)
}

func TestFindUsersBySkill(t *testing.T) {
	ctx := context.TODO()
	givenFilter := repository.UserFilter{
		Skills:      []string{"cabinetmaker"},
		Page:        1,
		RowsPerPage: 10,
	}
	expectedResult := repository.FindUsersResult{
		Users: []repository.User{
			{
				ID:        "125",
				City:      "Bogota",
				FirstName: "Cecilia",
				LastName:  "Quiroga",
				Skills:    []string{"cabinetmaker"},
			},
			{
				ID:        "126",
				City:      "Medellin",
				FirstName: "Armando",
				LastName:  "Lopez",
				Skills:    []string{"sculptor", "cabinetmaker"},
			},
		},
		Total:       2,
		Page:        1,
		RowsPerPage: 10,
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	countRow := sqlmock.NewRows([]string{"COUNT(*)"}).
		AddRow("2")

	mock.ExpectPrepare("SELECT (.+) FROM jobseeker WHERE skills").
		ExpectQuery().WithArgs([]byte(`["cabinetmaker"]`)).
		WillReturnRows(countRow)

	rows := sqlmock.NewRows([]string{"id", "firstname", "lastname", "city", "skills"}).
		AddRow("125", "Cecilia", "Quiroga", "Bogota", []byte(`["cabinetmaker"]`)).
		AddRow("126", "Armando", "Lopez", "Medellin", []byte(`["sculptor", "cabinetmaker"]`))

	mock.ExpectQuery("SELECT (.+) FROM jobseeker WHERE skills").
		WithArgs([]byte(`["cabinetmaker"]`), 10, 0).
		WillReturnRows(rows)

	userRepository := postgresql.NewUserRepository(db)

	// WHEN
	got, findError := userRepository.SearchWithFilters(ctx, givenFilter)

	assert.NoError(t, findError)
	assert.Equal(t, expectedResult, got)
}

func TestFindUsersBySkillAndCity(t *testing.T) {
	ctx := context.TODO()
	givenFilter := repository.UserFilter{
		City:        "Medellin",
		Skills:      []string{"cabinetmaker"},
		Page:        1,
		RowsPerPage: 10,
	}
	expectedResult := repository.FindUsersResult{
		Users: []repository.User{
			{
				ID:        "126",
				City:      "Medellin",
				FirstName: "Armando",
				LastName:  "Lopez",
				Skills:    []string{"sculptor", "cabinetmaker"},
			},
		},
		Total:       1,
		Page:        1,
		RowsPerPage: 10,
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	countRow := sqlmock.NewRows([]string{"COUNT(*)"}).
		AddRow("1")

	mock.ExpectPrepare("SELECT (.+) FROM jobseeker WHERE city").
		ExpectQuery().WithArgs("Medellin", []byte(`["cabinetmaker"]`)).
		WillReturnRows(countRow)

	rows := sqlmock.NewRows([]string{"id", "firstname", "lastname", "city", "skills"}).
		AddRow("126", "Armando", "Lopez", "Medellin", []byte(`["sculptor", "cabinetmaker"]`))

	mock.ExpectQuery("SELECT (.+) FROM jobseeker WHERE city").
		WithArgs("Medellin", []byte(`["cabinetmaker"]`), 10, 0).
		WillReturnRows(rows)

	userRepository := postgresql.NewUserRepository(db)

	// WHEN
	got, findError := userRepository.SearchWithFilters(ctx, givenFilter)

	assert.NoError(t, findError)
	assert.Equal(t, expectedResult, got)
}
