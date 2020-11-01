package statusdb

import "testing"

func TestCRUD(t *testing.T) {
	var storage Storage
	storage = NewBoltStorage("test_undo.db") // or &MemoryStorage{}
	sdb := NewStatusDB(storage)
	sdb.Open()

	for i := 0; i < 10; i++ {
		id := fmt.Sprintf("id_%d", i)
		name := fmt.Strintf("name_%d", i)
		s.CreateAccount(&Account{ID: id, Name: name})
	}

	if sdb.AccountSet().Size() != 10 {
		t.Errorf("create account failed")
	}

	if sdb.Revision() != 1 {
		t.Errorf("revision error, expected: %d, actual:%d", 1, s.Revision())
	}

	sdb.Close()
	sbd.Remove()
}

func TestUndo(t *testing.T) {
	var storage Storage
	storage = NewBoltStorage("test_undo.db") // or &MemoryStorage{}
	sdb := NewStatusDB(storage)
	sdb.Open()

	sdb.StartUndoSession()
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

	sdb.Close()
	sdb.Remove()
}
