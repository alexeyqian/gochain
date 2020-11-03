package ledger

import (
	"os"
	"testing"
	"time"

	"github.com/alexeyqian/gochain/utils"
)

const testDataDir = "test_data"

type XBlock struct {
	ID        string
	Num       int
	CreatedOn int64
	Witness   string
}

func TestLedgerAppendSingle(t *testing.T) {
	lgr := NewLedger()
	lgr.Open(testDataDir)

	id := utils.CreateUuid()
	b := XBlock{ID: id, Num: 0, CreatedOn: time.Now().Unix(), Witness: "init_miner"}
	lgr.Append(utils.Serialize(b))

	bdata, _ := lgr.Read(0)
	var b2 XBlock
	utils.Deserialize(&b2, bdata)

	if b2.ID != id {
		t.Errorf("id is: %s, want: %s", b2.ID, id)
	}

	lgr.Close()
	lgr.Remove()
	os.Remove(testDataDir)

	//fmt.Printf("%+v\n", block)
}

func TestLedgerAppendMultiple(t *testing.T) {
	lgr := NewLedger()
	lgr.Open(testDataDir)

	i := 0
	sec := time.Now().Unix()
	for i < 10 {
		id := utils.CreateUuid()
		createdOn := sec + int64(i)*int64(3)
		b := XBlock{ID: id, Num: i, CreatedOn: createdOn, Witness: "init_miner"}
		lgr.Append(utils.Serialize(b))
		bdata, _ := lgr.Read(i)
		var b2 XBlock
		utils.Deserialize(&b2, bdata)

		if b2.ID != id {
			t.Errorf("id is: %s, want: %s", b2.ID, id)
		}

		//fmt.Printf("%+v\n", block)
		i++
	}

	lgr.Close()
	lgr.Remove()
	os.Remove(testDataDir)
}
