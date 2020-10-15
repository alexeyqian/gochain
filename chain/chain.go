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
	statusdb.Open()
	if !_isGenesised {
		genesis()
	}
}

func Close() {
	ledger.Close()
	statusdb.Close()
}

func Remove() {
	ledger.Remove()
	statusdb.Remove()
}

func GetBlock(num int) (*core.Block, error) {
	// TODO: use cache to speed up reading
	bs, err := ledger.Read(num)
	if err != nil {
		return nil, err
	}
	return core.UnSerializeBlock(bs)
}

func BroadcastTx(tx core.Transactioner) {
	// broadcast to other peers
}

func AddPendingTx(tx core.Transactioner) error {
	err := eval.Validate(tx)
	if err == nil {
		statusdb.AddPendingTransaction(tx)
	}
	return err
}

func GenerateBlock() *core.Block {
	var b core.Block
	b.Id = utils.CreateUuid()
	b.PrevBlockId = statusdb.GetGpo().BlockId
	b.Num = statusdb.GetGpo().BlockNum + uint64(1)
	sec := time.Now().Unix()
	b.CreatedOn = uint64(sec)
	statusdb.MovePendingTransactionsToBlock(&b)

	for _, tx := range b.Transactions {
		terr := eval.Apply(tx)
		if terr != nil {
			// move tx to invalid tx
			//
		}
	}

	// update gpo
	gpo := statusdb.GetGpo()
	gpo.BlockId = b.Id
	gpo.BlockNum = b.Num
	gpo.Time = b.CreatedOn
	gpo.Supply += core.AmountPerBlock

	// append new block to ledger
	sb, _ := core.SerializeBlock(&b)
	ledger.Append(sb)

	return &b
}

func genesis() {
	// update global status
	gpo := statusdb.GetGpo()
	gpo.BlockId = core.BlockZeroId
	gpo.BlockNum = 0
	gpo.Witness = core.InitWitness
	gpo.Time = core.GenesisTime
	gpo.Supply = core.InitAmount

	// update chain database
	var acc core.Account
	acc.Id = utils.CreateUuid() // should be public key string
	acc.Name = core.InitWitness
	acc.CreatedOn = core.GenesisTime
	acc.Coin = core.InitAmount
	statusdb.AddAccount(acc)

	// update ledger, create a dummy block 0
	b := core.Block{Id: core.BlockZeroId, Num: 0, CreatedOn: core.GenesisTime, Witness: core.InitWitness}
	sb, _ := core.SerializeBlock(&b)
	ledger.Append(sb)
}
