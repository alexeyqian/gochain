package statusdb

func (sdb *StatusDB) HeadBlockNumber() uint64 {
	gpo, _ := sdb.GetGpo()
	return gpo.BlockNum
}
