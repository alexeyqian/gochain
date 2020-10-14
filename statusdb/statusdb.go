package statusdb

import (
	"github.com/alexeyqian/gochain/core"
)

var _gpo core.Gpo
var _accounts []core.Account
var _pendingTxs []core.Transactioner

func GetGpo() *core.Gpo {
	return &_gpo
}

func GetPendingTxs() []core.Transactioner {
	return _pendingTxs
}

func AddPendingTx(tx core.Transactioner) {
	_pendingTxs = append(_pendingTxs, tx)
}

func MovePendingTxsToBlock(b *core.Block) {
	i := 0
	for _, tx := range _pendingTxs {
		if i >= core.MaxTransactionsInBlock {
			break
		}
		b.AddTransaction(tx)
		i++
	}

	if len(_pendingTxs) > core.MaxTransactionsInBlock {
		_pendingTxs = _pendingTxs[core.MaxTransactionsInBlock:]
	}
}

func AddAccount(acc core.Account) {
	_accounts = append(_accounts, acc)
}

func GetAccounts() []core.Account {
	return _accounts
}

func GetAccount(name string) *core.Account {
	for _, acc := range _accounts {
		if acc.Name == name {
			return &acc
		}
	}
	return nil
}
