package statusdb

func (sdb *StatusDB) StartUndoSession() {
	sdb.reloadCachedValues()
	sdb.udb.StartUndoSession()
}

func (sdb *StatusDB) Commit() {
	sdb.udb.CommitLastSession()
	sdb.reloadCachedValues()
}

func (sdb *StatusDB) Undo() {
	sdb.udb.UndoLastSession()
	sdb.reloadCachedValues()
}

func (sdb *StatusDB) reloadCachedValues() {
	// no cached values so far
}
