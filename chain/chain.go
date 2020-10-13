package chain

import (
	"time"

	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/utils"
)

const InitWitness = "init"

type Gpo struct {
	BlockId   string
	BlockNum  int
	Witness   string
	CreatedOn uint64
	Version   string
	Supply    uint64
}

var _gpo Gpo
var _accounts []core.Account
var _pendingTxs []core.Transactioner

func Open() {

}

func Close() {

}

func AddPendingTx(tx core.Transactioner) {
	_pendingTxs = append(_pendingTxs, tx)
}

func BroadcastTx(tx core.Transactioner) {
	// broadcast to other peers
}

func GenerateBlock() *core.Block {
	var b core.Block
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

	return &b
}

func Genesis() {
	// create dummp block and push it to ledger
	// TODO

	// update global status
	bid := utils.CreateUuid()
	createdOn := uint64(time.Now().Unix())
	_gpo.BlockId = bid
	_gpo.BlockNum = 0
	_gpo.Witness = InitWitness
	_gpo.CreatedOn = createdOn

	// update chain database
	var acc core.Account
	acc.Id = utils.CreateUuid() // should be public key string
	acc.Name = InitWitness
	acc.CreatedOn = createdOn
	_accounts = append(_accounts, acc)
}

func GetGpo() *Gpo {
	return &_gpo
}

func GetAccounts() []core.Account {
	return _accounts
}
