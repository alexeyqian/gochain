package core

import (
	"errors"

	"github.com/alexeyqian/gochain/entity"
	"github.com/alexeyqian/gochain/statusdb"
)

func (tx CreateAccountTransaction) Validate() error {
	if tx.AccountName == "" {
		return errors.New("name cannot be empty.")
	}
	return nil
}

func (tx TransferCoinTransaction) Validate() error {
	var err error
	var fromAcc *entity.Account

	fromAcc, err = statusdb.GetAccountByName(tx.From)
	if err != nil {
		return errors.New("transfer coin: from account is not exist")
	} else if fromAcc.Coin < tx.Amount {
		return errors.New("transfer coin: no enough coin")
	}

	_, err = statusdb.GetAccountByName(tx.To)
	if err != nil {
		return errors.New("transfer coin: to account is not exist")
	}

	return nil
}

func (tx CreateArticleTransaction) Validate() error {
	return nil
}

func (tx CreateCommentTransaction) Validate() error {
	return nil
}

func (tx VoteTransaction) Validate() error {
	return nil
}

func (tx CreateAccountTransaction) Apply() error {
	err := tx.Validate()
	if err != nil {
		return err
	}

	var acc entity.Account
	acc.ID = tx.AccountId
	acc.Name = tx.AccountName
	statusdb.AddAccount(&acc)

	return nil
}

func (tx TransferCoinTransaction) Apply() error {
	err := tx.Validate()
	if err != nil {
		return err
	}

	fromAcc, _ := statusdb.GetAccountByName(tx.From)
	toAcc, _ := statusdb.GetAccountByName(tx.To)
	fromAcc.Coin -= tx.Amount
	toAcc.Coin += tx.Amount

	// TODO: check errors
	statusdb.UpdateAccount(fromAcc)
	statusdb.UpdateAccount(toAcc)

	return nil
}

func (tx CreateArticleTransaction) Apply() error {
	err := tx.Validate()
	if err != nil {
		return err
	}

	var article entity.Article
	article.ID = tx.ArticleId
	article.Author = tx.Author
	article.Title = tx.Title
	article.Body = tx.Body
	article.Meta = tx.Meta
	statusdb.AddArticle(&article)

	acc, _ := statusdb.GetAccountByName(tx.Author)
	acc.ArticleCount += 1
	// TODO: check error
	statusdb.UpdateAccount(acc)
	return nil
}

func (tx CreateCommentTransaction) Apply() error {
	err := tx.Validate()
	if err != nil {
		return err
	}

	var comment entity.Comment
	comment.ID = tx.CommentId
	comment.ParentId = tx.ParentId
	comment.Commentor = tx.Commentor
	comment.Body = tx.Body
	comment.CreatedOn = tx.CreatedOn
	statusdb.AddComment(&comment)

	return nil
}

func (tx VoteTransaction) Apply() error {
	err := tx.Validate()
	if err != nil {
		return err
	}

	var vote entity.Vote
	vote.ID = tx.ID
	vote.ParentId = tx.ParentId
	vote.ParentType = tx.ParentType
	vote.Direction = tx.Direction
	vote.VotePower = tx.VotePower
	statusdb.AddVote(&vote)

	if vote.ParentType == VoteParentTypeAccount {
		account, _ := statusdb.GetAccount(vote.ParentId)

		if vote.Direction > 0 {
			account.UpVotes += 1
			account.VotePower += vote.VotePower
		} else {
			account.DownVotes += 1
			account.VotePower -= vote.VotePower
		}

		statusdb.UpdateAccount(account)

	} else if vote.ParentType == VoteParentTypeArticle {
		article, _ := statusdb.GetArticle(vote.ParentId)
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
		statusdb.UpdateArticle(article)

	} else if vote.ParentType == VoteParentTypeComment {
		comment, _ := statusdb.GetComment(vote.ParentId)

		if vote.Direction > 0 {
			comment.UpVotes += 1
			comment.VotePower += vote.VotePower
		} else {
			comment.DownVotes += 1
			comment.VotePower -= vote.VotePower
		}

		statusdb.UpdateComment(comment)
	}

	return nil
}
