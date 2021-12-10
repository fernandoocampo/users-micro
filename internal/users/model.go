package users

import (
	"encoding/json"

	"github.com/fernandoocampo/users-micro/internal/adapter/repository"
)

// UserSkills a list of user skills
type UserSkills []string

// CreateUserResult standard response for create User
type CreateUserResult struct {
	ID  string
	Err string
}

// UpdateUserResult standard response for updating a user
type UpdateUserResult struct {
	Err string
}

// GetUserWithIDResult standard roespnse for get a User with an ID.
type GetUserWithIDResult struct {
	User *User
	Err  string
}

// SearchUsersDataResult standard roespnse for get a User with an ID.
type SearchUsersDataResult struct {
	SearchResult *SearchUsersResult
	Err          string
}

// SearchUserFilter contains filters to search users
type SearchUserFilter struct {
	// City user's city.
	City string
	// Skill skill of the user.
	Skills UserSkills
	// Page page to query
	Page int
	// rows per page
	RowsPerPage int
}

// SearchUsersResult contains search users result data.
type SearchUsersResult struct {
	Users       []User
	Total       int
	Page        int
	RowsPerPage int
}

// NewUser contains user data.
type NewUser struct {
	// City user's city.
	City string `json:"city"`
	// Skill skill of the user.
	Skills UserSkills `json:"skills"`
	// FirstName name of the person who is owner of this user.
	FirstName string `json:"first_name"`
	// LastName last name of the person who is owner of this user.
	LastName string `json:"last_name"`
}

// UpdateUser contains user data to update.
type UpdateUser struct {
	ID string `json:"id"`
	// City user's city.
	City string `json:"city"`
	// Skill skill of the user.
	Skills UserSkills `json:"skills"`
	// FirstName name of the person who is owner of this user.
	FirstName string `json:"first_name"`
	// LastName last name of the person who is owner of this user.
	LastName string `json:"last_name"`
}

// User contains user data.
type User struct {
	ID string `json:"id"`
	// City user's city.
	City string `json:"city"`
	// Skill skill of the user.
	Skills UserSkills `json:"skills"`
	// FirstName name of the person who is owner of this user.
	FirstName string `json:"first_name"`
	// LastName last name of the person who is owner of this user.
	LastName string `json:"last_name"`
}

// ToUserPortOut transforms new user to a user port out.
func (u User) ToUserPortOut() repository.User {
	return repository.User{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Skills:    u.Skills.toDBSkills(),
		City:      u.City,
	}
}

// NewUser transforms new user to a user port out.
func (n NewUser) NewUser(userID string) User {
	return User{
		ID:        userID,
		FirstName: n.FirstName,
		LastName:  n.LastName,
		Skills:    UserSkills(n.Skills),
		City:      n.City,
	}
}

// UpdateUser transforms update user to a user port out.
func (u UpdateUser) UpdateUser() User {
	return User{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Skills:    UserSkills(u.Skills),
		City:      u.City,
	}
}

func toUserSkills(repoSkills repository.Skills) UserSkills {
	userSkills := UserSkills(repoSkills)
	return userSkills
}

func (u UserSkills) toDBSkills() repository.Skills {
	repoSkills := repository.Skills(u)
	return repoSkills
}

// transformUserPortOuttoUser transforms the given user port out to service user.
func transformUserPortOuttoUser(userRepo *repository.User) *User {
	if userRepo == nil {
		return nil
	}
	newuser := User{
		ID:        userRepo.ID,
		City:      userRepo.City,
		FirstName: userRepo.FirstName,
		LastName:  userRepo.LastName,
		Skills:    toUserSkills(userRepo.Skills),
	}
	return &newuser
}

// newGetUserWithIDResult create a new GetUserWithIDResult
func newGetUserWithIDResult(user *User, err error) GetUserWithIDResult {
	var errmessage string
	if err != nil {
		errmessage = err.Error()
	}
	return GetUserWithIDResult{
		User: user,
		Err:  errmessage,
	}
}

// newSearchUsersResult create a new SearchUsersResult
func newSearchUsersDataResult(result *SearchUsersResult, err error) SearchUsersDataResult {
	var errmessage string
	if err != nil {
		errmessage = err.Error()
	}
	return SearchUsersDataResult{
		SearchResult: result,
		Err:          errmessage,
	}
}

// newCreateUserResult create a new CreateUserResponse
func newCreateUserResult(id string, err error) CreateUserResult {
	var errmessage string
	if err != nil {
		errmessage = err.Error()
	}
	return CreateUserResult{
		ID:  id,
		Err: errmessage,
	}
}

// newUpdateUserResult udpate a new UpdateUserResponse
func newUpdateUserResult(err error) UpdateUserResult {
	var errmessage string
	if err != nil {
		errmessage = err.Error()
	}
	return UpdateUserResult{
		Err: errmessage,
	}
}

func (u User) String() string {
	b, err := json.Marshal(u)
	if err != nil {
		return ""
	}
	return string(b)
}

func (n NewUser) String() string {
	b, err := json.Marshal(n)
	if err != nil {
		return ""
	}
	return string(b)
}

func (s SearchUserFilter) toRepositoryFilters() repository.UserFilter {
	return repository.UserFilter{
		City:        s.City,
		Skills:      s.Skills.toDBSkills(),
		Page:        s.Page,
		RowsPerPage: s.RowsPerPage,
	}
}

func toSearchUsersResult(repoResult repository.FindUsersResult) SearchUsersResult {
	var userCollection []User
	for _, v := range repoResult.Users {
		userFound := &v
		userToAdd := transformUserPortOuttoUser(userFound)
		userCollection = append(userCollection, *userToAdd)
	}
	return SearchUsersResult{
		Users:       userCollection,
		Total:       repoResult.Total,
		Page:        repoResult.Page,
		RowsPerPage: repoResult.RowsPerPage,
	}
}
