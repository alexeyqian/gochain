package tests

import (
	"testing"
	"time"

	"github.com/alexeyqian/gochain/chain"
	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/statusdb"
	"github.com/alexeyqian/gochain/utils"
)

func TestCreateArticle(t *testing.T) {
	chain.Open(TestDataDir)

	chain.AddPendingTx(CreateTestAccount("alice"))
	chain.GenerateBlock()

	var tx core.CreateArticleTransaction
	tx.ArticleId = utils.CreateUuid()
	tx.Author = "alice"
	tx.Title = "test_title"
	tx.Body = "test_body"
	tx.Meta = `{"tags": "test,science"}`
	chain.AddPendingTx(tx)
	chain.GenerateBlock()

	acc, _ := statusdb.GetAccountByName("alice")
	if acc.ArticleCount != 1 {
		t.Errorf("create article >>> article count expected: %d, actual: %d", 1, acc.ArticleCount)
	}

	articles := statusdb.GetArticles()
	if articles == nil || len(articles) != 1 {
		t.Errorf("create article error")
	}

	article, _ := statusdb.GetArticle(tx.ArticleId)
	//fmt.Printf("article: %v", article)
	if article == nil || article.Title != "test_title" {
		t.Errorf("create and get article error")
	}

	chain.Close()
	chain.Remove()
}

func TestCreateComment(t *testing.T) {
	var err error

	chain.Open(TestDataDir)

	chain.AddPendingTx(CreateTestAccount("alice"))
	chain.AddPendingTx(CreateTestAccount("bob"))
	chain.GenerateBlock()
	chain.AddPendingTx(CreateTestArticle("test_article_001", "alice", "test_aritcle_title"))
	chain.GenerateBlock()

	var tx core.CreateCommentTransaction
	tx.ParentId = "test_article_001"
	tx.CommentId = "test_comment_001"
	tx.Body = "comment_test_body"
	tx.Commentor = "bob"
	tx.CreatedOn = uint64(time.Now().Unix())
	chain.AddPendingTx(tx)
	chain.GenerateBlock()

	comments := statusdb.GetComments()
	if len(comments) != 1 {
		t.Errorf("create comment error")
	}

	comment, err := statusdb.GetComment(tx.CommentId)
	if err != nil {
		t.Errorf("comment is nil")
	} else {
		if comment.Body != "comment_test_body" {
			t.Errorf("create and get comment error")
		}
		if comment.ParentId != "test_article_001" {
			t.Errorf("create comment parent error")
		}
	}

	chain.Close()
	chain.Remove()
}

// TestCreateNestedComment -> comment level <= 5
