package core

import (
	"github.com/alexeyqian/gochain/entity"
	"github.com/alexeyqian/gochain/statusdb"
)

func (tx *CreateArticleTransaction) Validate(sdb *statusdb.StatusDB) error {
	return nil
}

func (tx *CreateArticleTransaction) Apply(sdb *statusdb.StatusDB) error {
	err := tx.Validate(sdb)
	if err != nil {
		return err
	}

	var article entity.Article
	article.ID = tx.ArticleId
	article.Author = tx.Author
	article.Title = tx.Title
	article.Body = tx.Body
	article.Meta = tx.Meta
	sdb.AddArticle(&article)

	acc, _ := c.sdb.GetAccountByName(tx.Author)
	acc.ArticleCount += 1
	// TODO: check error
	sdb.UpdateAccount(acc)
	return nil
}
