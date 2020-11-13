package chain

func (c *Chain) startUndoSession() {
	c.sdb.StartUndoSession()
}

func (c *Chain) pushUndoSession() {
	// nothing need to be done here
	// just leave the undo session/revision inside undoable db
}

func (c *Chain) commit() {
	c.sdb.Commit()
}

// undo last session/revision
func (c *Chain) undo() {
	c.sdb.Undo()
}
