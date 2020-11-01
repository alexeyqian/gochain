package tests

import (
	"testing"
	"time"

	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/ledger"
	"github.com/alexeyqian/gochain/utils"
)

func TestLedgerAppendSingle(t *testing.T) {
	ledger.Open(TestDataDir)

	sec := time.Now().Unix()
	id := utils.CreateUuid()
	b := core.Block{ID: id, Num: 0, CreatedOn: uint64(sec), Witness: "init_miner"}
	sb, _ := core.SerializeBlock(&b)
	ledger.Append(sb)

	br, _ := ledger.Read(0)
	block, _ := core.UnSerializeBlock(br)

	if block.ID != id {
		t.Errorf("id is: %s, want: %s", block.ID, id)
	}

	ledger.Close()
	ledger.Remove()

	//fmt.Printf("%+v\n", block)
}

func TestLedgerAppendMultiple(t *testing.T) {
	ledger.Open(TestDataDir)

	i := 0
	sec := time.Now().Unix()
	for i < 10 {
		id := utils.CreateUuid()
		createdOn := uint64(sec) + uint64(i)*uint64(core.BlockInterval)
		b := core.Block{ID: id, Num: uint64(i), CreatedOn: createdOn, Witness: "init_miner"}
		sb, _ := core.SerializeBlock(&b)
		ledger.Append(sb)
		br, _ := ledger.Read(i)
		block, _ := core.UnSerializeBlock(br)

		if block.ID != id {
			t.Errorf("id is: %s, want: %s", block.ID, id)
		}

		//fmt.Printf("%+v\n", block)
		i++
	}

	ledger.Close()
	ledger.Remove()
}
