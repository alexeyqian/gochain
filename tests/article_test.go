package tests

import (
	"testing"

	"github.com/alexeyqian/gochain/utils"

	"github.com/alexeyqian/gochain/chain"
	core "github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/statusdb"
)

func TestCreateArticle(t *testing.T) {
	chain.Open(TestDataDir)

	chain.AddPendingTx(TstCreateAccount("alice"))
	chain.GenerateBlock()

	var tx core.CreateArticleTransaction
	tx.ArticleId = utils.CreateUuid()
	tx.Author = "alice"
	tx.Title = "test_title"
	tx.Body = "test_body"
	tx.Meta = `{"tags": "test,science"}`
	chain.AddPendingTx(tx)
	chain.GenerateBlock()

	acc := statusdb.GetAccount("alice")
	if acc.ArticleCount != 1 {
		t.Errorf("create article >>> article count expected: %d, actual: %d", 1, acc.ArticleCount)
	}

	articles := statusdb.GetArticles()
	if articles == nil || len(articles) != 1 {
		t.Errorf("create article error")
	}

	article := statusdb.GetArticle(tx.ArticleId)
	//fmt.Printf("article: %v", article)
	if article == nil || article.Title != "test_title" {
		t.Errorf("create and get article error")
	}

	chain.Close()
	chain.Remove()
}

func TestCreateComment(t *testing.T) {

}

// TestCreateNestedComment -> comment level <= 5
