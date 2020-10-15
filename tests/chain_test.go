package tests

import (
	"fmt"
	"testing"

	"github.com/alexeyqian/gochain/chain"
	core "github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/statusdb"
)

func TestGenesis(t *testing.T) {
	chain.Open(TestDataDir)

	gpo := statusdb.GetGpo()
	if gpo.BlockId == "" || gpo.BlockNum != 0 || gpo.Witness != core.InitWitness || gpo.Time <= 0 {
		t.Errorf("Gpo failed.")
	}

	if gpo.Supply != core.InitAmount {
		t.Errorf("genesis supply expected: %d, actual: %d", core.InitAmount, gpo.Supply)
	}

	accounts := statusdb.GetAccounts()
	if len(accounts) != 1 || accounts[0].Name != core.InitWitness {
		t.Errorf("create account fail.")
	}

	acc := accounts[0]
	if acc.Id == "" || acc.Coin != core.InitAmount || acc.Credit != 0 {
		t.Errorf("init account error.")
	}

	b, _ := chain.GetBlock(0)
	if b.Num != 0 {
		t.Errorf("genesis zero block not generated")
	}

	chain.Close()
	chain.Remove()
}

func TestGenerateBlock(t *testing.T) {
	chain.Open(TestDataDir)

	i := 0
	countx := 20
	for i < countx {
		tx := TstCreateAccount(fmt.Sprintf("test_account_name_%d", i))
		chain.AddPendingTx(tx)
		//chain.BroadcastTx(tx)
		i++
	}

	tempb := chain.GenerateBlock()
	b, _ := chain.GetBlock(int(tempb.Num))

	gpo := statusdb.GetGpo()
	if gpo.BlockNum != 1 || gpo.BlockId != b.Id {
		t.Errorf("generate block error")
	}

	if gpo.Supply != core.InitAmount+core.AmountPerBlock {
		t.Errorf("generate block gpo amount expected: %d, actual: %d", core.InitAmount+core.AmountPerBlock, gpo.Supply)
	}

	if len(b.Transactions) != core.MaxTransactionsInBlock {
		t.Errorf("generate block error: transactions are not right, actual: %d", len(b.Transactions))
	}

	if len(statusdb.GetPendingTransactions()) != countx-core.MaxTransactionsInBlock {
		t.Errorf("generate block error: pending txs are not right")
	}

	i = 0
	for i < core.MaxTransactionsInBlock {
		accName := fmt.Sprintf("test_account_name_%d", i)
		acc := statusdb.GetAccount(accName)
		if acc == nil || acc.Name != accName {
			t.Errorf("cannot find account name: %s", accName)
		}
		i++
	}

	chain.Close()
	chain.Remove()
}

func TestGenerateBlocks(t *testing.T) {
	chain.Open(TestDataDir)

	i := 1
	for i <= 20 {
		tx := TstCreateAccount(fmt.Sprintf("test_account_name_%d", i))
		chain.AddPendingTx(tx)
		//chain.BroadcastTx(tx)
		b := chain.GenerateBlock()
		if b.Num != uint64(i) {
			t.Errorf("expected: %d, actual: %d", i, b.Num)
		}

		gpo := statusdb.GetGpo()
		if gpo.BlockNum != b.Num {
			t.Errorf("gpo num expected: %d, actual: %d", 20, gpo.BlockNum)
		}
		if gpo.BlockId != b.Id {
			t.Errorf("gpo id expected: %s, actual: %s", b.Id, gpo.BlockId)
		}
		if gpo.Time != b.CreatedOn {
			t.Errorf("gpo time expected: %d, actual: %d", b.CreatedOn, gpo.Time)
		}
		if gpo.Supply != (core.InitAmount + core.AmountPerBlock*uint64(i)) {
			t.Errorf("generate block gpo amount expected: %d, actual: %d", core.InitAmount+core.AmountPerBlock*i, gpo.Supply)
		}

		// TODO: validate block and previous block hash/linking
		prevb, _ := chain.GetBlock(i - 1)
		//fmt.Printf("prevb id: %s", prevb.Id)
		if b.PrevBlockId != prevb.Id {
			t.Errorf("block linking is broken")
		}

		i++
	}

	chain.Close()
	chain.Remove()
}
