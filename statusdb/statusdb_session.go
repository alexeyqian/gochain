package statusdb

func (sdb *StatusDB) StartUndoSession() {
	sdb.udb.StartUndoSession()
}

func (sdb *StatusDB) Commit() {
	sdb.udb.CommitLastSession()
}

func (sdb *StatusDB) Undo() {
	sdb.udb.UndoLastSession()
}
