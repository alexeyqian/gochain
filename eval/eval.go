package eval

import (
	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/statusdb"
)

func Apply(tx core.Transactioner) {
	t := tx.TxType()
	switch t {
	case core.CreateAccountTransactionType:
		ctx := tx.(core.CreateAccountTransaction)
		var acc core.Account
		acc.Id = ctx.AccountId
		acc.Name = ctx.AccountName

		statusdb.AddAccount(acc)
	}
}
