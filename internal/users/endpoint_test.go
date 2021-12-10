package users_test

import (
	"context"
	"testing"

	"github.com/fernandoocampo/users-micro/internal/adapter/repository"
	"github.com/fernandoocampo/users-micro/internal/users"
	"github.com/stretchr/testify/assert"
)

func TestGetUserSuccessfully(t *testing.T) {
	userID := "1234"
	expectedResponse := users.GetUserWithIDResult{
		User: &users.User{
			ID:        "1234",
			City:      "Cali",
			Skills:    users.UserSkills([]string{"jack"}),
			FirstName: "Alicia",
			LastName:  "Mendez",
		},
		Err: "",
	}
	userRepository := userRepoMock{
		repo: make(map[string]repository.User),
	}
	existingUser := repository.User{
		ID:        "1234",
		City:      "Cali",
		Skills:    repository.Skills([]string{"jack"}),
		FirstName: "Alicia",
		LastName:  "Mendez",
	}
	userRepository.repo[existingUser.ID] = existingUser
	userService := users.NewService(&userRepository)
	getUserEndpoint := users.MakeGetUserWithIDEndpoint(userService)
	ctx := context.TODO()

	userFound, err := getUserEndpoint(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, userFound)
}

func TestGetUserNotFound(t *testing.T) {
	userID := "1234"
	expectedResponse := users.GetUserWithIDResult{
		User: nil,
		Err:  "",
	}
	userRepository := userRepoMock{
		repo: make(map[string]repository.User),
	}
	userService := users.NewService(&userRepository)
	getUserEndpoint := users.MakeGetUserWithIDEndpoint(userService)
	ctx := context.TODO()

	userFound, err := getUserEndpoint(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, userFound)
}

func TestCreateUserSuccessfully(t *testing.T) {
	userRepository := userRepoMock{
		repo: make(map[string]repository.User),
	}
	newUser := users.NewUser{
		City:      "Cali",
		Skills:    users.UserSkills([]string{"jack"}),
		FirstName: "Alicia",
		LastName:  "Mendez",
	}
	userService := users.NewService(&userRepository)
	createUserEndpoint := users.MakeCreateUserEndpoint(userService)
	ctx := context.TODO()

	result, err := createUserEndpoint(ctx, &newUser)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestUpdateUserSuccessfully(t *testing.T) {
	expectedResponse := users.UpdateUserResult{
		Err: "",
	}
	userRepository := userRepoMock{
		repo: make(map[string]repository.User),
	}
	updatewUser := users.UpdateUser{
		ID:        "123",
		City:      "Cali",
		Skills:    users.UserSkills([]string{"jack"}),
		FirstName: "Alicia",
		LastName:  "Mendez",
	}
	userService := users.NewService(&userRepository)
	updateUserEndpoint := users.MakeUpdateUserEndpoint(userService)
	ctx := context.TODO()

	result, err := updateUserEndpoint(ctx, &updatewUser)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
}

func TestSearchUsersEndpointSuccessfully(t *testing.T) {
	givenFilter := users.SearchUserFilter{
		City:        "Cali",
		Skills:      users.UserSkills([]string{"gardener"}),
		Page:        1,
		RowsPerPage: 10,
	}
	expectedSearchResult := users.SearchUsersResult{
		Users: []users.User{
			{
				ID:        "1234",
				City:      "Cali",
				Skills:    users.UserSkills([]string{"jack", "gardener"}),
				FirstName: "Alicia",
				LastName:  "Mendez",
			},
			{
				ID:        "1240",
				City:      "Cali",
				Skills:    users.UserSkills([]string{"gardener", "painter"}),
				FirstName: "Oliver",
				LastName:  "Vasquez",
			},
		},
		Total:       2,
		Page:        1,
		RowsPerPage: 10,
	}
	expectedResult := users.SearchUsersDataResult{
		SearchResult: &expectedSearchResult,
	}
	searchResultFixture := repository.FindUsersResult{
		Users: []repository.User{
			{
				ID:        "1234",
				City:      "Cali",
				Skills:    repository.Skills([]string{"jack", "gardener"}),
				FirstName: "Alicia",
				LastName:  "Mendez",
			},
			{
				ID:        "1240",
				City:      "Cali",
				Skills:    repository.Skills([]string{"gardener", "painter"}),
				FirstName: "Oliver",
				LastName:  "Vasquez",
			},
		},
		Total:       2,
		Page:        1,
		RowsPerPage: 10,
	}
	userRepository := userRepoMock{
		repo:         make(map[string]repository.User),
		searchResult: searchResultFixture,
	}
	userService := users.NewService(&userRepository)
	searchUserEndpoint := users.MakeSearchUsersEndpoint(userService)
	ctx := context.TODO()

	usersFound, err := searchUserEndpoint(ctx, givenFilter)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, usersFound)
}
