package store

import (
	"fmt"
	"os"
	"time"

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

func (s *BoltStorage) Open() error {
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

	return nil
}

func (s *BoltStorage) Close() error {
	s.db.Close()
	return nil
}

func (s *BoltStorage) Remove() error {
	os.Remove(s.pathname)
	return nil
}

func (s *BoltStorage) Get(bucket, key string) ([]byte, error) {
	var temp []byte

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket: %s not exist", bucket)
		}

		data := b.Get([]byte(key))
		temp = make([]byte, len(data))
		// @attention Have to duplicate the data out, it will invalid out tx
		copy(temp, data)
		return nil
	})
	return temp, err
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

func (s *BoltStorage) Delete(bucket, key string) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil // no bucket, do nothing
		}

		return b.Delete([]byte(key))
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
			res = append(res, temp)
		}

		return nil
	})

	return res, err
}

func (s *BoltStorage) HasBucket(bucket string) bool {
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b != nil {
			return nil
		} else {
			return fmt.Errorf("not exist")
		}
	})

	return err == nil
}

func (s *BoltStorage) HasKey(bucket, key string) bool {
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket: %s not exist", bucket)
		}

		data := b.Get([]byte(key))
		if data == nil {
			return fmt.Errorf("cannot find key")
		}

		return nil
	})

	return err == nil
}

func (s *BoltStorage) CreateBucket(bucket string) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			tx.CreateBucket([]byte(bucket))
		}

		return nil
	})

	return err
}
