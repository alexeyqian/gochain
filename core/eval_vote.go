package core

import (
	"errors"

	"github.com/alexeyqian/gochain/entity"
	"github.com/alexeyqian/gochain/statusdb"
)

func (tx *VoteTransaction) Validate(sdb *statusdb.StatusDB) error {
	return nil
}

func (tx *VoteTransaction) Apply(sdb *statusdb.StatusDB) error {
	err := tx.Validate(sdb)
	if err != nil {
		return err
	}

	var vote entity.Vote
	vote.ID = tx.ID
	vote.ParentId = tx.ParentId
	vote.ParentType = tx.ParentType
	vote.Direction = tx.Direction
	vote.VotePower = tx.VotePower
	sdb.CreateVote(&vote)

	if vote.ParentType == VoteParentTypeAccount {
		account, _ := sdb.GetAccount(vote.ParentId)

		if vote.Direction > 0 {
			account.UpVotes += 1
			account.VotePower += vote.VotePower
		} else {
			account.DownVotes += 1
			account.VotePower -= vote.VotePower
		}

		sdb.UpdateAccount(account)

	} else if vote.ParentType == VoteParentTypeArticle {
		article, _ := sdb.GetArticle(vote.ParentId)
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
		sdb.UpdateArticle(article)

	} else if vote.ParentType == VoteParentTypeComment {
		comment, _ := sdb.GetComment(vote.ParentId)

		if vote.Direction > 0 {
			comment.UpVotes += 1
			comment.VotePower += vote.VotePower
		} else {
			comment.DownVotes += 1
			comment.VotePower -= vote.VotePower
		}

		sdb.UpdateComment(comment)
	}

	return nil
}
