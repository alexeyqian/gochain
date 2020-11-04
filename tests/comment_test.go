package tests

import (
	"testing"
	"time"

	"github.com/alexeyqian/gochain/core"
)

func TestCreateComment(t *testing.T) {
	var err error
	c := SetupTestChain()

	c.AddPendingTx(CreateTestAccount("alice"))
	c.AddPendingTx(CreateTestAccount("bob"))
	c.GenerateBlock()
	c.AddPendingTx(CreateTestArticle("test_article_001", "alice", "test_aritcle_title"))
	c.GenerateBlock()

	var tx core.CreateCommentTransaction
	tx.ParentId = "test_article_001"
	tx.CommentId = "test_comment_001"
	tx.Body = "comment_test_body"
	tx.Commentor = "bob"
	tx.CreatedOn = uint64(time.Now().Unix())
	c.AddPendingTx(&tx)
	c.GenerateBlock()

	comments := c.sdb.GetComments()
	if len(comments) != 1 {
		t.Errorf("create comment error")
	}

	comment, err := c.sdb.GetComment(tx.CommentId)
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

	TearDownTestChain(c)
}
