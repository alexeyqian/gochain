package statusdb

import (
	"strings"

	"github.com/alexeyqian/gochain/entity"
)

var _data map[string]entity.Entity

type MemoryStorage struct {
}

func (dp *MemoryStorage) Open() {
	_data = make(map[string]entity.Entity)
}

func (dp *MemoryStorage) Close() {

}

func (dp *MemoryStorage) Remove() {
	for k := range _data {
		delete(_data, k)
	}
}

func (dp *MemoryStorage) GetAll(table string) []entity.Entity {
	var res []entity.Entity
	for key, value := range _data {
		if strings.HasPrefix(key, table+"_") {
			res = append(res, value)
		}
	}
	return res
}

func (dp *MemoryStorage) Get(key string) (entity.Entity, error) {
	return _data[key], nil
}

func (dp *MemoryStorage) Put(key string, e entity.Entity) error {
	_data[key] = e
	return nil
}
