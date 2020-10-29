package undodb

import "fmt"

type RevisionRow struct {
	revision  uint32
	table     string
	operation string
	key       string
	data      []byte
}

type RevistionTable struct {
	rows []RevisionRow
}

type UndoableDB struct {
	storage        *Storage
	revisionTable  *RevistionTable
	dataBucket     string
	stateBucket    string
	latestRevision uint32
	count          uint32
}

func NewUndoableSet(s Storage, name, elementType string) (*UndoableSet, error) {
	if s.HasBucket(name) {
		return nil, fmt.Errorf("cannot create same set in storage")
	}

	us := UndoableSet{
		storage:        s,
		dataBucket:     name + "_data",
		stateBucket:    name + "_state",
		latestRevision: 0,
		count:          0,
	}

	// create data bucket and state bucket
	err := s.CreateBucket(us.dataBucket)
	if err != nil {
		return nil, err
	}
	err = s.CreateBucket(us.stateBucket)
	if err != nil {
		return nil, err
	}

	return &us, nil
}
