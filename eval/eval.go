package eval

import (
	"errors"

	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/statusdb"
)

func Validate(tx core.Transactioner) error {
	txType := tx.TxType()
	var err error
	switch txType {
	case core.CreateAccountTransactionType:
		t := tx.(core.CreateAccountTransaction)
		if t.AccountName == "" {
			err = errors.New("name cannot be empty.")
		}

	case core.TransferCoinTransactionType:
		t := tx.(core.TransferCoinTransaction)
		//fmt.Printf("from: %s, to: %s\n", tct.From, tct.To)
		fromAcc := statusdb.GetAccount(t.From)
		toAcc := statusdb.GetAccount(t.To)

		if fromAcc == nil {
			err = errors.New("transfer coin: from account is not exist")
		}
		if toAcc == nil {
			err = errors.New("transfer coin: to account is not exist")
		}
		if fromAcc.Coin < t.Amount {
			err = errors.New("transfer coin: no enough coin")
		}
	}

	return err
}

func Apply(tx core.Transactioner) error {
	err := Validate(tx)
	if err != nil {
		return err
	}

	txType := tx.TxType()
	switch txType {
	case core.CreateAccountTransactionType:
		cat := tx.(core.CreateAccountTransaction)
		var acc core.Account
		acc.Id = cat.AccountId
		acc.Name = cat.AccountName

		statusdb.AddAccount(acc)
	case core.TransferCoinTransactionType:
		tct := tx.(core.TransferCoinTransaction)
		//fmt.Printf("from: %s, to: %s\n", tct.From, tct.To)
		fromAcc := statusdb.GetAccount(tct.From)
		toAcc := statusdb.GetAccount(tct.To)

		if fromAcc == nil {
			return errors.New("transfer coin: from account is not exist")
		}
		if toAcc == nil {
			return errors.New("transfer coin: to account is not exist")
		}
		if fromAcc.Coin < tct.Amount {
			return errors.New("transfer coin: no enough coin")
		}
		fromAcc.Coin -= tct.Amount
		toAcc.Coin += tct.Amount

	case core.CreateArticleTransactionType:
		t := tx.(core.CreateArticleTransaction)
		var article core.Article
		article.ArticleId = t.ArticleId
		article.Author = t.Author
		article.Title = t.Title
		article.Body = t.Body
		article.Meta = t.Meta
		statusdb.AddArticle(article)

		acc := statusdb.GetAccount(t.Author)
		acc.ArticleCount += 1
	}

	return nil
}
