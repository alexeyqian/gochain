package statusdb

import "github.com/alexeyqian/gochain/entity"

// used as DI for easy testing
// memory data provider for tesging
// file data provider for production
type Storage interface {
	Open()
	Close()
	Remove()

	// New GP style DB interface
	GetAll(table string) []entity.Entity
	Get(key string) (entity.Entity, error)
	Put(key string, e entity.Entity) error

	// Traditinal CRUD Style DB interface
	//GetByID(id string) (*Entity, error)
	//Find() ([]*Entity, error)
	//Create(user *Entity) error
	//Update(user *Entity) error
	//Delete(id string) error
}
