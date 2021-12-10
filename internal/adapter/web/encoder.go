package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/fernandoocampo/users-micro/internal/users"
)

func encodeCreateUserResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	result, ok := response.(users.CreateUserResult)
	if !ok {
		log.Println("level", "ERROR", "msg", "cannot transform to users.CreateUserResult", "received", fmt.Sprintf("%+v", response))
		return errors.New("cannot build create user response")
	}
	w.Header().Set("Content-Type", "application/json")
	message := toCreateUserResponse(result)
	return json.NewEncoder(w).Encode(message)
}

func encodeUpdateUserResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	result, ok := response.(users.UpdateUserResult)
	if !ok {
		log.Println("level", "ERROR", "msg", "cannot transform to users.UpdateUserResult", "received", fmt.Sprintf("%+v", response))
		return errors.New("cannot build update user response")
	}
	w.Header().Set("Content-Type", "application/json")
	message := toUpdateUserResponse(result)
	return json.NewEncoder(w).Encode(message)
}

func encodeGetUserWithIDResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	result, ok := response.(users.GetUserWithIDResult)
	if !ok {
		log.Println("level", "ERROR", "msg", "cannot transform to users.GetUserWithIDResult", "received", fmt.Sprintf("%+v", response))
		return errors.New("cannot build get user response")
	}
	w.Header().Set("Content-Type", "application/json")
	message := toGetUserWithIDResponse(result)
	return json.NewEncoder(w).Encode(message)
}

func encodeSearchUsersResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	result, ok := response.(users.SearchUsersDataResult)
	if !ok {
		log.Println("level", "ERROR", "msg", "cannot transform to users.SearchUsersDataResult", "received", fmt.Sprintf("%T", response))
		return errors.New("cannot build search users response")
	}
	w.Header().Set("Content-Type", "application/json")
	message := toSearchUsersResponse(result)
	return json.NewEncoder(w).Encode(message)
}
