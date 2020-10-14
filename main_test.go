package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/alexeyqian/gochain/chain"
	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/ledger"
	"github.com/alexeyqian/gochain/statusdb"
	"github.com/alexeyqian/gochain/utils"
)

func TestLedgerAppend(t *testing.T) {
	ledger.Open("test_data")

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

func TestGenesis(t *testing.T) {
	chain.Genesis()

	gpo := statusdb.GetGpo()
	if gpo.BlockId == "" || gpo.BlockNum != 0 || gpo.Witness != core.InitWitness || gpo.CreatedOn <= 0 {
		t.Errorf("Gpo failed.")
	}

	accounts := statusdb.GetAccounts()
	if len(accounts) != 1 || accounts[0].Name != core.InitWitness {
		t.Errorf("create account fail.")
	}

	acc := accounts[0]
	if acc.Id == "" || acc.Coin != 0 || acc.Credit != 0 {
		t.Errorf("init account error.")
	}

}
