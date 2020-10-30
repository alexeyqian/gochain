package store

import "encoding/binary"

// will instruct the storage to create an auto increment key
const AutoIncrementKey = "_auto_increment_key_"

// @IMPORTANT: if change the endian, then add keys will be incorrect
// has to be big endian now, other wise the bit order sort of the bolt storage might not working
const IsIntKeyEncodedInBigEndian = true

type KeyValuePair struct {
	Key   []byte
	Value []byte
}

// Simple key/value storage interface
type Storage interface {
	Open() error
	Close() error
	Remove() error

	// Get all items from bucket
	GetAll(bucket string) ([]KeyValuePair, error)

	/*
	 * return error if buket or key is not found
	 */
	Get(bucket string, key []byte) ([]byte, error)

	/*
	 * if buket is not exist, create a new bucket
	 * if key not exist, create new key and value pair
	 * if key is already exist, replace the old value
	 * @attention: if key is nil or equals "AutoIncrementedId", then using auto-increment id as key
	 */
	Put(bucket string, key []byte, data []byte) error

	/*
	 * delete key/value pairs, if buket or key not exist, just do nothing
	 */
	Delete(bucket string, key []byte) error

	HasKey(bucket string, key []byte) bool

	CreateBucket(bucket string) error
	HasBucket(bucket string) bool
	RowCount(bucket string) int

	/* NOT USED
	 * delete a bucket, if not exist, do nothing
	 */
	//DeleteBucket(bucket string) error
}

func IntKeyToBytes(key uint64) []byte {
	b := make([]byte, 8)
	if IsIntKeyEncodedInBigEndian {
		binary.BigEndian.PutUint64(b, key)
	} else {
		binary.LittleEndian.PutUint64(b, key)
	}
	return b
}

func BytesToIntKey(data []byte) uint64 {
	var key uint64
	if IsIntKeyEncodedInBigEndian {
		key = binary.BigEndian.Uint64(data)
	} else {
		key = binary.LittleEndian.Uint64(data)
	}
	return key
}
