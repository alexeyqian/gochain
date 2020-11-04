package chain

import (
	"fmt"
	"testing"

	"github.com/alexeyqian/gochain/core"
)

func TestGenesis(t *testing.T) {
	c := SetupTestChain()

	gpo, err := c.sdb.GetGpo()
	if err != nil {
		panic(err)
	}
	if gpo.BlockId == "" || gpo.BlockNum != 0 || gpo.Witness != core.InitWitness || gpo.Time <= 0 {
		t.Errorf("Gpo failed.")
	}
	if gpo.Supply != core.InitAmount {
		t.Errorf("genesis supply expected: %d, actual: %d", core.InitAmount, gpo.Supply)
	}

	accounts := c.sdb.GetAccounts()
	fmt.Printf("account len: %d", len(accounts))
	if len(accounts) != 1 || accounts[0].Name != core.InitWitness {
		t.Errorf("create account fail.")
	}

	acc := accounts[0]
	if acc.ID == "" || acc.Coin != core.InitAmount || acc.Credit != 0 {
		t.Errorf("init account error.")
	}

	b, _ := c.GetBlock(0)
	if b.Num != 0 {
		t.Errorf("genesis zero block not generated")
	}

	TearDownTestChain(c)
}

func TestGenerateBlock(t *testing.T) {
	c := SetupTestChain()

	i := 0
	countx := 20
	for i < countx {
		tx := CreateTestAccount(fmt.Sprintf("test_account_name_%d", i))
		c.AddPendingTx(tx)
		i++
	}

	tempb := c.GenerateBlock()
	b, _ := c.GetBlock(int(tempb.Num))

	gpo, _ := c.sdb.GetGpo()
	if gpo.BlockNum != 1 || gpo.BlockId != b.ID {
		t.Errorf("generate block error")
	}

	if gpo.Supply != core.InitAmount+core.AmountPerBlock {
		t.Errorf("generate block gpo amount expected: %d, actual: %d", core.InitAmount+core.AmountPerBlock, gpo.Supply)
	}

	if len(b.Transactions) != core.MaxTransactionsInBlock {
		t.Errorf("generate block error: transactions are not right, actual: %d", len(b.Transactions))
	}

	if len(c.GetPendingTx()) != countx-core.MaxTransactionsInBlock {
		t.Errorf("generate block error: pending txs are not right")
	}

	i = 0
	for i < core.MaxTransactionsInBlock {
		accName := fmt.Sprintf("test_account_name_%d", i)
		acc, _ := c.sdb.GetAccountByName(accName)
		if acc == nil || acc.Name != accName {
			t.Errorf("cannot find account name: %s", accName)
		}
		i++
	}

	TearDownTestChain(c)
}

func TestGenerateBlocks(t *testing.T) {
	c := SetupTestChain()

	i := 1
	for i <= 20 {
		tx := CreateTestAccount(fmt.Sprintf("test_account_name_%d", i))
		c.AddPendingTx(tx)
		//chain.BroadcastTx(tx)
		b := c.GenerateBlock()
		if b.Num != uint64(i) {
			t.Errorf("expected: %d, actual: %d", i, b.Num)
		}

		gpo, _ := c.sdb.GetGpo()
		if gpo.BlockNum != b.Num {
			t.Errorf("gpo num expected: %d, actual: %d", 20, gpo.BlockNum)
		}
		if gpo.BlockId != b.ID {
			t.Errorf("gpo id expected: %s, actual: %s", b.ID, gpo.BlockId)
		}
		if gpo.Time != b.CreatedOn {
			t.Errorf("gpo time expected: %d, actual: %d", b.CreatedOn, gpo.Time)
		}
		if gpo.Supply != (core.InitAmount + core.AmountPerBlock*uint64(i)) {
			t.Errorf("generate block gpo amount expected: %d, actual: %d", core.InitAmount+core.AmountPerBlock*i, gpo.Supply)
		}

		// TODO: validate block and previous block hash/linking
		prevb, _ := c.GetBlock(i - 1)
		//fmt.Printf("prevb id: %s", prevb.ID)
		if b.PrevBlockId != prevb.ID {
			t.Errorf("block linking is broken")
		}

		i++
	}

	TearDownTestChain(c)
}
