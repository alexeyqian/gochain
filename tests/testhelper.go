package tests

import (
	"time"

	"github.com/alexeyqian/gochain/core"
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
