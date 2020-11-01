package undodb

import (
	"fmt"
	"testing"

	"github.com/alexeyqian/gochain/entity"
	"github.com/alexeyqian/gochain/store"
)

func NewTestDB() *UndoableDB {
	pathname := "test.db"
	storage := store.NewMemoryStorage(pathname) //store.NewBoltStorage(pathname)
	udb := NewUndoableDB(storage)
	return udb
}

func TestCreate(t *testing.T) {
	udb := NewTestDB()

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
			ID:            fmt.Sprintf("id_%d", i),
			Title:         fmt.Sprintf("title_%d", i),
			Author:        fmt.Sprintf("author_%d", i),
			PublishedYear: 1980 + i,
			Price:         100.00 + float32(i),
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

	id := fmt.Sprintf("id_%d", 1)

	// test update
	{
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
		//fmt.Printf("updated book: %+v\n", bookUpdated)
	}

	// test deletion
	{
		udb.Delete(table, id)
		_, err := udb.Get(table, id)
		if err == nil {
			t.Errorf("delete failed")
		}

		count := udb.RowCount(table)
		//fmt.Printf("row count: %d\n", count)
		if count != 9 {
			t.Errorf("[delete] expected: %d, actual: %d", 9, count)
		}
	}

	udb.Close()
	udb.Remove()
}
