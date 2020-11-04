package ledger

import "os"

const testDataDir = "test_data"

func SetupLedger() {
	lgr := NewFileLedger(testDataDir)
	lgr.Open()
	return lgr
}

func TearDownLedger(lgr Ledger) {
	lgr.Close()
	lgr.Remove()
	os.Remove(testDataDir)
}
