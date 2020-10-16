package eval

import (
	"errors"

	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/entity"
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
		fromAcc := statusdb.GetAccountByName(t.From)
		toAcc := statusdb.GetAccountByName(t.To)

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
		var acc entity.Account
		acc.Id = cat.AccountId
		acc.Name = cat.AccountName

		statusdb.AddAccount(acc)
	case core.TransferCoinTransactionType:
		tct := tx.(core.TransferCoinTransaction)
		//fmt.Printf("from: %s, to: %s\n", tct.From, tct.To)
		fromAcc := statusdb.GetAccountByName(tct.From)
		toAcc := statusdb.GetAccountByName(tct.To)

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
		var article entity.Article
		article.ArticleId = t.ArticleId
		article.Author = t.Author
		article.Title = t.Title
		article.Body = t.Body
		article.Meta = t.Meta
		statusdb.AddArticle(article)

		acc := statusdb.GetAccountByName(t.Author)
		acc.ArticleCount += 1

	case core.CreateCommentTransactionType:
		t := tx.(core.CreateCommentTransaction)
		var comment entity.Comment
		comment.ParentId = t.ParentId
		comment.CommentId = t.CommentId
		comment.Commentor = t.Commentor
		comment.Body = t.Body
		comment.CreatedOn = t.CreatedOn
		statusdb.AddComment(comment)

	case core.VoteTransactionType:
		t := tx.(core.VoteTransaction)

		var vote entity.Vote
		vote.Id = t.Id
		vote.ParentId = t.ParentId
		vote.ParentType = t.ParentType
		vote.Direction = t.Direction
		vote.VotePower = t.VotePower
		statusdb.AddVote(vote)

		if vote.ParentType == core.VoteParentTypeAccount {
			account := statusdb.GetAccount(vote.ParentId)
			if account == nil {
				return errors.New("vote account not exist")
			}
			if vote.Direction > 0 {
				account.UpVotes += 1
				account.VotePower += vote.VotePower
			} else {
				account.DownVotes += 1
				account.VotePower -= vote.VotePower
			}

		} else if vote.ParentType == core.VoteParentTypeArticle {
			article := statusdb.GetArticle(vote.ParentId)
			if article == nil {
				return errors.New("vote articel not exist")
			}
			if vote.Direction > 0 {
				article.UpVotes += 1
				article.VotePower += vote.VotePower
			} else {
				article.DownVotes += 1
				article.VotePower -= vote.VotePower
			}

		} else if vote.ParentType == core.VoteParentTypeComment {
			comment := statusdb.GetComment(vote.ParentId)
			if comment == nil {
				return errors.New("vote comment not exist")
			}
			if vote.Direction > 0 {
				comment.UpVotes += 1
				comment.VotePower += vote.VotePower
			} else {
				comment.DownVotes += 1
				comment.VotePower -= vote.VotePower
			}
		}

	}

	return nil
}
