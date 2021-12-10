package memorydb_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/fernandoocampo/users-micro/internal/adapter/memorydb"
	"github.com/fernandoocampo/users-micro/internal/adapter/repository"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserInMemoryDB(t *testing.T) {
	userID := "sfsfsf-sdfsf1234"
	newUser := repository.User{
		ID:        userID,
		FirstName: "users-micro",
		LastName:  "Wayne",
		City:      "Medellin",
		Skills:    []string{"work"},
	}
	newDB := memorydb.NewDryRunRepository()
	ctx := context.TODO()

	err := newDB.Save(ctx, userID, newUser)
	savedUser, readErr := newDB.FindByID(ctx, userID)

	assert.NoError(t, err)
	assert.NoError(t, readErr)
	assert.Equal(t, newUser, savedUser)
}

func TestCreateUserWithRepository(t *testing.T) {
	userID := "sfsfsf-sdfsf1234"
	newUser := repository.User{
		ID:        userID,
		FirstName: "users-micro",
		LastName:  "Wayne",
		City:      "Medellin",
		Skills:    []string{"work"},
	}
	newDB := memorydb.NewUserDryRunRepository()
	ctx := context.TODO()

	err := newDB.Save(ctx, newUser)
	savedUser, readErr := newDB.FindByID(ctx, userID)

	assert.NoError(t, err)
	assert.NoError(t, readErr)
	assert.Equal(t, &newUser, savedUser)
}

func TestCreateUserInMemoryDBWithLimit(t *testing.T) {
	newDB := memorydb.NewDryRunRepository()
	ctx := context.TODO()
	newUser := repository.User{
		FirstName: "users-micro",
		LastName:  "Wayne",
		City:      "Medellin",
		Skills:    []string{"work"},
	}

	for i := 0; i < 100; i++ {
		userID := strconv.Itoa(i)
		newUser.ID = userID
		err := newDB.Save(ctx, userID, newUser)
		if !assert.NoError(t, err) {
			t.FailNow()
		}
	}
	userID := "100"
	err := newDB.Save(ctx, userID, newUser)
	assert.Error(t, err)
	assert.Equal(t, 100, newDB.Count())
}
