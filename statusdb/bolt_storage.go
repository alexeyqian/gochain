package statusdb

import (
	"fmt"
	"os"
	"time"

	"github.com/alexeyqian/gochain/config"

	"github.com/alexeyqian/gochain/entity"
	"github.com/boltdb/bolt"
)

// BoltStorage implemented with Bolt DB
type BoltStorage struct {
}

var _db *bolt.DB

// Open database
func (dp *BoltStorage) Open() {
	var err error

	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	// Bolt obtains a file lock on the data file so multiple processes cannot open the same database at the same time.
	// Opening an already open Bolt database will cause it to hang until the other process closes it.
	// To prevent an indefinite wait you can pass a timeout option to the Open()
	_db, err = bolt.Open(config.BoltDbFileName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err)
	}
	/*
		err = _db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(GpoTable))

			if b == nil { // create all buckets
				tx.CreateBucket([]byte(GpoTable))
				tx.CreateBucket([]byte(AccountTable))
				tx.CreateBucket([]byte(ArticleTable))
				tx.CreateBucket([]byte(CommentTable))
				tx.CreateBucket([]byte(VoteTable))
			}

			return nil
		})*/

}

// Close database
func (dp *BoltStorage) Close() {
	_db.Close()
}

// Remove all data in db
func (dp *BoltStorage) RemoveAll() {
	os.Remove(config.BoltDbFileName)
}

// Get an entity
func (dp *BoltStorage) Get(bucket, key string) ([]byte, error) {
	err = _db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(entityType))
		if b == nil {
			return nil, fmt.Errorf("bucket: %s not exist", entityType)
		}
		data := b.Get([]byte(key))
		temp = make([]byte, len(data))
		// @attention Have to duplicate the data out, it will invalid out tx
		copy(temp, data)
		return temp, err
	})
	return result, err
}

func (dp *BoltStorage) Put(bucket, key string, data []byte) error {
	err := _db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			tx.CreateBucket([]byte(bucket))
		}

		err := b.Put([]byte(key), data)

		return err
	})

	return err
}

func (s *BoldStorage) Remove(key string) error {
	var err error
	err = _db.Update(func(tx *bolt.Tx) error {
		data, _ := serializeEntity(entityType, e)
		b := tx.Bucket([]byte(entityType))
		err = b.Delete([]byte(key))

		return err
	})

	return err
}

// GetAll data from table
func (dp *BoltStorage) GetAll(table string) []entity.Entity {
	var err error
	var temp []byte
	var res []entity.Entity

	err = _db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			temp = make([]byte, len(v))
			copy(temp, v) // @attention Have to duplicate v out, which will be invalid out side of tx
			e, _ := deserializeEntity(table, temp)
			res = append(res, e)
		}

		return err
	})

	return res
}

/*
func serializeEntity(entityType string, e entity.Entity) ([]byte, error) {
	var err error
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	switch entityType {
	case GpoTable:
		ce := e.(entity.Gpo)
		err = encoder.Encode(ce)
		return result.Bytes(), err
	case AccountTable:
		ce := e.(entity.Account)
		err = encoder.Encode(ce)
		return result.Bytes(), err
	case ArticleTable:
		ce := e.(entity.Article)
		err = encoder.Encode(ce)
		return result.Bytes(), err
	case CommentTable:
		ce := e.(entity.Comment)
		err = encoder.Encode(ce)
		return result.Bytes(), err
	case VoteTable:
		ce := e.(entity.Vote)
		err = encoder.Encode(ce)
		return result.Bytes(), err
	}

	panic("serialize: unknown entity type")
}

func deserializeEntity(entityType string, data []byte) (entity.Entity, error) {
	var err error
	decoder := gob.NewDecoder(bytes.NewReader(data))

	switch entityType {
	case GpoTable:
		var e entity.Gpo
		err = decoder.Decode(&e)
		return e, err
	case AccountTable:
		var e entity.Account
		err = decoder.Decode(&e)
		return e, err
	case ArticleTable:
		var e entity.Article
		err = decoder.Decode(&e)
		return e, err
	case CommentTable:
		var e entity.Comment
		err = decoder.Decode(&e)
		return e, err
	case VoteTable:
		var e entity.Vote
		err = decoder.Decode(&e)
		return e, err
	}

	panic("deserialize: unknown entity type")
}

func getPrefix(key string) string {
	index := bytes.Index([]byte(key), []byte("_"))
	return key[0:index]
}
*/
