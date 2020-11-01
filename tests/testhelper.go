package tests

import (
	"fmt"
	"time"

	"github.com/alexeyqian/gochain/chain"
	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/statusdb"
	"github.com/alexeyqian/gochain/utils"
)

const TestDataDir = "data"

func CreateTestAccount(name string) core.Transactioner {
	var tx core.CreateAccountTransaction
	tx.AccountId = utils.CreateUuid()
	tx.AccountName = name
	tx.CreatedOn = uint64(time.Now().Unix())
	return tx
}

func CreateTestArticle(id string, author string, title string) core.Transactioner {
	var tx core.CreateArticleTransaction
	tx.ArticleId = id
	tx.Author = author
	tx.Title = title
	tx.Body = "test_body"
	tx.Meta = `{"tags": "test,science"}`
	return tx
}

func CreateTestBlocks(count int, datadir string) {
	chain.Open(datadir)

	i := 1
	for i <= 20 {
		tx := CreateTestAccount(fmt.Sprintf("test_account_name_%d", i))
		chain.AddPendingTx(tx)
		//chain.BroadcastTx(tx)
		b := chain.GenerateBlock()
		if b.Num != uint64(i) {
			panic(fmt.Sprintf("expected: %d, actual: %d", i, b.Num))
		}

		gpo, _ := statusdb.GetGpo()
		if gpo.BlockNum != b.Num {
			panic(fmt.Sprintf("gpo num expected: %d, actual: %d", 20, gpo.BlockNum))
		}
		if gpo.BlockId != b.ID {
			panic(fmt.Sprintf("gpo id expected: %s, actual: %s", b.ID, gpo.BlockId))
		}
		if gpo.Time != b.CreatedOn {
			panic(fmt.Sprintf("gpo time expected: %d, actual: %d", b.CreatedOn, gpo.Time))
		}
		if gpo.Supply != (core.InitAmount + core.AmountPerBlock*uint64(i)) {
			panic(fmt.Sprintf("generate block gpo amount expected: %d, actual: %d", core.InitAmount+core.AmountPerBlock*i, gpo.Supply))
		}

		// TODO: validate block and previous block hash/linking
		prevb, _ := chain.GetBlock(i - 1)
		//fmt.Printf("prevb id: %s", prevb.ID)
		if b.PrevBlockId != prevb.ID {
			panic(fmt.Sprintf(("block linking is broken"))
		}

		i++
	}

	//chain.Close()
	//chain.Remove()
}
