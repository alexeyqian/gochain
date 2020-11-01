package statusdb

import (
	"fmt"

	"github.com/alexeyqian/gochain/entity"
	"github.com/alexeyqian/gochain/store"
	"github.com/alexeyqian/gochain/undodb"
)

const GpoKey = "gpo_1"

const GpoBucket = "gpo"
const GpoStateBucket = "gpostate"
const AccountBucket = "account"
const AccountStateBucket = "accountstate"
const ArticleBucket = "article"
const ArticleStateBucket = "articlestate"
const CommentBucket = "comment"
const CommentStateBucket = "commentstate"
const VoteBucket = "vote"
const VoteStateBucket = "votestate"

type StatusDB struct {
	udb undodb.UndoableDB
}

func NewStatusDB(s store.Storage) *StatusDB {
	return &StatusDB{
		udb: undodb.NewUndoableDB(store.Storage),
	}
}

// Open has parameter MemoryStorage
func (sdb *StatusDB) Open() {
	sdb.udb.Open()
}

func (sdb *StatusDB) Close() {
	sdb.udb.Close()
}

func (sdb *StatusDB) Remove() {
	sdb.udb.Remove()
}

// =========== gpo ================

func (sdb *StatusDB) CreateGpo(e *entity.Gpo) error {
	return sdb.createEntity(GpoBucket, *e)
}

func (sdb *StatusDB) UpdateGpo(e *entity.Gpo) error {
	return sdb.updateEntity(GpoBucket, *e)
}

func (sdb *StatusDB) GetGpo() (*entity.Gpo, error) {
	var e entity.Gpo
	err := sdb.getEntityByID(GpoBucket, GpoKey, e)
	return &e, err
}

// =========== account ================

func (sdb *StatusDB) CreateAccount(e *entity.Account) error {
	return sdb.createEntity(AccountBucket, *e)
}

func (sdb *StatusDB) UpdateAccount(e *entity.Account) error {
	return sdb.updateEntity(AccountBucket, *e)
}

func (sdb *StatusDB) GetAccount(id string) (*entity.Account, error) {
	var e entity.Account
	err := sdb.getEntityByID(AccountBucket, id, e)
	return &e, err
}

func (sdb *StatusDB) GetAccounts() []*entity.Account {
	var res []*entity.Account
	for _, value := range sdb.getAll(AccountBucket) {
		var e entity.Account
		entity.Deserialize(e, value)
		res = append(res, &e)
	}
	return res
}

func GetAccountByName(name string) (*entity.Account, error) {
	var res *entity.Account
	for _, acc := range GetAccounts() {
		if acc.Name == name {
			res = acc
			break
		}
	}

	return res, nil
}

// =========== article ================

func (sdb *StatusDB) CreateArticle(e *entity.Article) error {
	return sdb.createEntity(ArticleBucket, *e)
}

func (sdb *StatusDB) UpdateArticle(e *entity.Article) error {
	return sdb.updateEntity(ArticleBucket, *e)
}

func (sdb *StatusDB) GetArticle(id string) (*entity.Article, error) {
	var e entity.Article
	err := sdb.getEntityByID(ArticleBucket, id, e)
	return &e, err
}

func (sdb *StatusDB) GetArticles() []*entity.Article {
	var res []*entity.Article
	for _, value := range sdb.getAll(ArticleBucket) {
		var e entity.Article
		entity.Deserialize(e, value)
		res = append(res, &e)
	}
	return res
}

// =========== comment ================

func (sdb *StatusDB) CreateComment(e *entity.Comment) error {
	return sdb.createEntity(CommentBucket, *e)
}

func (sdb *StatusDB) UpdateComment(e *entity.Comment) error {
	return sdb.updateEntity(CommentBucket, *e)
}

func (sdb *StatusDB) GetComment(id string) (*entity.Comment, error) {
	var e entity.Comment
	err := sdb.getEntityByID(CommentBucket, id, e)
	return &e, err
}

func (sdb *StatusDB) GetComments() []*entity.Comment {
	var res []*entity.Comment
	for _, value := range sdb.getAll(CommentBucket) {
		var e entity.Comment
		entity.Deserialize(e, value)
		res = append(res, &e)
	}
	return res
}

// =========== vote ================

func (sdb *StatusDB) CreateVote(e *entity.Vote) error {
	return sdb.createEntity(VoteBucket, *e)
}

func (sdb *StatusDB) UpdateVote(e *entity.Vote) error {
	return sdb.updateEntity(VoteBucket, *e)
}

func (sdb *StatusDB) GetVote(id string) (*entity.Vote, error) {
	var e entity.Vote
	err := sdb.getEntityByID(VoteBucket, id, e)
	return &e, err
}

func (sdb *StatusDB) GetVotes() []*entity.Vote {
	var res []*entity.Vote
	for _, value := range sdb.getAll(VoteBucket) {
		var e entity.Vote
		entity.Deserialize(e, value)
		res = append(res, &e)
	}
	return res
}

// ============ internel functions ==============
func (sdb *StatusDB) createEntity(bucket string, e entity.Entity) error {
	if !entity.HasID(e) {
		return fmt.Errorf("create: entity doesn't have ID")
	}

	return sdb.udb.Create(bucket, entity.GetID(e), entity.Serialize(e))
}

func (sdb *StatusDB) updateEntity(bucket string, e entity.Entity) error {
	if !entity.HasID(e) {
		return fmt.Errorf("update: entity doesn't have ID")
	}

	return sdb.udb.Update(bucket, entity.GetID(e), entity.Serialize(e))
}

func (sdb *StatusDB) getEntityByID(bucket, key string, e entity.Entity) error {
	data, err := sdb.udb.Get(bucket, key)
	if err != nil {
		return err
	}

	entity.Deserialize(e, data)
	return nil
}

func (sdb *StatusDB) getAllEntities(bucket string) []entity.Entity {
	var res []entity.Entity
	for k, v := range sdb.udb.GetAll(bucket) {
		var e entity.Entity

		switch bucket {
		case GpoBucket:
			e = entity.Gpo{}
		case AccountBucket:
			e = entity.Account{}
		case ArticleBucket:
			e = entity.Article{}
		case CommentBucket:
			e = entity.Comment{}
		case VoteBucket:
			e = entity.Vote{}
		defalut:
			panic("unknown entity type.")
		}

		temp, _ := entity.Deserialize(e, v)
		res = append(res, &temp)
	}

	return res
}

/*
func (sdb *StatusDB) getEntityType(bucket string) reflect.Type {
	switch bucket {
	case GpoBucket:
		var e entity.Gpo
		return reflect.TypeOf(e)
	case AccountBucket:
		var e entity.Account
		return reflect.TypeOf(e)
	case ArticleBucket:
		var e entity.Article
		return reflect.TypeOf(e)
	case CommentBucket:
		var e entity.Comment
		return reflect.TypeOf(e)
	case VoteBucket:
		var e entity.Vote
		return reflect.TypeOf(e)

	}

	panic("unknown type")
}
*/
