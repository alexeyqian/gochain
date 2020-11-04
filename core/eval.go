package chain

import (
	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/statusdb"
)

type Evaluator struct {
	sdb *statusdb.StatusDB
}

func NewEvaluator(db *statusdb.StatusDB) *Evaluator {
	return &Evaluator{
		sdb: db,
	}
}

func (e *Evaluator) ApplyTx(tx core.Transactioner) error {

}
