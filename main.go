package main

import (
	"fmt"
	"path/filepath"

	"github.com/alexeyqian/gochain/chain"
	"github.com/alexeyqian/gochain/ledger"
	"github.com/alexeyqian/gochain/store"
)

func main() {
	fmt.Printf("starting ...\n")
	testDataDir := "test_data"
	lgr := ledger.NewFileLedger(testDataDir)
	storage := store.NewBoltStorage(filepath.Join(testDataDir, "status.db"))
	c := chain.NewChain(lgr, storage)
	fmt.Println("arrive here 2")
	c.Open()

}
