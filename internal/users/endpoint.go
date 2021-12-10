package users

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints is a wrapper for endpoints
type Endpoints struct {
	GetUserWithIDEndpoint endpoint.Endpoint
	CreateUserEndpoint    endpoint.Endpoint
	UpdateUserEndpoint    endpoint.Endpoint
	SearchUsersEndpoint   endpoint.Endpoint
}

// NewEndpoints Create the endpoints for users-micro application.
func NewEndpoints(service *Service) Endpoints {
	return Endpoints{
		GetUserWithIDEndpoint: MakeGetUserWithIDEndpoint(service),
		CreateUserEndpoint:    MakeCreateUserEndpoint(service),
		UpdateUserEndpoint:    MakeUpdateUserEndpoint(service),
		SearchUsersEndpoint:   MakeSearchUsersEndpoint(service),
	}
}

// MakeGetUserWithIDEndpoint create endpoint for get a user with ID service.
func MakeGetUserWithIDEndpoint(srv *Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		userID, ok := request.(string)
		if !ok {
			log.Println("level", "ERROR", "msg", "invalid user id", "received", fmt.Sprintf("%t", request))
			return nil, errors.New("invalid user id")
		}

		userFound, err := srv.GetUserWithID(ctx, userID)
		if err != nil {
			log.Println(
				"level", "ERROR",
				"msg", "something went wrong trying to get an user with the given id",
				"error", err,
			)
		}
		log.Println("level", "DEBUG", "msg", "find user by id endpoint", "result", userFound)
		return newGetUserWithIDResult(userFound, err), nil
	}
}

// MakeCreateUserEndpoint create endpoint for create user service.
func MakeCreateUserEndpoint(srv *Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		newUser, ok := request.(*NewUser)
		if !ok {
			log.Println("level", "ERROR", "msg", "invalid new user type", "received", fmt.Sprintf("%t", request))
			return nil, errors.New("invalid new user type")
		}

		newid, err := srv.Create(ctx, *newUser)
		if err != nil {
			log.Println(
				"level", "ERROR",
				"msg", "something went wrong trying to create an user with the given id",
				"error", err,
			)
		}
		return newCreateUserResult(newid, err), nil
	}
}

// MakeUpdateUserEndpoint create endpoint for update user service.
func MakeUpdateUserEndpoint(srv *Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		updateUser, ok := request.(*UpdateUser)
		if !ok {
			log.Println("level", "ERROR", "msg", "invalid update user type", "received", fmt.Sprintf("%t", request))
			return nil, errors.New("invalid update user type")
		}

		err := srv.Update(ctx, *updateUser)
		if err != nil {
			log.Println(
				"level", "ERROR",
				"msg", "something went wrong trying to update an user with the given id",
				"error", err,
			)
		}
		return newUpdateUserResult(err), nil
	}
}

// MakeSearchUsersEndpoint user endpoint to search users with filters.
func MakeSearchUsersEndpoint(srv *Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		userFilters, ok := request.(SearchUserFilter)
		if !ok {
			log.Println("level", "ERROR", "msg", "invalid user filters", "received", fmt.Sprintf("%t", request))
			return nil, errors.New("invalid user filters")
		}

		searchResult, err := srv.SearchUsers(ctx, userFilters)
		if err != nil {
			log.Println(
				"level", "ERROR",
				"msg", "something went wrong trying to search users with the given filter",
				"error", err,
			)
		}
		log.Println("level", "DEBUG", "msg", "search users endpoint", "result", searchResult)
		return newSearchUsersDataResult(searchResult, err), nil
	}
}
