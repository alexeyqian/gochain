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

func TestCreate(t *testing.T) {
	pathname := "test.db"
	storage := store.NewBoltStorage(pathname)
	udb := NewUndoableDB(storage)

	err := udb.Open()
	if err != nil {
		t.Errorf("cannot open")
	}

	if !udb.HasTable(metaTable) {
		t.Errorf("meta table is not created")
	}
	if !udb.HasTable(revisionTable) {
		t.Errorf("revision table is not created")
	}

	metaData := udb.getMetaData()
	if metaData.Revision != 0 {
		t.Errorf("revision is not correct")
	}

	table := "book"
	err = udb.CreateTable(table)
	if err != nil || !udb.HasTable(table) {
		t.Errorf("account table is not created")
	}

	for i := 0; i < 10; i++ {
		book := Book{
			ID:           fmt.Sprintf("id_%d", i),
			Title:        fmt.Sprintf("title_%d", i),
			Author:       fmt.Sprintf("author_%d", i),
			PublishedYer: 1980 + i,
			Price:        100.00 + float32(i),
		}

		err = udb.Create(table, book.ID, entity.Serialize(book))
		if err != nil {
			t.Error(err)
		}
	}

	rowCount := udb.RowCount(table)
	if rowCount != 10 {
		t.Errorf("row count expected: %d, actual: %d", 10, rowCount)
	}

	for i := 0; i < 10; i++ {
		data, err := udb.Get(table, fmt.Sprintf("id_%d", i))
		if err != nil {
			t.Errorf("get book: %d error", i)
		}
		var book Book
		entity.Deserialize(&book, data)
		if book.ID != fmt.Sprintf("id_%d", i) {
			t.Errorf("book id error")
		}
		if book.Title != fmt.Sprintf("title_%d", i) {
			t.Errorf("book title error")
		}

		//fmt.Printf("got book %d is %+v\n", i, book)
	}

	udb.Close()
	udb.Remove()
}

// func test db reopen

func TestUndoSession(t *testing.T) {
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
