package web

import (
	"net/http"

	"github.com/fernandoocampo/users-micro/internal/users"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// NewHTTPServer is a factory to create http servers for this project.
func NewHTTPServer(endpoints users.Endpoints) http.Handler {
	router := mux.NewRouter()
	router.Methods(http.MethodGet).Path("/users/{id}").Handler(
		httptransport.NewServer(
			endpoints.GetUserWithIDEndpoint,
			decodeGetUserWithIDRequest,
			encodeGetUserWithIDResponse),
	)
	router.Methods(http.MethodPost).Path("/users").Handler(
		httptransport.NewServer(
			endpoints.CreateUserEndpoint,
			decodeCreateUserRequest,
			encodeCreateUserResponse),
	)
	router.Methods(http.MethodPut).Path("/users").Handler(
		httptransport.NewServer(
			endpoints.UpdateUserEndpoint,
			decodeUpdateUserRequest,
			encodeUpdateUserResponse),
	)
	router.Methods(http.MethodGet).Path("/users").Handler(
		httptransport.NewServer(
			endpoints.SearchUsersEndpoint,
			decodeSearchUsersRequest,
			encodeSearchUsersResponse),
	)
	return router
}
