package statusdb

import (
	"fmt"

	"github.com/alexeyqian/gochain/entity"
	"github.com/alexeyqian/gochain/utils"
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

func (sdb *StatusDB) GetGpo() (*entity.Gpo, error) {
	data, err := sdb.store.Get(GpoBucket, GpoKey)
	if err != nil {
		return nil, err
	}

	var e Gpo
	e, err = utils.Deserialize(&e, data)
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (sdb *StatusDB) CreateGpo(e *entity.Gpo) error {
	data, err := utils.SerializeEntity(e)
	if err != nil {
		return err
	}

	return sdb.store.Put(GpoBucket, GpoKey, *e)
}

func (sdb *StatusDB) UpdateGpo(e *entity.Gpo) error {
	data, err := utils.SerializeEntity(e)
	if err != nil {
		return err
	}

	return sdb.store.Put(GpoBucket, GpoKey, *e)
}

func (sdb *StatusDB) CreateAccount(e *entity.Account) error {
	// TODO: validate if account is already exist
	_dp.Put(addPrefix(AccountBucket, e.Id), *e)
	return nil // TODO: check if account is already exist
}

func UpdateAccount(e *entity.Account) error {
	_dp.Put(addPrefix(AccountBucket, e.Id), *e)
	return nil // TODO: check update errors
}

func GetAccounts() []*entity.Account {
	var res []*entity.Account
	for _, value := range _dp.GetAll(AccountBucket) {
		temp, _ := value.(entity.Account)
		res = append(res, &temp)
	}
	return res
}

func GetAccount(id string) (*entity.Account, error) {
	e, err := _dp.Get(addPrefix(AccountBucket, id))
	if err != nil {
		return nil, err
	}
	ce, _ := e.(entity.Account)
	return &ce, nil
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

func AddArticle(e *entity.Article) error {
	_dp.Put(addPrefix(ArticleBucket, e.Id), *e)
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
	_dp.Put(addPrefix(ArticleBucket, e.Id), *e)
	return nil // TODO: check update errors
}

func AddComment(e *entity.Comment) error {
	_dp.Put(addPrefix(CommentBucket, e.Id), *e)
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
	_dp.Put(addPrefix(CommentBucket, e.Id), *e)
	return nil // TODO: check update errors
}

func AddVote(e *entity.Vote) error {
	_dp.Put(addPrefix(VoteBucket, e.Id), *e)
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
	_dp.Put(addPrefix(VoteBucket, e.Id), *e)
	return nil // TODO: check update errors
}

func addPrefix(table string, key string) string {
	return table + "_" + key
}

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
		return fmt.Errorf("entity doesn't have ID")
	}

	data, err := entity.SerializeEntity(e)
	if err != nil {
		return err
	}

	return sdb.store.Put(bucket, entity.GetID(e), data)
}
