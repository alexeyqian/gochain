package chain

import (
	"crypto/ecdsa"
	"path/filepath"
	"testing"

	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/ledger"
	"github.com/alexeyqian/gochain/store"
	"github.com/alexeyqian/gochain/utils"
)

var alicekey *ecdsa.PrivateKey

func setup() {
	alicekey, _ = utils.GenerateKey()
}

func newChain(dir string) *Chain {
	lgr := ledger.NewFileLedger(dir)
	storage := store.NewBoltStorage(filepath.Join(dir, "status.db"))
	c := NewChain(lgr, storage)

	c.Open()
	return c
}

func createOneBlock(c *Chain) *core.Block {
	slottime := getNextBlockTime(c.Gpo())
	name := getNextWitness(c.Gpo(), c.Wso())
	witness, _ := c.sdb.GetAccountByName(name)

	b, _ := c.GenerateBlock(slottime, witness, alicekey)
	return b
}

func TestBranch1(t *testing.T) {
	setup()
	chain1 := newChain("test_chain1_data")
	chain2 := newChain("test_chain2_data")

	count := 10
	// create 10 empty blocks in both chain1 and chain2
	for i := 0; i < count; i++ {
		b := createOneBlock(chain1)
		chain1.PushBlock(b)
		chain2.PushBlock(b)
	}
}
