package users_test

import (
	"context"
	"errors"
	"testing"

	"github.com/fernandoocampo/users-micro/internal/adapter/repository"
	"github.com/fernandoocampo/users-micro/internal/users"
	"github.com/stretchr/testify/assert"
)

func TestFindUserSuccessfully(t *testing.T) {
	userID := "1234"
	expectedUser := users.User{
		ID:        "1234",
		City:      "Cali",
		Skills:    users.UserSkills([]string{"jack"}),
		FirstName: "Alicia",
		LastName:  "Mendez",
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
	ctx := context.TODO()

	userFound, err := userService.GetUserWithID(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, &expectedUser, userFound)
}

func TestFindUserNotFound(t *testing.T) {
	userID := "1234"
	userRepository := userRepoMock{
		repo: make(map[string]repository.User),
	}
	userService := users.NewService(&userRepository)
	ctx := context.TODO()

	userFound, err := userService.GetUserWithID(ctx, userID)

	assert.NoError(t, err)
	assert.Nil(t, userFound)
}

func TestFindUserWithError(t *testing.T) {
	userID := "1234"
	userRepository := userRepoMock{
		repo: make(map[string]repository.User),
		err:  errors.New("any error"),
	}
	userService := users.NewService(&userRepository)
	ctx := context.TODO()

	userFound, err := userService.GetUserWithID(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, userFound)
	assert.Equal(t, errors.New("any error"), err)
}

func TestSearchUsersSuccessfully(t *testing.T) {
	givenFilter := users.SearchUserFilter{
		City:        "Cali",
		Skills:      users.UserSkills([]string{"gardener"}),
		Page:        1,
		RowsPerPage: 10,
	}
	expectedResult := users.SearchUsersResult{
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
	ctx := context.TODO()

	usersFound, err := userService.SearchUsers(ctx, givenFilter)

	assert.NoError(t, err)
	assert.Equal(t, &expectedResult, usersFound)
}

type userRepoMock struct {
	err          error
	repo         map[string]repository.User
	searchResult repository.FindUsersResult
}

func (u *userRepoMock) FindByID(_ context.Context, userID string) (*repository.User, error) {
	if u.err != nil {
		return nil, u.err
	}
	result, ok := u.repo[userID]
	if !ok {
		return nil, nil
	}
	return &result, nil
}

func (u *userRepoMock) Save(ctx context.Context, user repository.User) error {
	if u.err != nil {
		return u.err
	}
	u.repo[user.ID] = user
	return nil
}

func (u *userRepoMock) Update(ctx context.Context, user repository.User) error {
	if u.err != nil {
		return u.err
	}
	u.repo[user.ID] = user
	return nil
}

func (u *userRepoMock) SearchWithFilters(ctx context.Context, filter repository.UserFilter) (repository.FindUsersResult, error) {
	var result repository.FindUsersResult
	if u.err != nil {
		return result, u.err
	}
	return u.searchResult, nil
}
