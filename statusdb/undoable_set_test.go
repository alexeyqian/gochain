package statusdb

import "testing"

type Book struct {
	Title   string
	Author  string
	Content string
}

func NewBook(title, author, content string) *Book {
	return &Book{
		Title:   title,
		Author:  author,
		Content: content,
	}
}

type Storage interface {
}

func TestUndo(t *testing.T) {
	var storage Storage
	storage = NewBoltStorage("test_undo.db") // or &MemoryStorage{}

	s := NewUndoableSet(storage, "book", Book)

	s.StartUndoSession()
	for i := 0; i < 10; i++ {
		title := fmt.Sprintf("book_%d", i)
		author := fmt.Strintf("author_%d", i)
		content := fmt.Sprintf("content_%d", i)
		s.Create(NewBook(title, author, content))
	}

	if s.Size() != 10 {
		t.Errorf("undoable set: data not saved")
	}

	if s.ReadRevision() != 1 {
		t.Errorf("undoable set: read revision error, expected: %d, actual:%d", 1, s.ReadRevision())
	}

	//if s.ReadNextId() != 1 {
	//		t.Errorf("undoable set: read revision error, expected: %d, actual:%d", 1, s.ReadRevision())
	//}

	s.UndoChangesFromLastSession()
	if s.Size() != 0 {
		t.Errorf("undoable set: size undo failed")
	}

	if s.ReadRevision() != 0 {
		t.Errorf("undoable set: revision undo error, expected: %d, actual:%d", 0, s.ReadRevision())
	}
}
