package core

import (
	"github.com/alexeyqian/gochain/entity"
	"github.com/alexeyqian/gochain/statusdb"
)

func (tx *CreateCommentTransaction) Validate(sdb *statusdb.StatusDB) error {
	return nil
}

func (tx *CreateCommentTransaction) Apply(sdb *statusdb.StatusDB) error {
	err := tx.Validate(sdb)
	if err != nil {
		return err
	}

	var comment entity.Comment
	comment.ID = tx.CommentId
	comment.ParentId = tx.ParentId
	comment.Commentor = tx.Commentor
	comment.Body = tx.Body
	comment.CreatedOn = tx.CreatedOn
	sdb.CreateComment(&comment)

	return nil
}
