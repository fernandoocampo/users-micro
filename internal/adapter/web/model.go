package web

import "github.com/fernandoocampo/users-micro/internal/users"

// Result standard result for the service
type Result struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Errors  []string    `json:"errors"`
}

// User contains user data.
type User struct {
	ID string `json:"id"`
	// City user's city.
	City string `json:"city"`
	// Skill skill of the user.
	Skills []string `json:"skills"`
	// FirstName name of the person who is owner of this user.
	FirstName string `json:"first_name"`
	// LastName last name of the person who is owner of this user.
	LastName string `json:"last_name"`
}

// NewUser contains the expected data for a new user.
type NewUser struct {
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	City      string   `json:"city"`
	Skills    []string `json:"skills"`
}

// UpdateUser contains the expected data to update an user.
type UpdateUser struct {
	ID        string   `json:"id"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	City      string   `json:"city"`
	Skills    []string `json:"skills"`
}

// CreateUserResponse standard response for create User
type CreateUserResponse struct {
	ID  string `json:"id"`
	Err string `json:"err,omitempty"`
}

// GetUserWithIDResponse standard response for get a User with an ID.
type GetUserWithIDResponse struct {
	User *User  `json:"user"`
	Err  string `json:"err,omitempty"`
}

// SearchUsersResponse standard response for searching users with filters.
type SearchUsersResponse struct {
	Users *SearchUsersResult `json:"result"`
	Err   string             `json:"err,omitempty"`
}

// SearchUserFilter contains filters to search users
type SearchUserFilter struct {
	// City user's city.
	City string
	// Skill skill of the user.
	Skills []string
	// Page page to query
	Page int
	// rows per page
	PageSize int
}

// SearchUsersResult contains search users result data.
type SearchUsersResult struct {
	Users    []User `json:"users"`
	Total    int    `json:"total"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

// toUser transforms new user to a user object.
func toUser(user *users.User) *User {
	if user == nil {
		return nil
	}
	webUser := User{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Skills:    user.Skills,
		City:      user.City,
	}
	return &webUser
}

// toSearchUserResult transforms new user to a user object.
func toSearchUserResult(result *users.SearchUsersResult) *SearchUsersResult {
	if result == nil {
		return nil
	}
	usersFound := make([]User, 0)
	for _, v := range result.Users {
		userFound := toUser(&v)
		usersFound = append(usersFound, *userFound)
	}
	webUser := SearchUsersResult{
		Users:    usersFound,
		Total:    result.Total,
		Page:     result.Page,
		PageSize: result.RowsPerPage,
	}
	return &webUser
}

// toUser transforms new user to a user object.
func (n *NewUser) toUser() *users.NewUser {
	if n == nil {
		return nil
	}
	userDomain := users.NewUser{
		FirstName: n.FirstName,
		LastName:  n.LastName,
		Skills:    n.Skills,
		City:      n.City,
	}
	return &userDomain
}

// toUser transforms udpate user to a user object.
func (u *UpdateUser) toUser() *users.UpdateUser {
	if u == nil {
		return nil
	}
	userDomain := users.UpdateUser{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Skills:    u.Skills,
		City:      u.City,
	}
	return &userDomain
}

func toCreateUserResponse(userResult users.CreateUserResult) Result {
	var message Result
	if userResult.Err == "" {
		message.Success = true
		message.Data = userResult.ID
	}
	if userResult.Err != "" {
		message.Errors = []string{userResult.Err}
	}
	return message
}

func toUpdateUserResponse(userResult users.UpdateUserResult) Result {
	var message Result
	if userResult.Err == "" {
		message.Success = true
	}
	if userResult.Err != "" {
		message.Errors = []string{userResult.Err}
	}
	return message
}

func toGetUserWithIDResponse(userResult users.GetUserWithIDResult) Result {
	var message Result
	newUser := toUser(userResult.User)
	if userResult.Err == "" {
		message.Success = true
		message.Data = newUser
	}
	if userResult.Err != "" {
		message.Errors = []string{userResult.Err}
	}
	return message
}

func toSearchUsersResponse(userResult users.SearchUsersDataResult) Result {
	var message Result

	if userResult.Err == "" {
		message.Success = true
		message.Data = toSearchUserResult(userResult.SearchResult)
	}
	if userResult.Err != "" {
		message.Errors = []string{userResult.Err}
	}
	return message
}

func (s SearchUserFilter) toSearchUserFilter() users.SearchUserFilter {
	return users.SearchUserFilter{
		City:        s.City,
		Skills:      s.Skills,
		Page:        s.Page,
		RowsPerPage: s.PageSize,
	}
}
