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

func (sdb *StatusDB) CreateGpo(e *entity.Gpo) error {
	return sdb.createEntity(GpoBucket, *e)
}

func (sdb *StatusDB) UpdateGpo(e *entity.Gpo) error {
	return sdb.updateEntity(GpoBucket, *e)
}

func (sdb *StatusDB) GetGpo() (*entity.Gpo, error) {
	var e Gpo
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
	for _, value := range sdb.store.GetAll(AccountBucket) {
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

func CreateArticle(e *entity.Article) error {
	return sdb.createEntity(ArticleBucket, *e)
}

func UpdateArticle(e *entity.Article) error {
	return sdb.updateEntity(ArticleBucket, *e)
}

func GetArticle(id string) (*entity.Article, error) {
	var e entity.Article
	err := sdb.getEntityByID(ArticleBucket, id, e)
	return &e, err
}

func GetArticles() []*entity.Article {
	var res []*entity.Article
	for _, value := range sdb.store.GetAll(ArticleBucket) {
		var e entity.Article
		entity.Deserialize(e, value)
		res = append(res, &e)
	}
	return res
}

// =========== comment ================

func CreateComment(e *entity.Comment) error {
	return sdb.createEntity(CommentBucket, *e)
}

func UpdateComment(e *entity.Comment) error {
	return sdb.updateEntity(CommentBucket, *e)
}

func GetComment(id string) (*entity.Comment, error) {
	var e entity.Comment
	err := sdb.getEntityByID(CommentBucket, id, e)
	return &e, err
}

func GetComments() []*entity.Comment {
	var res []*entity.Comment
	for _, value := range sdb.store.GetAll(CommentBucket) {
		var e entity.Comment
		entity.Deserialize(e, value)
		res = append(res, &e)
	}
	return res
}

// =========== vote ================

func CreateVote(e *entity.Vote) error {
	return sdb.createEntity(VoteBucket, *e)
}

func UpdateVote(e *entity.Vote) error {
	return sdb.updateEntity(VoteBucket, *e)
}

func GetVote(id string) (*entity.Vote, error) {
	var e entity.Vote
	err := sdb.getEntityByID(VoteBucket, id, e)
	return &e, err
}

func GetVotes() []*entity.Vote {
	var res []*entity.Vote
	for _, value := range sdb.store.GetAll(VoteBucket) {
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

	data, err := entity.Serialize(e)
	if err != nil {
		return err
	}
	// sdb.sets[bucket].Create(id, data)
	return sdb.store.Put(bucket, entity.GetID(e), data)
}

func (sdb *StatusDB) updateEntity(bucket string, e entity.Entity) error {
	if !entity.HasID(e) {
		return fmt.Errorf("update: entity doesn't have ID")
	}

	// TODO: check existence

	data, err := entity.Serialize(e)
	if err != nil {
		return err
	}

	return sdb.store.Put(bucket, entity.GetID(e), data)
}

func (sdb *StatusDB) getEntityByID(bucket, key string, e entity.Entity) error {
	data, err := sdb.store.Get(bucket, key)
	if err != nil {
		return err
	}

	err = entity.Deserialize(e, data)
	if err != nil {
		return err
	}

	return nil
}

/*
func (sdb *StatusDB) getAllEntities() []entity.Entity {
	var res []entity.Entity
	for _, data := range sdb.store.GetAll(bucket) {
		var e entity.Entity // TODo: need redesign here, NOT WORKING
		temp, _ := entity.Deserialize(e, data)
		res = append(res, &temp)
	}
	return res
}*/
