package web_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fernandoocampo/users-micro/internal/adapter/web"
	"github.com/fernandoocampo/users-micro/internal/users"
	"github.com/go-kit/kit/endpoint"
	"github.com/stretchr/testify/assert"
)

type webResultGetUser struct {
	Success bool      `json:"success"`
	Data    *web.User `json:"data"`
	Errors  []string  `json:"errors"`
}

type webResultSearchUsers struct {
	Success bool                   `json:"success"`
	Data    *web.SearchUsersResult `json:"data"`
	Errors  []string               `json:"errors"`
}

type webResultCreateUser struct {
	Success bool     `json:"success"`
	Data    string   `json:"data"`
	Errors  []string `json:"errors"`
}

type webResultUpdateUser struct {
	Success bool     `json:"success"`
	Data    string   `json:"data"`
	Errors  []string `json:"errors"`
}

func TestGetUserSuccessfully(t *testing.T) {
	userID := "1234"
	expectedResponse := webResultGetUser{
		Success: true,
		Data: &web.User{
			ID:        "1234",
			City:      "Cali",
			Skills:    []string{"jack"},
			FirstName: "Lucia",
			LastName:  "Mendez",
		},
		Errors: nil,
	}
	userToReturn := users.User{
		ID:        "1234",
		City:      "Cali",
		Skills:    []string{"jack"},
		FirstName: "Lucia",
		LastName:  "Mendez",
	}
	userEndpoints := users.Endpoints{
		GetUserWithIDEndpoint: makeDummyGetUserWithIDSuccessfullyEndpoint(t, &userToReturn, nil),
	}
	httpHandler := web.NewHTTPServer(userEndpoints)
	dummyServer := httptest.NewServer(httpHandler)
	defer dummyServer.Close()

	response, err := http.Get(dummyServer.URL + "/users/" + userID)
	if err != nil {
		t.Error("unexpected error", err)
		t.FailNow()
	}
	defer response.Body.Close()

	var result webResultGetUser

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Error("unexpected error", err)
		t.FailNow()
	}

	assert.Equal(t, expectedResponse, result)
}

func TestSearchUsersSuccessfully(t *testing.T) {
	queryParams := "?city=Cali&skills=gardener&page=1&pagesize=10"
	expectedFilter := users.SearchUserFilter{
		City:        "Cali",
		Skills:      []string{"gardener"},
		Page:        1,
		RowsPerPage: 10,
	}
	expectedResponse := webResultSearchUsers{
		Success: true,
		Data: &web.SearchUsersResult{
			Users: []web.User{
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
			Total:    2,
			Page:     1,
			PageSize: 10,
		},
		Errors: nil,
	}

	serviceResult := users.SearchUsersResult{
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
	userEndpoints := users.Endpoints{
		SearchUsersEndpoint: makeDummySearchUsersSuccessfullyEndpoint(t, expectedFilter, &serviceResult, nil),
	}
	httpHandler := web.NewHTTPServer(userEndpoints)
	dummyServer := httptest.NewServer(httpHandler)
	defer dummyServer.Close()

	response, err := http.Get(dummyServer.URL + "/users" + queryParams)
	if err != nil {
		t.Error("unexpected error", err)
		t.FailNow()
	}
	defer response.Body.Close()

	var result webResultSearchUsers

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Error("unexpected error", err)
		t.FailNow()
	}

	assert.Equal(t, expectedResponse, result)
}

func TestGetUserNotFound(t *testing.T) {
	userID := "1234"
	expectedResponse := webResultGetUser{
		Success: true,
		Data:    nil,
		Errors:  nil,
	}
	userEndpoints := users.Endpoints{
		GetUserWithIDEndpoint: makeDummyGetUserWithIDSuccessfullyEndpoint(t, nil, nil),
	}
	httpHandler := web.NewHTTPServer(userEndpoints)
	dummyServer := httptest.NewServer(httpHandler)
	defer dummyServer.Close()

	response, err := http.Get(dummyServer.URL + "/users/" + userID)
	if err != nil {
		t.Error("unexpected error", err)
		t.FailNow()
	}
	defer response.Body.Close()

	var result webResultGetUser

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Error("unexpected error", err)
		t.FailNow()
	}

	assert.Equal(t, expectedResponse, result)
}

func TestGetUserWithError(t *testing.T) {
	userID := "1234"
	expectedResponse := webResultGetUser{
		Success: false,
		Data:    nil,
		Errors:  []string{"any error"},
	}
	errorToReturn := errors.New("any error")
	userEndpoints := users.Endpoints{
		GetUserWithIDEndpoint: makeDummyGetUserWithIDSuccessfullyEndpoint(t, nil, errorToReturn),
	}
	httpHandler := web.NewHTTPServer(userEndpoints)
	dummyServer := httptest.NewServer(httpHandler)
	defer dummyServer.Close()

	response, err := http.Get(dummyServer.URL + "/users/" + userID)
	if err != nil {
		t.Error("unexpected error", err)
		t.FailNow()
	}
	defer response.Body.Close()

	var result webResultGetUser

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Error("unexpected error", err)
		t.FailNow()
	}

	assert.Equal(t, expectedResponse, result)
}

func TestPostUserSuccessfully(t *testing.T) {
	newUser := web.NewUser{
		FirstName: "lucia",
		LastName:  "mendez",
		City:      "Cali",
		Skills:    []string{"jack"},
	}
	newUserJson, err := json.Marshal(newUser)
	if err != nil {
		t.Errorf("unexpected error marshalling new user: %s", err)
		t.FailNow()
	}
	userEndpoints := users.Endpoints{
		CreateUserEndpoint: makeDummyCreateUserSuccessfullyEndpoint(t, "1234", nil),
	}
	userHandler := web.NewHTTPServer(userEndpoints)

	dummyServer := httptest.NewServer(userHandler)
	defer dummyServer.Close()

	expectedResponse := webResultCreateUser{
		Success: true,
		Data:    "1234",
		Errors:  nil,
	}

	response, err := http.Post(dummyServer.URL+"/users", "application/json", bytes.NewBuffer(newUserJson))
	if err != nil {
		t.Errorf("unexected error creating post request: %s", err)
	}
	defer response.Body.Close()

	var result webResultCreateUser

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Error("unexpected error", err)
		t.FailNow()
	}

	assert.Equal(t, expectedResponse, result)
}

func TestPutUserSuccessfully(t *testing.T) {
	updateUser := web.UpdateUser{
		ID:        "123",
		FirstName: "lucia",
		LastName:  "mendez",
		City:      "Cali",
		Skills:    []string{"jack", "painter"},
	}
	updateUserJson, err := json.Marshal(updateUser)
	if err != nil {
		t.Errorf("unexpected error marshalling update user: %s", err)
		t.FailNow()
	}
	userEndpoints := users.Endpoints{
		UpdateUserEndpoint: makeDummyUpdateUserSuccessfullyEndpoint(t, nil),
	}
	userHandler := web.NewHTTPServer(userEndpoints)

	dummyServer := httptest.NewServer(userHandler)
	defer dummyServer.Close()

	expectedResponse := webResultUpdateUser{
		Success: true,
		Data:    "",
		Errors:  nil,
	}

	updateRequest, err := http.NewRequest("PUT", dummyServer.URL+"/users", bytes.NewBuffer(updateUserJson))
	if err != nil {
		t.Errorf("unexpected error creating update request: %s", err)
	}

	client := &http.Client{}
	response, err := client.Do(updateRequest)
	if err != nil {
		t.Errorf("unexected error creating put request: %s", err)
	}
	defer response.Body.Close()

	var result webResultUpdateUser

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Error("unexpected error", err)
		t.FailNow()
	}

	assert.Equal(t, expectedResponse, result)
}

func TestPostUserWithError(t *testing.T) {
	newUser := web.NewUser{
		FirstName: "lucia",
		LastName:  "mendez",
		City:      "Cali",
		Skills:    []string{"jack"},
	}
	newUserJson, err := json.Marshal(newUser)
	if err != nil {
		t.Errorf("unexpected error marshalling new user: %s", err)
		t.FailNow()
	}
	userEndpoints := users.Endpoints{
		CreateUserEndpoint: makeDummyCreateUserSuccessfullyEndpoint(t, "", errors.New("any error")),
	}
	userHandler := web.NewHTTPServer(userEndpoints)

	dummyServer := httptest.NewServer(userHandler)
	defer dummyServer.Close()

	expectedResponse := webResultCreateUser{
		Success: false,
		Data:    "",
		Errors:  []string{"any error"},
	}

	response, err := http.Post(dummyServer.URL+"/users", "application/json", bytes.NewBuffer(newUserJson))
	if err != nil {
		t.Errorf("unexected error creating post request: %s", err)
	}
	defer response.Body.Close()

	var result webResultCreateUser

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Error("unexpected error", err)
		t.FailNow()
	}

	assert.Equal(t, expectedResponse, result)
}

func makeDummyGetUserWithIDSuccessfullyEndpoint(t *testing.T, userToReturn *users.User, err error) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		t.Helper()
		userID, ok := request.(string)
		if !ok {
			t.Errorf("user id parameter is not valid: %+v", userID)
			t.FailNow()
		}

		var errMessage string
		if err != nil {
			errMessage = err.Error()
		}
		result := users.GetUserWithIDResult{
			User: userToReturn,
			Err:  errMessage,
		}
		return result, nil
	}
}

func makeDummySearchUsersSuccessfullyEndpoint(t *testing.T, expectedFilter users.SearchUserFilter, resultToReturn *users.SearchUsersResult, err error) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		t.Helper()
		filter, ok := request.(users.SearchUserFilter)
		if !ok {
			t.Errorf("user filter parameter is not valid: %T", request)
			t.FailNow()
		}

		assert.Equal(t, expectedFilter, filter)

		var errMessage string
		if err != nil {
			errMessage = err.Error()
		}
		result := users.SearchUsersDataResult{
			SearchResult: resultToReturn,
			Err:          errMessage,
		}
		return result, nil
	}
}

func makeDummyCreateUserSuccessfullyEndpoint(t *testing.T, newUserID string, err error) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		t.Helper()
		_, ok := request.(*users.NewUser)
		if !ok {
			t.Errorf("user parameter is not valid: %T", request)
			t.FailNow()
		}
		var errMessage string
		if err != nil {
			errMessage = err.Error()
		}
		result := users.CreateUserResult{
			ID:  newUserID,
			Err: errMessage,
		}
		return result, nil
	}
}

func makeDummyUpdateUserSuccessfullyEndpoint(t *testing.T, err error) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		t.Helper()
		_, ok := request.(*users.UpdateUser)
		if !ok {
			t.Errorf("update user parameter is not valid: %T", request)
			t.FailNow()
		}
		var errMessage string
		if err != nil {
			errMessage = err.Error()
		}
		result := users.UpdateUserResult{
			Err: errMessage,
		}
		return result, nil
	}
}
