package core

import "testing"

func TestAddTransactions(t *testing.T) {
	var b Block
	var tx CreateAccountTransaction
	b.AddTransaction(tx)
	if len(b.Transactions) != 1 {
		t.Errorf("add transaction failed.")
	}
}
