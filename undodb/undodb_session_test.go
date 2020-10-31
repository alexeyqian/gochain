package undodb

import (
	"fmt"
	"testing"

	"github.com/alexeyqian/gochain/entity"
)

type Book struct {
	ID            string
	Title         string
	Author        string
	PublishedYear int
	Price         float32
}

// func test db reopen

func TestSessionUndo(t *testing.T) {
	udb := NewTestDB()

	udb.Open()
	table := "book"
	udb.CreateTable(table)

	revision := udb.StartUndoSession()
	if revision != 1 {
		t.Errorf("revision expected: %d, actual: %d", 1, revision)
	}

	for i := 0; i < 10; i++ {
		book := Book{
			ID:            fmt.Sprintf("id_%d", i),
			Title:         fmt.Sprintf("title_%d", i),
			Author:        fmt.Sprintf("author_%d", i),
			PublishedYear: 1980 + i,
			Price:         100.00 + float32(i),
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

func TestSessionUndoExtra(t *testing.T) {
	udb := NewTestDB()

	udb.Open()
	table := "book"
	udb.CreateTable(table)

	for i := 0; i < 10; i++ {
		book := Book{
			ID:            fmt.Sprintf("id_%d", i),
			Title:         fmt.Sprintf("title_%d", i),
			Author:        fmt.Sprintf("author_%d", i),
			PublishedYear: 1980 + i,
			Price:         100.00 + float32(i),
		}

		udb.Create(table, book.ID, entity.Serialize(book))
	}

	revision := udb.StartUndoSession()
	if revision != 1 {
		t.Errorf("revision expected: %d, actual: %d", 1, revision)
	}

	rowCount := udb.RowCount(table)
	if rowCount != 10 {
		t.Errorf("row count expected: %d, actual: %d", 10, rowCount)
	}

	{
		// update book1
		id := fmt.Sprintf("id_%d", 1)
		var book Book
		data, _ := udb.Get(table, id)
		entity.Deserialize(&book, data)
		book.Title = "updated title"
		book.PublishedYear = 2000
		book.Price = 205.00

		udb.Update(table, book.ID, entity.Serialize(book))

		var bookUpdated Book
		data, _ = udb.Get(table, id)
		entity.Deserialize(&bookUpdated, data)

		if bookUpdated.Title != book.Title || bookUpdated.Price != book.Price {
			t.Errorf("book update failed")
		}

		// delete book2
		id = fmt.Sprintf("id_%d", 2)
		udb.Delete(table, id)
		_, err := udb.Get(table, id)
		if err == nil {
			t.Errorf("delete failed")
		}
	}

	udb.UndoLastSession()

	{
		// verify undo book1
		id := fmt.Sprintf("id_%d", 1)
		var book Book
		data, _ := udb.Get(table, id)
		entity.Deserialize(&book, data)
		//fmt.Printf("undo book: %+v\n", book)
		if book.Title != fmt.Sprintf("title_%d", 1) || book.PublishedYear != 1981 {
			t.Errorf("undo book update failed")
		}

		// verify undo deleted book2
		id = fmt.Sprintf("id_%d", 2)
		_, err := udb.Get(table, id)
		if err != nil {
			t.Errorf("undo delete failed")
		}

	}

	revision = udb.getCurrentRevision()
	if revision != 0 {
		t.Errorf("revision expected: %d, actual: %d", 0, revision)
	}

	rowCount = udb.RowCount(table)
	if rowCount != 10 {
		t.Errorf("row count expected: %d, actual: %d", 10, rowCount)
	}

	udb.Close()
	udb.Remove()
}

// commit session should empty revision table, and keep all data

func TestSessionCommit(t *testing.T) {
	udb := NewTestDB()

	udb.Open()
	table := "book"
	udb.CreateTable(table)

	revision := udb.StartUndoSession()
	if revision != 1 {
		t.Errorf("revision expected: %d, actual: %d", 1, revision)
	}

	for i := 0; i < 10; i++ {
		book := Book{
			ID:            fmt.Sprintf("id_%d", i),
			Title:         fmt.Sprintf("title_%d", i),
			Author:        fmt.Sprintf("author_%d", i),
			PublishedYear: 1980 + i,
			Price:         100.00 + float32(i),
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
