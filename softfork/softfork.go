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

func (fork *SoftFork) reset() {
	fork.sdb.TruncateSoftFork()
}

// multiple branches are stored in same table
// push a block into linked table if it can be linked to one of the branches in the table
// and also move head to biggest block num if larger than current num
func (fork *SoftFork) PushBlock(b core.Block) error {
	// validate block before push, expired, max_depth of fork
	// make sure head is not empty
	// check duplication of BlockID before insert
	// check if block can be linked to current branches, if cannot, then put it into unlinked pool

	item, err := fork.sdb.GetSoftForkItem(b.ID)
	if item != nil || err == nil {
		return fmt.Errorf("soft fork: cannot insert duplicated item")
	}

	item.ID = b.ID
	item.BlockNum = b.Num
	item.BlockDatam, _ = core.SerializeBlock(&b)
	item.PrevBlockID = b.PrevBlockId

	fork.sdb.CreateSoftForkItem(&item)
	// switch head if num is bigger.
	if( item.BlockNum > _head->num ){
		_head = item
	} 

	return nil
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

// might have multiple blocks with same block num
func (fork *SoftFork) FetchBlocksByNumber(num uint64) []*core.Block {
	var blocks []*core.Block

	allItems := fork.sdb.GetSoftForkItems()
	for _, item := range allItems {
		if item.BlockNum == num {
			tempB, _ := core.UnSerializeBlock(item.BlockData)
			blocks = append(blocks, tempB)
		}
	}
	return blocks
}
