package undodb

import (
	"fmt"
	"testing"

	"github.com/alexeyqian/gochain/entity"
	"github.com/alexeyqian/gochain/store"
)

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

	// test update

	// test deletion

	udb.Close()
	udb.Remove()
}
