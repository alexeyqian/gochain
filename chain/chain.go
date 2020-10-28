package chain

import (
	"time"

	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/entity"
	"github.com/alexeyqian/gochain/ledger"
	"github.com/alexeyqian/gochain/statusdb"
	"github.com/alexeyqian/gochain/utils"
)

var _dataFolder string = "data"
var _isGenesised = false

var _pendingTransactions []core.Transactioner

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
	_pendingTransactions = nil
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

func ReceiveTx(tx core.Transactioner) error {
	// check if already has the tx
	// validate tx: two validations, fast validate and full validate

	return nil
}

func GetPendingTx() []core.Transactioner {
	return _pendingTransactions
}

func movePendingTransactionsToBlock(b *core.Block) {
	i := 0
	for _, tx := range _pendingTransactions {
		if i >= core.MaxTransactionsInBlock {
			break
		}
		b.AddTransaction(tx)
		i++
	}

	if len(_pendingTransactions) > core.MaxTransactionsInBlock {
		_pendingTransactions = _pendingTransactions[core.MaxTransactionsInBlock:]
	}
}

func AddPendingTx(tx core.Transactioner) error {
	err := tx.Validate()
	if err == nil {
		_pendingTransactions = append(_pendingTransactions, tx)
	}
	return err
}

func GenerateBlock() *core.Block {
	var b core.Block
	var gpo *entity.Gpo

	gpo, _ = statusdb.GetGpo()

	b.ID = utils.CreateUuid()
	b.PrevBlockId = gpo.BlockId
	b.Num = gpo.BlockNum + uint64(1)
	sec := time.Now().Unix()
	b.CreatedOn = uint64(sec)
	movePendingTransactionsToBlock(&b)

	// TODO: should gpo be updated during tx.Apply ??
	for _, tx := range b.Transactions {
		terr := tx.Apply() // gpo might be updated during tx.Apply()
		if terr != nil {
			// move tx to invalid tx
			//
		}
	}

	gpo, _ = statusdb.GetGpo()
	gpo.BlockId = b.ID
	gpo.BlockNum = b.Num
	gpo.Time = b.CreatedOn
	gpo.Supply += core.AmountPerBlock
	statusdb.UpdateGpo(gpo)

	// append new block to ledger
	sb, _ := core.SerializeBlock(&b)
	ledger.Append(sb)

	return &b
}

func genesis() {
	// update global status
	var gpo entity.Gpo
	gpo.BlockId = core.BlockZeroId
	gpo.BlockNum = 0
	gpo.Witness = core.InitWitness
	gpo.Time = core.GenesisTime
	gpo.Supply = core.InitAmount
	statusdb.AddGpo(&gpo)

	// update chain database
	var acc entity.Account
	acc.ID = utils.CreateUuid() // should be public key string
	acc.Name = core.InitWitness
	acc.CreatedOn = core.GenesisTime
	acc.Coin = core.InitAmount
	statusdb.AddAccount(&acc)

	// update ledger, create a dummy block 0
	b := core.Block{ID: core.BlockZeroId, Num: 0, CreatedOn: core.GenesisTime, Witness: core.InitWitness}
	sb, _ := core.SerializeBlock(&b)
	ledger.Append(sb)
}
