package statusdb

import (
	"bytes"
	"encoding/gob"

	"github.com/alexeyqian/gochain/entity"
)

type BoltDataProvider struct {
}

func serializeEntity(entityType string, e entity.Entity) ([]byte, error) {
	var err error
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	switch entityType {
	case GpoTable:
		ce := e.(entity.Gpo)
		err = encoder.Encode(ce)
		return result.Bytes(), err
	case AccountTable:
		ce := e.(entity.Account)
		err = encoder.Encode(ce)
		return result.Bytes(), err
	case ArticleTable:
		ce := e.(entity.Article)
		err = encoder.Encode(ce)
		return result.Bytes(), err
	case CommentTable:
		ce := e.(entity.Comment)
		err = encoder.Encode(ce)
		return result.Bytes(), err
	case VoteTable:
		ce := e.(entity.Vote)
		err = encoder.Encode(ce)
		return result.Bytes(), err
	}

	panic("serialize: unknown entity type")
}

func deserializeEntity(entityType string, data []byte) (entity.Entity, error) {
	var err error
	decoder := gob.NewDecoder(bytes.NewReader(data))

	switch entityType {
	case GpoTable:
		var e entity.Gpo
		err = decoder.Decode(&e)
		return &e, err
	case AccountTable:
		var e entity.Account
		err = decoder.Decode(&e)
		return &e, err
	case ArticleTable:
		var e entity.Article
		err = decoder.Decode(&e)
		return &e, err
	case CommentTable:
		var e entity.Comment
		err = decoder.Decode(&e)
		return &e, err
	case VoteTable:
		var e entity.Vote
		err = decoder.Decode(&e)
		return &e, err
	}

	panic("deserialize: unknown entity type")
}
