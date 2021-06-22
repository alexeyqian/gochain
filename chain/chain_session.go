package chain

func (c *Chain) startUndoSession() {
	c.reloadCachedValues()
	c.sdb.StartUndoSession()
}

func (c *Chain) pushUndoSession() {
	c.reloadCachedValues()
	// nothing more need to be done
	// just leave the undo session/revision inside undoable db
}

func (c *Chain) commit() {
	c.sdb.Commit()
	c.reloadCachedValues()
}

// undo last session/revision
func (c *Chain) undo() {
	c.sdb.Undo()
	c.reloadCachedValues() // all undo operations should reload.
}

func (c *Chain) reloadCachedValues() {
	c.reloadHead()
}
