package store

import (
	"fmt"
)

type MemKey []byte
type MemValue []byte
type MemBucket map[string]MemValue

// or MmeBucket{[]KeyValuePair}

type MemoryStorage struct {
	pathname  string // not used
	buckets   map[string]MemBucket
	sequences map[string]int
}

func NewMemoryStorage(name string) *MemoryStorage {
	return &MemoryStorage{
		pathname:  name,
		buckets:   make(map[string]MemBucket),
		sequences: make(map[string]int),
	}
}

func (s *MemoryStorage) Open() error {
	return nil
}

func (s *MemoryStorage) Close() error {
	return nil
}

func (s *MemoryStorage) Remove() error {
	s.buckets = nil
	return nil
}

func (s *MemoryStorage) Get(bucket string, key []byte) ([]byte, error) {
	if !s.HasBucket(bucket) {
		return nil, fmt.Errorf("bucket: %s not exist", bucket)
	}
	data, ok := s.buckets[bucket][string(key)]
	if len(data) > 0 && ok { // ?? some times returns ok == true, but return []
		return data, nil
	} else {
		return nil, fmt.Errorf("not found!")
	}
}

func (s *MemoryStorage) Put(bucket string, key []byte, data []byte) error {
	if !s.HasBucket(bucket) {
		s.buckets[bucket] = make(map[string]MemValue)
	}

	if key == nil || string(key) == AutoIncrementKey {
		id := s.NextSequence(bucket) // return int, error
		key = IntKeyToBytes(id)
	}

	s.buckets[bucket][string(key)] = data
	return nil
}

// autokey is int
func (s *MemoryStorage) PutWithAutoKey(bucket string, data []byte) error {
	return s.Put(bucket, nil, data)
}

func (s *MemoryStorage) Delete(bucket string, key []byte) error {
	if !s.HasBucket(bucket) {
		return nil
	}

	delete(s.buckets[bucket], string(key))

	return nil
}

func (s *MemoryStorage) GetAll(bucket string) ([]KeyValuePair, error) {
	if !s.HasBucket(bucket) {
		return nil, fmt.Errorf("bucket: %s not exist", bucket)
	}

	var res []KeyValuePair
	for k, v := range s.buckets[bucket] {
		res = append(res, KeyValuePair{Key: []byte(k), Value: v})
	}
	return res, nil
}

func (s *MemoryStorage) HasKey(bucket string, key []byte) bool {
	if !s.HasBucket(bucket) {
		return false
	}

	data := s.buckets[bucket][string(key)]
	return data != nil && len(data) > 0
}

func (s *MemoryStorage) HasBucket(bucket string) bool {
	_, ok := s.buckets[bucket]
	return ok
}

func (s *MemoryStorage) CreateBucket(bucket string) error {
	s.buckets[bucket] = make(map[string]MemValue)
	return nil
}

func (s *MemoryStorage) DeleteBucket(bucket string) error {
	s.buckets[bucket] = nil
	return nil
}

func (s *MemoryStorage) RowCount(bucket string) int {
	if !s.HasBucket(bucket) {
		return 0
	}

	return len(s.buckets[bucket])
}

func (s *MemoryStorage) NextSequence(bucket string) int {
	res := s.sequences[bucket]
	res++
	s.sequences[bucket] = res
	return res
}
