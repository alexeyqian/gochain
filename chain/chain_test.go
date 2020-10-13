package chain

import (
	"fmt"
	"testing"
	"time"

	"github.com/alexeyqian/gochain/core"
)

func TestGenesis(t *testing.T) {
	Genesis()

	gpo := GetGpo()
	if gpo.BlockId == "" || gpo.BlockNum != 0 || gpo.Witness != InitWitness || gpo.CreatedOn <= 0 {
		t.Errorf("Gpo failed.")
	}

	accounts := GetAccounts()
	if len(accounts) != 1 || accounts[0].Name != InitWitness {
		t.Errorf("create account fail.")
	}

	acc := accounts[0]
	if acc.Id == "" || acc.Coin != 0 || acc.Credit != 0 {
		t.Errorf("init account error.")
	}

}

func TestGenerateBlock(t *testing.T) {
	Open()

	i := 0
	countx := 20
	for i < countx {
		var tx core.CreateAccountTransaction
		tx.AccountId = fmt.Sprintf("test_id_%d", i)
		tx.AccountName = fmt.Sprintf("test_account_%d", i)
		tx.CreatedOn = uint64(time.Now().Unix())
		tx.ExpiredOn = tx.CreatedOn + uint64(1000) + uint64(i)

		AddPendingTx(tx)
		BroadcastTx(tx)
		i++
	}

	b := GenerateBlock()
	if len(b.Transactions) != core.MaxTransactionsInBlock {
		t.Errorf("generate block error: transactions are not right, actual: %d", len(b.Transactions))
	}

	if len(_pendingTxs) != countx-core.MaxTransactionsInBlock {
		t.Errorf("generate block error: pending txs are not right")
	}

	Close()
}
