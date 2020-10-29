package statusdb

import (
	"fmt"
	"os"
	"time"

	"github.com/alexeyqian/gochain/config"

	"github.com/boltdb/bolt"
)

type BoltStorage struct {
	pathname string
	db       *bolt.DB
}

func NewBoltStorage(path string) *BoltStorage {
	return &BoltStorage{
		pathname: path,
	}
}

func (s *BoltStorage) Open() {
	var err error

	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	// Bolt obtains a file lock on the data file so multiple processes cannot open the same database at the same time.
	// Opening an already open Bolt database will cause it to hang until the other process closes it.
	// To prevent an indefinite wait you can pass a timeout option to the Open()
	s.db, err = bolt.Open(s.pathname, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err)
	}
}

func (s *BoltStorage) Close() {
	s.db.Close()
}

func (s *BoltStorage) RemoveAll() {
	os.Remove(config.BoltDbFileName)
}

func (s *BoltStorage) Get(bucket, key string) ([]byte, error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil, fmt.Errorf("bucket: %s not exist", bucket)
		}
		data := b.Get([]byte(key))
		temp = make([]byte, len(data))
		// @attention Have to duplicate the data out, it will invalid out tx
		copy(temp, data)
		return temp, err
	})
	return result, err
}

func (s *BoltStorage) Put(bucket, key string, data []byte) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			tx.CreateBucket([]byte(bucket))
		}

		err := b.Put([]byte(key), data)

		return err
	})

	return err
}

func (s *BoldStorage) Delete(bucket, key string) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil // no bucket, do nothing
		}

		err := b.Delete([]byte(key))
		return err
	})

	return err
}

func (s *BoltStorage) GetAll(bucket string) ([][]byte, error) {
	var temp []byte
	var res [][]byte

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			temp = make([]byte, len(v))
			// @attention Have to duplicate v out, which will be invalid out side of tx
			copy(temp, v)
			res = append(res, e)
		}

		return err
	})

	return res, err
}

func (s *BoltStorage) HasBucket(bucket string) bool {
	return s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b != nil
	})
}

func (s *BoltStorage) CreateBucket(bucket) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			tx.CreateBucket([]byte(bucket))
		}

		return err
	})

	return err
}

/*
func serializeEntity(bucket string, e entity.Entity) ([]byte, error) {
	var err error
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	switch bucket {
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

func deserializeEntity(bucket string, data []byte) (entity.Entity, error) {
	var err error
	decoder := gob.NewDecoder(bytes.NewReader(data))

	switch bucket {
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
