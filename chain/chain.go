package chain

import (
	"time"

	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/eval"
	"github.com/alexeyqian/gochain/ledger"
	"github.com/alexeyqian/gochain/statusdb"
	"github.com/alexeyqian/gochain/utils"
)

var _dataFolder string = "data"
var _isGenesised = false

func Open(dir string) {
	_dataFolder = dir
	ledger.Open(dir)
	// statusdb.Open()
	if !_isGenesised {
		genesis()
	}
}

func Close() {
	ledger.Close()
	//statusdb.Close()
}

func Remove() {
	ledger.Remove()
	//statusdb.Remove()
}

func BroadcastTx(tx core.Transactioner) {
	// broadcast to other peers
}

func GenerateBlock() *core.Block {
	var b core.Block
	b.Id = utils.CreateUuid()
	b.Num = statusdb.GetGpo().BlockNum + uint64(1)
	sec := time.Now().Unix()
	b.CreatedOn = uint64(sec)
	statusdb.MovePendingTxsToBlock(&b)

	applyBlock(&b)

	// update gpo
	gpo := statusdb.GetGpo()
	gpo.BlockId = b.Id
	gpo.BlockNum = b.Num
	gpo.Time = b.CreatedOn
	gpo.Supply += core.AmountPerBlock
	return &b
}

func applyBlock(b *core.Block) {
	for _, tx := range b.Transactions {
		eval.Apply(tx)
	}
}

func genesis() {
	// create dummp block and push it to ledger
	// TODO

	// update global status
	bid := utils.CreateUuid()
	createdOn := uint64(time.Now().Unix())
	gpo := statusdb.GetGpo()
	gpo.BlockId = bid
	gpo.BlockNum = 0
	gpo.Witness = core.InitWitness
	gpo.Time = createdOn
	gpo.Supply = core.InitAmount

	// update chain database
	var acc core.Account
	acc.Id = utils.CreateUuid() // should be public key string
	acc.Name = core.InitWitness
	acc.CreatedOn = createdOn
	acc.Coin = core.InitAmount
	statusdb.AddAccount(acc)
}
