package chain

import (
	"time"

	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/eval"
	"github.com/alexeyqian/gochain/statusdb"
	"github.com/alexeyqian/gochain/utils"
)

func Open() {

}

func Close() {

}

func BroadcastTx(tx core.Transactioner) {
	// broadcast to other peers
}

func GenerateBlock() *core.Block {
	var b core.Block
	statusdb.MovePendingTxsToBlock(&b)
	applyBlock(&b)

	return &b
}

func applyBlock(b *core.Block) {
	for _, tx := range b.Transactions {
		eval.Apply(tx)
	}
}

func Genesis() {
	// create dummp block and push it to ledger
	// TODO

	// update global status
	bid := utils.CreateUuid()
	createdOn := uint64(time.Now().Unix())
	gpo := statusdb.GetGpo()
	gpo.BlockId = bid
	gpo.BlockNum = 0
	gpo.Witness = core.InitWitness
	gpo.CreatedOn = createdOn

	// update chain database
	var acc core.Account
	acc.Id = utils.CreateUuid() // should be public key string
	acc.Name = core.InitWitness
	acc.CreatedOn = createdOn
	statusdb.AddAccount(acc)
}
