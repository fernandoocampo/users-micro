package memorydb

import (
	"context"
	"errors"
	"fmt"
	"log"
)

const recordsLimit = 100

// DryRunRepository is the repository handler for s in a relational db.
type DryRunRepository struct {
	storage map[string]interface{}
}

// NewDryRunRepository creates a new  repository that will use a rdb.
func NewDryRunRepository() *DryRunRepository {
	newRepo := DryRunRepository{
		storage: make(map[string]interface{}),
	}
	return &newRepo
}

// Save store the given entity
func (u *DryRunRepository) Save(ctx context.Context, entityID string, entity interface{}) error {
	log.Println("level", "DEBUG", "msg", "storing entity", "method", "memory.DryRunRepository.Save", "entity", entity)
	if len(u.storage) >= recordsLimit {
		log.Println("level", "ERROR", "msg", "cannot save given entity, because the limit of allowed records was exceeded", "limit", recordsLimit)
		return errors.New("cannot save given entity, contact administrator")
	}
	if entityID == "" {
		log.Println("level", "ERROR", "msg", "cannot save given entity, because it doesn't contain a valid id", "entity", entity)
		return fmt.Errorf("cannot save given entity %v, because it doesn't contain a valid id", entity)
	}
	u.storage[entityID] = entity
	return nil
}

// Update update the given entity
func (u *DryRunRepository) Update(ctx context.Context, entityID string, entity interface{}) error {
	log.Println("level", "DEBUG", "msg", "updating entity", "method", "memory.DryRunRepository.Update", "entity", entity)
	if entityID == "" {
		log.Println("level", "ERROR", "msg", "cannot update given entity, because it doesn't contain a valid id", "entity", entity)
		return fmt.Errorf("cannot update given entity %v, because it doesn't contain a valid id", entity)
	}
	entity, ok := u.storage[entityID]
	if !ok {
		return errors.New("given entity doesn't exist")
	}
	u.storage[entityID] = entity
	return nil
}

// FindByID finds a entity with the given id in the memory storate of this dry run database.
func (u *DryRunRepository) FindByID(ctx context.Context, entityID string) (interface{}, error) {
	log.Println("level", "DEBUG", "msg", "reading entity", "method", "memory.DryRunRepository.FindByID", "entity id", entityID)
	entity, ok := u.storage[entityID]
	if !ok {
		return nil, nil
	}
	log.Println("level", "DEBUG", "msg", "entity found", "method", "memory.DryRunRepository.FindByID", "entity id", entityID, "entity", entity)
	return entity, nil
}

// FindAll return all entities
func (u *DryRunRepository) FindAll(ctx context.Context) ([]interface{}, error) {
	log.Println("level", "DEBUG", "msg", "reading all entities", "method", "memory.DryRunRepository.FindAll")
	result := make([]interface{}, 0)
	for _, v := range u.storage {
		result = append(result, v)
	}
	return result, nil
}

// Count counts records in the memory repo
func (u *DryRunRepository) Count() int {
	return len(u.storage)
}
