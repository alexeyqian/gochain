package tests

import (
	"testing"

	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/utils"
)

func TestCreateArticle(t *testing.T) {
	c := SetupTestChain()

	c.AddPendingTx(CreateTestAccount("alice"))
	c.GenerateBlock()

	var tx core.CreateArticleTransaction
	tx.ArticleId = utils.CreateUuid()
	tx.Author = "alice"
	tx.Title = "test_title"
	tx.Body = "test_body"
	tx.Meta = `{"tags": "test,science"}`
	c.AddPendingTx(&tx)
	c.GenerateBlock()

	acc, _ := c.sdb.GetAccountByName("alice")
	if acc.ArticleCount != 1 {
		t.Errorf("create article >>> article count expected: %d, actual: %d", 1, acc.ArticleCount)
	}

	articles := c.sdb.GetArticles()
	if articles == nil || len(articles) != 1 {
		t.Errorf("create article error")
	}

	article, _ := c.sdb.GetArticle(tx.ArticleId)
	//fmt.Printf("article: %v", article)
	if article == nil || article.Title != "test_title" {
		t.Errorf("create and get article error")
	}

	TearDownTestChain(c)
}

// TestCreateNestedComment -> comment level <= 5
