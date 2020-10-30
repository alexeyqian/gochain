package undodb

import (
	"fmt"
	"testing"

	"github.com/alexeyqian/gochain/entity"

	"github.com/alexeyqian/gochain/store"
)

type Book struct {
	ID           string
	Title        string
	Author       string
	PublishedYer int
	Price        float32
}

// func test db reopen

func TestSessionUndo(t *testing.T) {
	pathname := "test.db"
	storage := store.NewBoltStorage(pathname)
	udb := NewUndoableDB(storage)

	udb.Open()
	table := "book"
	udb.CreateTable(table)

	revision := udb.StartUndoSession()
	if revision != 1 {
		t.Errorf("revision expected: %d, actual: %d", 1, revision)
	}

	for i := 0; i < 10; i++ {
		book := Book{
			ID:           fmt.Sprintf("id_%d", i),
			Title:        fmt.Sprintf("title_%d", i),
			Author:       fmt.Sprintf("author_%d", i),
			PublishedYer: 1980 + i,
			Price:        100.00 + float32(i),
		}

		udb.Create(table, book.ID, entity.Serialize(book))
	}

	rowCount := udb.RowCount(table)
	if rowCount != 10 {
		t.Errorf("row count expected: %d, actual: %d", 10, rowCount)
	}

	udb.UndoLastSession()

	revision = udb.getCurrentRevision()
	if revision != 0 {
		t.Errorf("revision expected: %d, actual: %d", 0, revision)
	}

	rowCount = udb.RowCount(table)
	if rowCount != 0 {
		t.Errorf("row count expected: %d, actual: %d", 0, rowCount)
	}

	udb.Close()
	udb.Remove()
}

// commit session should empty revision table, and keep all data

func TestSessionCommit(t *testing.T) {
	pathname := "test.db"
	storage := store.NewBoltStorage(pathname)
	udb := NewUndoableDB(storage)

	udb.Open()
	table := "book"
	udb.CreateTable(table)

	revision := udb.StartUndoSession()
	if revision != 1 {
		t.Errorf("revision expected: %d, actual: %d", 1, revision)
	}

	for i := 0; i < 10; i++ {
		book := Book{
			ID:           fmt.Sprintf("id_%d", i),
			Title:        fmt.Sprintf("title_%d", i),
			Author:       fmt.Sprintf("author_%d", i),
			PublishedYer: 1980 + i,
			Price:        100.00 + float32(i),
		}

		udb.Create(table, book.ID, entity.Serialize(book))
	}

	rowCount := udb.RowCount(table)
	if rowCount != 10 {
		t.Errorf("row count expected: %d, actual: %d", 10, rowCount)
	}

	udb.CommitLastSession()

	revision = udb.getCurrentRevision()
	if revision != 0 {
		t.Errorf("revision expected: %d, actual: %d", 0, revision)
	}

	rowCount = udb.RowCount(table)
	if rowCount != 10 {
		t.Errorf("book row count expected: %d, actual: %d", 10, rowCount)
	}

	rowCount = udb.RowCount(revisionTable)
	if rowCount != 0 {
		t.Errorf("revision row count expected: %d, actual: %d", 0, rowCount)
	}

	for i := 0; i < 10; i++ {
		bookdata, _ := udb.Get(table, fmt.Sprintf("id_%d", i))
		var book Book
		entity.Deserialize(&book, bookdata)
		//fmt.Printf("commited book: %+v\n", book)
	}

	udb.Close()
	udb.Remove()
}
