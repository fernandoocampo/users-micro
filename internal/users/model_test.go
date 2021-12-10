package users_test

import (
	"testing"

	"github.com/fernandoocampo/users-micro/internal/adapter/repository"
	"github.com/fernandoocampo/users-micro/internal/users"
	"github.com/stretchr/testify/assert"
)

func TestToUserPortOut(t *testing.T) {
	expectedRepoUser := repository.User{
		ID:        "1234",
		City:      "Cali",
		FirstName: "Lucia",
		LastName:  "Mendez",
		Skills:    repository.Skills([]string{"worker", "happy"}),
	}
	givenUser := users.User{
		ID:        "1234",
		City:      "Cali",
		FirstName: "Lucia",
		LastName:  "Mendez",
		Skills:    users.UserSkills([]string{"worker", "happy"}),
	}

	got := givenUser.ToUserPortOut()

	assert.Equal(t, expectedRepoUser, got)
}

func TestNewUser(t *testing.T) {
	expectedUser := users.User{
		ID:        "1234",
		City:      "Cali",
		FirstName: "Lucia",
		LastName:  "Mendez",
		Skills:    users.UserSkills([]string{"worker", "happy"}),
	}
	givenUserID := "1234"
	givenNewUser := users.NewUser{
		City:      "Cali",
		FirstName: "Lucia",
		LastName:  "Mendez",
		Skills:    users.UserSkills([]string{"worker", "happy"}),
	}

	got := givenNewUser.NewUser(givenUserID)

	assert.Equal(t, expectedUser, got)
}
