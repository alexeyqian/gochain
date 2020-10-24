package statusdb

import (
	"github.com/alexeyqian/gochain/config"
	"github.com/alexeyqian/gochain/entity"
)

const GpoKey = "gpo_1"
const GpoTable = "gpo"
const AccountTable = "account"
const ArticleTable = "article"
const CommentTable = "comment"
const VoteTable = "vote"

// used as DI for easy testing
// memory data provider for tesging
// file data provider for production
type DataProvider interface {
	Open()
	Close()
	Remove()

	GetAll(table string) []entity.Entity
	Get(key string) (entity.Entity, error)
	Put(key string, e entity.Entity) error
}

/*
type Repository struct{
	Open()
	Close()
	Remove()

	GetByID(id string) (*Account, error)
	Find() ([]*Account, error)
	Create(user *Account) error
	Update(user *Account) error
	Delete(id string) error

}*/

//var _lastSavedPoint int = 0 // used for fast replay from giving point
var _dp DataProvider

// NewDB(){

//}

// Open has parameter MemDataProvider
func Open() {
	if config.DataProvider == "MemDataProvider" {
		_dp = &MemDataProvider{}
	} else if config.DataProvider == "BoltDataProvider" {
		_dp = &BoltDataProvider{}
	} else {
		panic("Unknown data provider")
	}

	_dp.Open()
}

func Close() {
	_dp.Close()
}

func Remove() {
	_dp.Remove()
}

func GetGpo() (*entity.Gpo, error) {
	e, _ := _dp.Get(GpoKey)
	ce, _ := e.(entity.Gpo)
	return &ce, nil
}

func AddGpo(e *entity.Gpo) error {
	_dp.Put(GpoKey, *e)
	return nil
}

func UpdateGpo(e *entity.Gpo) error {
	_dp.Put(GpoKey, *e)
	return nil
}

func AddAccount(e *entity.Account) error {
	_dp.Put(addPrefix(AccountTable, e.Id), *e)
	return nil // TODO: check if account is already exist
}

func UpdateAccount(e *entity.Account) error {
	_dp.Put(addPrefix(AccountTable, e.Id), *e)
	return nil // TODO: check update errors
}

func GetAccounts() []*entity.Account {
	var res []*entity.Account
	for _, value := range _dp.GetAll(AccountTable) {
		temp, _ := value.(entity.Account)
		res = append(res, &temp)
	}
	return res
}

func GetAccount(id string) (*entity.Account, error) {
	e, err := _dp.Get(addPrefix(AccountTable, id))
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
	_dp.Put(addPrefix(ArticleTable, e.Id), *e)
	return nil
}

func GetArticles() []*entity.Article {
	var res []*entity.Article
	for _, value := range _dp.GetAll(ArticleTable) {
		temp, _ := value.(entity.Article)
		res = append(res, &temp)
	}
	return res
}

func GetArticle(id string) (*entity.Article, error) {
	e, _ := _dp.Get(addPrefix(ArticleTable, id))
	ce, _ := e.(entity.Article)
	return &ce, nil
}

func UpdateArticle(e *entity.Article) error {
	_dp.Put(addPrefix(ArticleTable, e.Id), *e)
	return nil // TODO: check update errors
}

func AddComment(e *entity.Comment) error {
	_dp.Put(addPrefix(CommentTable, e.Id), *e)
	return nil
}

func GetComments() []*entity.Comment {
	var res []*entity.Comment
	for _, value := range _dp.GetAll(CommentTable) {
		temp, _ := value.(entity.Comment)
		res = append(res, &temp)
	}
	return res
}

func GetComment(id string) (*entity.Comment, error) {
	e, _ := _dp.Get(addPrefix(CommentTable, id))
	ce, _ := e.(entity.Comment)
	return &ce, nil
}

func UpdateComment(e *entity.Comment) error {
	_dp.Put(addPrefix(CommentTable, e.Id), *e)
	return nil // TODO: check update errors
}

func AddVote(e *entity.Vote) error {
	_dp.Put(addPrefix(VoteTable, e.Id), *e)
	return nil
}

func GetVotes() []*entity.Vote {
	var res []*entity.Vote
	for _, value := range _dp.GetAll(VoteTable) {
		temp, _ := value.(entity.Vote)
		res = append(res, &temp)
	}
	return res
}

func GetVote(id string) (*entity.Vote, error) {
	e, _ := _dp.Get(addPrefix(VoteTable, id))
	ce, _ := e.(entity.Vote)
	return &ce, nil
}

func UpdateVote(e *entity.Vote) error {
	_dp.Put(addPrefix(VoteTable, e.Id), *e)
	return nil // TODO: check update errors
}

func addPrefix(table string, key string) string {
	return table + "_" + key
}
