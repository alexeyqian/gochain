package store

import (
	"fmt"
)

type MemPair struct {
	Key  string
	Data []byte
}

type MemBucket struct {
	Pairs []MemPair
}

type MemoryStorage struct {
	buckets map[string]MemBucket
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		buckets: make(map[string]MemBucket),
	}
}

func (s *MemoryStorage) Open() {
}

func (s *MemoryStorage) Close() {
}

func (s *MemoryStorage) RemoveAll() {
	s.buckets = nil
}

func (s *MemoryStorage) GetAll(bucket string) ([][]byte, error) {
	if !s.HasBucket(bucket) {
		return nil, fmt.Errorf("bucket: %s not exist", bucket)
	}

	var res [][]byte
	for _, value := range s.buckets[bucket] {
		res = append(res, value)
	}
	return res, nil
}

func (s *MemoryStorage) Get(bucket, key string) ([]byte, error) {
	if !s.HasBucket(bucket) {
		return nil, fmt.Errorf("bucket: %s not exist", bucket)
	}
	return s.buckets[bucket][key], nil
}

func (s *MemoryStorage) Put(bucket, key string, data []byte) error {
	if !s.HasBucket(bucket) {
		s.buckets[bucket] = make([]MemPair)
	}

	s.buckets[bucket][key] = data
	return nil
}

func (s *MemoryStorage) HasBucket(bucket string) bool {
	_, ok := s.buckets[bucket]
	return ok
}
