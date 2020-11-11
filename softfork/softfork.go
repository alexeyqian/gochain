package softfork

import (
	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/statusdb"
)

type SoftFork struct {
	Head string
	sdb  *statusdb.StatusDB
}

func NewSoftFork(b core.Block, db *statusdb.StatusDB) *SoftFork {
	var fork SoftFork
	fork.sdb = db
	fork.reset()

	sb := fork.insert(b)
	fork.Head = &sb
	return &fork
}

func (fork *SoftFork) insert(b core.Block) {

}

// fork should always have at lease one item,
// every time the database opens, it will starts with the last irriversable block.
func (fork *SoftFork) PopBlock() {
	fork.Head = fork.Head.prev
	// TODO: handle empty scenario
}
