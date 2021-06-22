package statusdb

func (sdb *StatusDB) HeadBlockNumber() int {
	gpo, _ := sdb.GetGpo()
	return gpo.BlockNum
}
