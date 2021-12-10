package users

import (
	"context"
	"log"

	"github.com/fernandoocampo/users-micro/internal/adapter/repository"
	"github.com/google/uuid"
)

// Repository defines portout behavior to send user data to external platforms.
type Repository interface {
	FindByID(ctx context.Context, userID string) (*repository.User, error)
	Save(ctx context.Context, user repository.User) error
	Update(ctx context.Context, user repository.User) error
	SearchWithFilters(ctx context.Context, filter repository.UserFilter) (repository.FindUsersResult, error)
}

// Service implements user management logic.
type Service struct {
	userRepository Repository
}

// NewService creates a new application service
func NewService(userRepository Repository) *Service {
	return &Service{
		userRepository: userRepository,
	}
}

// GetUserWithID get the user with the given id.
func (s *Service) GetUserWithID(ctx context.Context, userID string) (*User, error) {
	log.Println(
		"level", "DEBUG",
		"msg", "getting user with id",
		"method", "Service.GetUserWithID",
		"userID", userID)
	result, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		log.Println("level", "ERROR",
			"msg", "something went wrong trying to get an user",
			"method", "Service.GetUserWithID", "userID", userID,
		)
		return nil, err
	}
	log.Println("level", "DEBUG", "endpoint result 1", result)

	user := transformUserPortOuttoUser(result)
	log.Println(
		"level", "DEBUG",
		"msg", "user was found",
		"method", "Service.GetUserWithID",
		"user", user)
	return user, nil
}

// Create creates an user
func (s *Service) Create(ctx context.Context, newuser NewUser) (string, error) {
	log.Println(
		"level", "DEBUG",
		"msg", "creating user",
		"method", "Service.Create",
		"newuser", newuser)
	id := uuid.New().String()
	user := newuser.NewUser(id)
	log.Println(
		"level", "DEBUG",
		"msg", "creating user",
		"method", "Service.Create",
		"user", user)
	err := s.userRepository.Save(ctx, user.ToUserPortOut())
	if err != nil {
		log.Println("level", "ERROR",
			"msg", "something goes wrong creating user",
			"method", "Service.Create", "user", user,
		)
		return "", err
	}
	log.Println(
		"level", "INFO",
		"msg", "user was created successfuly",
		"method", "Service.Create",
		"user", user)
	return id, nil
}

// Update updates an user
func (s *Service) Update(ctx context.Context, userToUpdate UpdateUser) error {
	log.Println(
		"level", "DEBUG",
		"msg", "updating user",
		"method", "Service.Update",
		"user", userToUpdate)
	user := userToUpdate.UpdateUser()
	log.Println(
		"level", "DEBUG",
		"msg", "updating user",
		"method", "Service.Update",
		"user", user)
	err := s.userRepository.Update(ctx, user.ToUserPortOut())
	if err != nil {
		log.Println("level", "ERROR",
			"msg", "something goes wrong updating user",
			"method", "Service.Update", "user", user,
		)
		return err
	}
	log.Println(
		"level", "INFO",
		"msg", "user was updated successfuly",
		"method", "Service.Update",
		"user", user)
	return nil
}

// SearchUsers search users who match the given filters
func (s *Service) SearchUsers(ctx context.Context, givenFilter SearchUserFilter) (*SearchUsersResult, error) {
	log.Println(
		"level", "DEBUG",
		"msg", "searching users",
		"method", "Service.SearchUsers",
		"filter", givenFilter,
	)
	filters := givenFilter.toRepositoryFilters()

	repoResult, err := s.userRepository.SearchWithFilters(ctx, filters)
	if err != nil {
		log.Println("level", "ERROR",
			"msg", "something goes wrong searching users",
			"method", "Service.SearchUsers",
			"filter", givenFilter,
		)
		return nil, err
	}

	result := toSearchUsersResult(repoResult)

	return &result, nil
}
