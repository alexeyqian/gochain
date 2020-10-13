package main

import (
	"fmt"
	"testing"
	"time"

	core "github.com/alexeyqian/gochain/core"
	ledger "github.com/alexeyqian/gochain/ledger"
	utils "github.com/alexeyqian/gochain/utils"
)

func TestLedgerOpen(t *testing.T) {
	ledger.Open("test_data")
	sec := time.Now().Unix()
	id := utils.CreateUuid()
	b := core.Block{Id: id, Num: 0, CreatedOn: uint64(sec), Witness: "init_miner"}
	ledger.Append(core.SerializeBlock(&b))
	br := ledger.Read(0)
	block := core.UnSerializeBlock(br)
	fmt.Printf("%+v\n", block)
	ledger.Close()
	ledger.Remove()
	if block.Id != id {
		t.Errorf("id is: %s, want: %s", block.Id, id)
	}
}
