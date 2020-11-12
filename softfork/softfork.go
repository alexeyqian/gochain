package softfork

import (
	"fmt"

	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/entity"
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

	var item entity.SoftForkItem
	item.ID = b.ID
	item.PrevBlockID = b.PrevBlockId
	item.BlockNum = b.Num
	data, _ := core.SerializeBlock(b)
	item.BlockData = data
	fork.sdb.CreateSoftForkItem(&item)
	fork.Head = item.ID
	return &fork
}

// fork should always have at lease one item,
// every time the database opens, it will starts with the last irriversable block.
func (fork *SoftFork) PopBlock() {
	if fork.Head == "" {
		panic("fork head is empty.")
	}

	item, err := fork.sdb.GetSoftForkItem(fork.Head)
	if err != nil {
		panic(fmt.Sprintf("cannot find the soft fork item: %s", fork.Head))
	}
	// TODO: remove previous head block
	// check if it's still the longest branch?

	fork.Head = item.PrevBlockID
}
