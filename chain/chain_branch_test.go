package chain

import (
	"path/filepath"
	"testing"

	"github.com/alexeyqian/gochain/chain"
	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/ledger"
	"github.com/alexeyqian/gochain/store"
	"github.com/alexeyqian/gochain/utils"
)

func TNewChain(dir string) *Chain {
	lgr := ledger.NewFileLedger(dir)
	storage := store.NewBoltStorage(filepath.Join(dir, "status.db"))
	c := NewChain(lgr, storage)

	c.Open()
	return c
}

func TCreateOneBlock(c *Chain, privkey string) *core.Block {
	slottime := chain.GetNextBlockTime(c.Gpo())
	name := chain.GetNextWitness(c.Wso())
	witness, _ := c.sdb.GetAccountByName(name)

	return c.GenerateBlock(slottime, witness, privkey)
}

func TestBranch1(t *testing.T) {

	chain1 := TNewChain("test_chain1_data")
	chain2 := TNewChain("test_chain2_data")

	alicekey, _ := utils.GenerateKey()
	//bobkey, _ := utils.GenerateKey()
	slotnum := 1
	count := 10
	// create 10 empty blocks in both chain1 and chain2
	for i := 0; i < count; i++ {
		b := TCreateOneBlock(chain1, alicekey)
		chain1.PushBlock(b)
		chain2.PushBlock(b)
	}
}
