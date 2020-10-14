package tests

import (
	"fmt"
	"testing"
	"time"

	core "github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/ledger"
	utils "github.com/alexeyqian/gochain/utils"
)

func TestLedgerAppend(t *testing.T) {
	ledger.Open(TestDataDir)

	sec := time.Now().Unix()
	id := utils.CreateUuid()
	b := core.Block{Id: id, Num: 0, CreatedOn: uint64(sec), Witness: "init_miner"}
	ledger.Append(core.SerializeBlock(&b))

	br := ledger.Read(0)
	block := core.UnSerializeBlock(br)

	if block.Id != id {
		t.Errorf("id is: %s, want: %s", block.Id, id)
	}

	ledger.Close()
	ledger.Remove()

	fmt.Printf("%+v\n", block)
}
