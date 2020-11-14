package chain

import (
	"path/filepath"
	"testing"

	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/ledger"
	"github.com/alexeyqian/gochain/store"
)

func TNewChain(dir string) *Chain {
	lgr := ledger.NewFileLedger(dir)
	storage := store.NewBoltStorage(filepath.Join(dir, "status.db"))
	c := NewChain(lgr, storage)

	c.Open()
	return c
}

func TCreateOneBlock(c *Chain, privKey string, slotNum int) *core.Block {
	slotTime := getSlotTime(c.Gpo(), slotNum)
	witness := getScheduledWitness(c.Gpo(), c.Wso(), slotNum)

	b := c.GenerateBlock(slotTime, witness, privKey)
}

func TestBranch1(t *testing.T) {

	chain1 := TNewChain("test_chain1_data")
	chain2 := TNewChain("test_chain2_data")

	// create 10 empty blocks in both chain1 and chain2
	for i := 0; i < 10; i++ {
		b := createOneBlock()
		chain1.PushBlock(b)
		chain2.PushBlock(b)
	}
}
