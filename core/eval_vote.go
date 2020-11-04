package core

import (
	"errors"

	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/entity"
)

func (c *Chain) ValidateVote(tx core.VoteTransaction) error {
	return nil
}

func (c *Chain) ApplyVote(tx *core.VoteTransaction) error {
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
	c.sdb.AddVote(&vote)

	if vote.ParentType == VoteParentTypeAccount {
		account, _ := c.sdb.GetAccount(vote.ParentId)

		if vote.Direction > 0 {
			account.UpVotes += 1
			account.VotePower += vote.VotePower
		} else {
			account.DownVotes += 1
			account.VotePower -= vote.VotePower
		}

		c.sdb.UpdateAccount(account)

	} else if vote.ParentType == VoteParentTypeArticle {
		article, _ := c.sdb.GetArticle(vote.ParentId)
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
		c.sdb.UpdateArticle(article)

	} else if vote.ParentType == VoteParentTypeComment {
		comment, _ := c.sdb.GetComment(vote.ParentId)

		if vote.Direction > 0 {
			comment.UpVotes += 1
			comment.VotePower += vote.VotePower
		} else {
			comment.DownVotes += 1
			comment.VotePower -= vote.VotePower
		}

		c.sdb.UpdateComment(comment)
	}

	return nil
}
