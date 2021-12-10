package memorydb

import (
	"context"
	"errors"
	"log"

	"github.com/fernandoocampo/users-micro/internal/adapter/repository"
)

// UserMemoryRepository is the repository handler for users in a memory db.
type UserMemoryRepository struct {
	storage *DryRunRepository
}

// NewUserDryRunRepository creates a new user repository in a dry run repository
func NewUserDryRunRepository() *UserMemoryRepository {
	newRepo := UserMemoryRepository{
		storage: NewDryRunRepository(),
	}
	return &newRepo
}

// Save save the given user in the postgresql database.
func (u *UserMemoryRepository) Save(ctx context.Context, user repository.User) error {
	log.Println("level", "DEBUG", "msg", "storing user", "method", "repository.UserMemoryRepository.Save", "data", user)
	err := u.storage.Save(ctx, user.ID, user)
	if err != nil {
		log.Println("level", "ERROR", "msg", "storing user", "method", "repository.UserMemoryRepository.Save", "error", err)
		return errors.New("given user could not be stored")
	}
	return nil
}

// Update update the given user in the postgresql database.
func (u *UserMemoryRepository) Update(ctx context.Context, user repository.User) error {
	log.Println("level", "DEBUG", "msg", "updating user", "method", "repository.UserMemoryRepository.Update", "data", user)
	err := u.storage.Update(ctx, user.ID, user)
	if err != nil {
		log.Println("level", "ERROR", "msg", "updating user", "method", "repository.UserMemoryRepository.Update", "error", err)
		return errors.New("given user could not be updated")
	}
	return nil
}

// FindByID look for an user with the given id
func (u *UserMemoryRepository) FindByID(ctx context.Context, userID string) (*repository.User, error) {
	log.Println("level", "DEBUG", "msg", "reading user", "method", "repository.UserMemoryRepository.FindByID", "user id", userID)
	result, err := u.storage.FindByID(ctx, userID)
	if err != nil {
		log.Println("level", "ERROR", "msg", "reading user", "method", "repository.UserMemoryRepository.FindByID", "error", err)
		return nil, errors.New("something went wrong trying to get the given user id")
	}
	if result == nil {
		return nil, nil
	}
	user, ok := result.(repository.User)
	if !ok {
		log.Println("level", "ERROR", "msg", "reading user", "method", "repository.UserMemoryRepository.FindByID", "error", "unexpected object", "object", result)
		return nil, errors.New("something went wrong trying to get the given user id")
	}
	return &user, nil
}

// SearchWithFilters memory search
func (u *UserMemoryRepository) SearchWithFilters(ctx context.Context, filter repository.UserFilter) (repository.FindUsersResult, error) {
	return repository.FindUsersResult{}, nil
}
