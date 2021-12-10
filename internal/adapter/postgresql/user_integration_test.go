package postgresql_test

import (
	"context"
	"testing"

	"github.com/fernandoocampo/users-micro/internal/adapter/postgresql"
	"github.com/fernandoocampo/users-micro/internal/adapter/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	if !*integration {
		t.Skip("this is an integration test, to execute this test send integration flag to true")
	}
	ctx := context.TODO()

	givenParameters := postgresql.Parameters{
		DBName:   "postgres",
		Host:     "localhost",
		User:     "postgres",
		Password: "postgres",
		Port:     5432,
	}

	client := createClient(t, givenParameters)
	defer client.Close()

	newUserID := uuid.New().String()
	newUser := repository.User{
		ID:        newUserID,
		FirstName: "users-micro",
		LastName:  "Wayne",
		City:      "Medellin",
		Skills:    []string{"work", "happy"},
	}
	userRepository := postgresql.NewUserRepository(client)

	saveErr := userRepository.Save(ctx, newUser)
	assert.NoError(t, saveErr)
	if saveErr != nil {
		t.Errorf("unexpected error trying to create a user: %s", saveErr)
		t.FailNow()
	}
}

func TestUpdateExistingUser(t *testing.T) {
	if !*integration {
		t.Skip("this is an integration test, to execute this test send integration flag to true")
	}
	ctx := context.TODO()

	givenParameters := postgresql.Parameters{
		DBName:   "postgres",
		Host:     "localhost",
		User:     "postgres",
		Password: "postgres",
		Port:     5432,
	}
	client := createClient(t, givenParameters)
	defer client.Close()

	newUserID := uuid.New().String()
	newUser := repository.User{
		ID:        newUserID,
		FirstName: "users-micro",
		LastName:  "Wayne",
		City:      "Medellin",
		Skills:    []string{"work", "happy"},
	}
	givenUser := repository.User{
		ID:        newUserID,
		City:      "Cali",
		FirstName: "Alonso",
		LastName:  "Ojeda",
		Skills:    []string{"painter"},
	}

	userRepository := postgresql.NewUserRepository(client)

	// WHEN
	saveErr := userRepository.Save(ctx, newUser)
	if saveErr != nil {
		t.Errorf("unexpected error trying to save an user in update user test: %s", saveErr)

	}
	saveError := userRepository.Update(ctx, givenUser)

	// THEN
	assert.NoError(t, saveError)
	storedUser, getErr := userRepository.FindByID(ctx, newUserID)
	if getErr != nil {
		t.Errorf("unexpected error trying to get a user by its id: %s", getErr)
		t.FailNow()
	}

	assert.Equal(t, &givenUser, storedUser)
}

func TestReadUser(t *testing.T) {
	if !*integration {
		t.Skip("this is an integration test, to execute this test send integration flag to true")
	}

	givenParameters := postgresql.Parameters{
		DBName:   "postgres",
		Host:     "localhost",
		User:     "postgres",
		Password: "postgres",
		Port:     5432,
	}

	client := createClient(t, givenParameters)
	defer client.Close()

	ctx := context.TODO()
	newUser := repository.User{
		ID:        "sfsfsf-sdfsf1234",
		FirstName: "users-micro",
		LastName:  "Wayne",
		City:      "Medellin",
		Skills:    []string{"work"},
	}
	userRepository := postgresql.NewUserRepository(client)

	saveErr := userRepository.Save(ctx, newUser)
	assert.NoError(t, saveErr)
	if saveErr != nil {
		t.Errorf("unexpected error trying to create a user: %s", saveErr)
		t.FailNow()
	}

	storedUser, getErr := userRepository.FindByID(ctx, "sfsfsf-sdfsf1234")
	assert.NoError(t, getErr)
	if getErr != nil {
		t.Errorf("unexpected error trying to get a user by its id: %s", getErr)
		t.FailNow()
	}

	assert.Equal(t, &newUser, storedUser)
}

func TestFindUsersByCityIntegration(t *testing.T) {
	if !*integration {
		t.Skip("this is an integration test, to execute this test send integration flag to true")
	}

	givenParameters := postgresql.Parameters{
		DBName:   "postgres",
		Host:     "localhost",
		User:     "postgres",
		Password: "postgres",
		Port:     5432,
	}

	client := createClient(t, givenParameters)
	defer client.Close()

	ctx := context.TODO()
	newUsers := []repository.User{
		{
			ID:        "123",
			City:      "Cali",
			FirstName: "Wayne",
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
		{
			ID:        "125",
			City:      "Bogota",
			FirstName: "Alicia",
			LastName:  "Cifuentes",
			Skills:    []string{"sculptor"},
		},
	}
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
				FirstName: "Wayne",
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
	userRepository := postgresql.NewUserRepository(client)

	for _, v := range newUsers {
		newUser := v
		saveErr := userRepository.Save(ctx, newUser)
		assert.NoError(t, saveErr)
		if saveErr != nil {
			t.Errorf("unexpected error trying to create a user: %s", saveErr)
			t.FailNow()
		}
	}

	searchResult, getErr := userRepository.SearchWithFilters(ctx, givenFilter)
	assert.NoError(t, getErr)
	if getErr != nil {
		t.Errorf("unexpected error trying to get a user by its id: %s", getErr)
		t.FailNow()
	}

	assert.Equal(t, expectedResult, searchResult)
}

func TestFindUsersBySkillsIntegration(t *testing.T) {
	if !*integration {
		t.Skip("this is an integration test, to execute this test send integration flag to true")
	}

	givenParameters := postgresql.Parameters{
		DBName:   "postgres",
		Host:     "localhost",
		User:     "postgres",
		Password: "postgres",
		Port:     5432,
	}

	client := createClient(t, givenParameters)
	defer client.Close()

	ctx := context.TODO()
	newUsers := []repository.User{
		{
			ID:        "126",
			City:      "Cali",
			FirstName: "Wayne",
			LastName:  "Ojeda",
			Skills:    []string{"painter", "decorator"},
		},
		{
			ID:        "127",
			City:      "Cali",
			FirstName: "Alicia",
			LastName:  "Cifuentes",
			Skills:    []string{"sculptor", "cabinetmaker", "painter"},
		},
		{
			ID:        "128",
			City:      "Bogota",
			FirstName: "Liliana",
			LastName:  "Marino",
			Skills:    []string{"painter", "cabinetmaker"},
		},
	}
	givenFilter := repository.UserFilter{
		Skills:      []string{"cabinetmaker", "painter"},
		Page:        1,
		RowsPerPage: 10,
	}
	expectedResult := repository.FindUsersResult{
		Users: []repository.User{
			{
				ID:        "127",
				City:      "Cali",
				FirstName: "Alicia",
				LastName:  "Cifuentes",
				Skills:    []string{"sculptor", "cabinetmaker", "painter"},
			},
			{
				ID:        "128",
				City:      "Bogota",
				FirstName: "Liliana",
				LastName:  "Marino",
				Skills:    []string{"painter", "cabinetmaker"},
			},
		},
		Total:       2,
		Page:        1,
		RowsPerPage: 10,
	}
	userRepository := postgresql.NewUserRepository(client)

	for _, v := range newUsers {
		newUser := v
		saveErr := userRepository.Save(ctx, newUser)
		assert.NoError(t, saveErr)
		if saveErr != nil {
			t.Errorf("unexpected error trying to create a user: %s", saveErr)
			t.FailNow()
		}
	}

	searchResult, getErr := userRepository.SearchWithFilters(ctx, givenFilter)
	assert.NoError(t, getErr)
	if getErr != nil {
		t.Errorf("unexpected error trying to get a user by its id: %s", getErr)
		t.FailNow()
	}

	assert.Equal(t, expectedResult, searchResult)
}
