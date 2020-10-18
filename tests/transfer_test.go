package tests

import (
	"testing"

	"github.com/alexeyqian/gochain/chain"
	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/statusdb"
)

func TestCreateAccount(t *testing.T) {
	chain.Open(TestDataDir)

	chain.AddPendingTx(CreateTestAccount("alice"))
	chain.GenerateBlock()

	acc, _ := statusdb.GetAccountByName("alice")
	if acc.Name != "alice" {
		t.Errorf("create account faile")
	}

	chain.Close()
	chain.Remove()
}

func TestTransfer(t *testing.T) {
	chain.Open(TestDataDir)

	chain.AddPendingTx(CreateTestAccount("alice"))
	chain.GenerateBlock()

	var tx core.TransferCoinTransaction
	tx.From = "init"
	tx.To = "alice"
	tx.Amount = 100
	chain.AddPendingTx(tx)
	chain.GenerateBlock()

	acc, _ := statusdb.GetAccountByName("alice")
	if acc.Coin != 100 {
		t.Errorf("expected: %d, actual: %d", 100, acc.Coin)
	}

	chain.Close()
	chain.Remove()
}

func TestInvalidTxWillBeIgnored(t *testing.T) {
	chain.Open(TestDataDir)

	var tx core.CreateAccountTransaction
	tx.AccountName = ""
	err := chain.AddPendingTx(tx)

	if err == nil {
		t.Errorf("cannot detect invalid tx")
	}

	if len(chain.GetPendingTx()) != 0 {
		t.Errorf("should not add invalid tx to pending tx list")
	}

	chain.Close()
	chain.Remove()
}

func TestTransferNoSufficientFund(t *testing.T) {
	chain.Open(TestDataDir)

	chain.AddPendingTx(CreateTestAccount("alice"))
	chain.GenerateBlock()

	var tx core.TransferCoinTransaction
	tx.From = "init"
	tx.To = "alice"
	tx.Amount = 10000
	err := chain.AddPendingTx(tx)

	if err == nil {
		t.Errorf("transfer not sufficent found, cannot detect error")
	}

	accInit, _ := statusdb.GetAccountByName(core.InitWitness)
	if accInit.Coin != core.InitAmount {
		t.Errorf("expected: %d, actual: %d", core.InitAmount, accInit.Coin)
	}

	accAlice, _ := statusdb.GetAccountByName("alice")
	if accAlice.Coin != 0 {
		t.Errorf("expected: %d, actual: %d", 0, accAlice.Coin)
	}

	chain.Close()
	chain.Remove()
}
