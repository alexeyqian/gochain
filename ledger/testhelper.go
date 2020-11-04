package ledger

import "os"

const testDataDir = "test_data"

func SetupLedger() Ledger {
	lgr := NewFileLedger(testDataDir)
	lgr.Open()
	return lgr
}

func TearDownLedger(lgr Ledger) {
	lgr.Close()
	lgr.Remove()
	os.Remove(testDataDir)
}
