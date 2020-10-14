package tests

import (
	"testing"

	"github.com/alexeyqian/gochain/chain"
	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/statusdb"
)

func TestCreateAccount(t *testing.T) {
	chain.Open(TestDataDir)

	statusdb.AddPendingTx(TstCreateAccount("alice"))
	chain.GenerateBlock()

	acc := statusdb.GetAccount("alice")
	if acc.Name != "alice" {
		t.Errorf("create account faile")
	}

	chain.Close()
	chain.Remove()
}

func TestTransfer(t *testing.T) {
	chain.Open(TestDataDir)

	statusdb.AddPendingTx(TstCreateAccount("alice"))
	chain.GenerateBlock()

	var tx core.TransferCoinTransaction
	tx.From = "init"
	tx.To = "alice"
	tx.Amount = 100
	statusdb.AddPendingTx(tx)
	chain.GenerateBlock()

	acc := statusdb.GetAccount("alice")
	if acc.Coin != 100 {
		t.Errorf("expected: %d, actual: %d", 100, acc.Coin)
	}

	chain.Close()
	chain.Remove()
}
