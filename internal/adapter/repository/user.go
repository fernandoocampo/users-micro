package repository

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Skills user skills type
type Skills []string

// User contains user data.
type User struct {
	ID string `json:"id"`
	// City user's city.
	City string `json:"city"`
	// FirstName name of the person who is owner of this user.
	FirstName string `json:"first_name"`
	// LastName last name of the person who is owner of this user.
	LastName string `json:"last_name"`
	// Skill skill of the user.
	Skills Skills `json:"skills"`
}

// FindUsersResult contains the list of users found plus some metadata.
type FindUsersResult struct {
	Users       []User
	Total       int
	Page        int
	RowsPerPage int
}

// UserFilter contains filters to search users
type UserFilter struct {
	// City user's city.
	City string
	// Skill skill of the user.
	Skills Skills
	// Page page to query
	Page int
	// rows per page
	RowsPerPage int
}

// Value make the Skills struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (s Skills) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan make the skills type implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (s *Skills) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, s)
}
