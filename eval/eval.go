package eval

import (
	"fmt"

	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/statusdb"
)

func Apply(tx core.Transactioner) {
	txType := tx.TxType()
	switch txType {
	case core.CreateAccountTransactionType:
		cat := tx.(core.CreateAccountTransaction)
		var acc core.Account
		acc.Id = cat.AccountId
		acc.Name = cat.AccountName

		statusdb.AddAccount(acc)
	case core.TransferCoinTransactionType:
		tct := tx.(core.TransferCoinTransaction)
		fmt.Printf("from: %s, to: %s\n", tct.From, tct.To)
		fromAcc := statusdb.GetAccount(tct.From)
		toAcc := statusdb.GetAccount(tct.To)

		if fromAcc == nil {
			panic("from account is not exist")
		}
		if toAcc == nil {
			panic("to account is not exist")
		}
		if fromAcc.Coin < tct.Amount {
			panic("no enough coin")
		}
		fromAcc.Coin -= tct.Amount
		toAcc.Coin += tct.Amount
	}
}
