package statusdb

import (
	"fmt"

	"github.com/alexeyqian/gochain/entity"
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
	store Storage
}

func NewStatusDB(s Storage) *StatusDB {
	return &StatusDB{
		store: Storage,
	}
}

// Open has parameter MemoryStorage
func (sdb *StatusDB) Open() {
	sdb.store.Open()
}

func (sdb *StatusDB) Close() {
	sdb.store.Close()
}

func (sdb *StatusDB) Remove() {
	sdb.store.RemoveAll()
}

// =========== gpo ================

func (sdb *StatusDB) GetGpo() (*entity.Gpo, error) {
	var e Gpo
	err := sdb.getEntityByID(GpoBucket, GpoKey, e)
	return &e, err
}

func (sdb *StatusDB) CreateGpo(e *entity.Gpo) error {
	return sdb.createEntity(GpoBucket, *e)
}

func (sdb *StatusDB) UpdateGpo(e *entity.Gpo) error {
	return sdb.updateEntity(GpoBucket, *e)
}

// =========== account ================

func (sdb *StatusDB) CreateAccount(e *entity.Account) error {
	return sdb.createEntity(AccountBucket, *e)
}

func UpdateAccount(e *entity.Account) error {
	return sdb.updateEntity(AccountBucket, *e)
}

func GetAccount(id string) (*entity.Account, error) {
	var e entity.Account
	err := sdb.getEntityByID(AccountBucket, id, e)
	return &e, err
}

func (sdb *StatusDB) GetAccounts() []*entity.Account {
	var res []*entity.Account
	for _, value := range sdb.store.GetAll(AccountBucket) {
		var e entity.Account
		entity.DeserializeEntity(e, value)
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

func AddArticle(e *entity.Article) error {
	_dp.Put(addPrefix(ArticleBucket, e.ID), *e)
	return nil
}

func GetArticles() []*entity.Article {
	var res []*entity.Article
	for _, value := range _dp.GetAll(ArticleBucket) {
		temp, _ := value.(entity.Article)
		res = append(res, &temp)
	}
	return res
}

func GetArticle(id string) (*entity.Article, error) {
	e, _ := _dp.Get(addPrefix(ArticleBucket, id))
	ce, _ := e.(entity.Article)
	return &ce, nil
}

func UpdateArticle(e *entity.Article) error {
	_dp.Put(addPrefix(ArticleBucket, e.ID), *e)
	return nil // TODO: check update errors
}

// =========== comment ================

func AddComment(e *entity.Comment) error {
	_dp.Put(addPrefix(CommentBucket, e.ID), *e)
	return nil
}

func GetComments() []*entity.Comment {
	var res []*entity.Comment
	for _, value := range _dp.GetAll(CommentBucket) {
		temp, _ := value.(entity.Comment)
		res = append(res, &temp)
	}
	return res
}

func GetComment(id string) (*entity.Comment, error) {
	e, _ := _dp.Get(addPrefix(CommentBucket, id))
	ce, _ := e.(entity.Comment)
	return &ce, nil
}

func UpdateComment(e *entity.Comment) error {
	_dp.Put(addPrefix(CommentBucket, e.ID), *e)
	return nil // TODO: check update errors
}

// =========== vote ================
func AddVote(e *entity.Vote) error {
	_dp.Put(addPrefix(VoteBucket, e.ID), *e)
	return nil
}

func GetVotes() []*entity.Vote {
	var res []*entity.Vote
	for _, value := range _dp.GetAll(VoteBucket) {
		temp, _ := value.(entity.Vote)
		res = append(res, &temp)
	}
	return res
}

func GetVote(id string) (*entity.Vote, error) {
	e, _ := _dp.Get(addPrefix(VoteBucket, id))
	ce, _ := e.(entity.Vote)
	return &ce, nil
}

func UpdateVote(e *entity.Vote) error {
	_dp.Put(addPrefix(VoteBucket, e.ID), *e)
	return nil // TODO: check update errors
}

/*
func (sdb *StatusDB) getAllEntities() []entity.Entity {
	var res []entity.Entity
	for _, data := range sdb.store.GetAll(bucket) {
		var e entity.Entity // TODo: need redesign here, NOT WORKING
		temp, _ := entity.DeserializeEntity(e, data)
		res = append(res, &temp)
	}
	return res
}*/

func (sdb *StatusDB) getEntityByID(bucket, key string, e entity.Entity) error {
	data, err := sdb.store.Get(bucket, key)
	if err != nil {
		return err
	}

	err = entity.DeserializeEntity(e, data)
	if err != nil {
		return err
	}

	return nil
}

func (sdb *StatusDB) createEntity(bucket string, e entity.Entity) error {
	if !entity.HasID(e) {
		return fmt.Errorf("create: entity doesn't have ID")
	}

	data, err := entity.SerializeEntity(e)
	if err != nil {
		return err
	}

	return sdb.store.Put(bucket, entity.GetID(e), data)
}

func (sdb *StatusDB) updateEntity(bucket string, e entity.Entity) error {
	if !entity.HasID(e) {
		return fmt.Errorf("update: entity doesn't have ID")
	}

	// TODO: check existence

	data, err := entity.SerializeEntity(e)
	if err != nil {
		return err
	}

	return sdb.store.Put(bucket, entity.GetID(e), data)
}
