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

var _lastSavedPoint int = 0 // used for fast replay from giving point
var _dp DataProvider

// Open has parameter MemDataProvider
func Open() {
	if config.DataProvider == "MemDataProvider" {
		_dp = &MemDataProvider{}
		//} else if config.DataProvider == "BoldDataProvider" {
		//	_dp = BoldDataProvider{}
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
	ce := e.(entity.Gpo)
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
		temp := value.(entity.Account)
		res = append(res, &temp)
	}
	return res
}

func GetAccount(id string) (*entity.Account, error) {
	e, err := _dp.Get(addPrefix(AccountTable, id))
	if err != nil {
		return nil, err
	}
	ce := e.(entity.Account)
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

func AddArticle(e *entity.Article) {
	_dp.Put(addPrefix(ArticleTable, e.Id), *e)
}

func GetArticles() []*entity.Article {
	var res []*entity.Article
	for _, value := range _dp.GetAll(ArticleTable) {
		temp := value.(entity.Article)
		res = append(res, &temp)
	}
	return res
}

func GetArticle(id string) *entity.Article {
	e, _ := _dp.Get(addPrefix(ArticleTable, id))
	ce := e.(entity.Article)
	return &ce
}

func AddComment(e *entity.Comment) {
	_dp.Put(addPrefix(CommentTable, e.Id), *e)
}

func GetComments() []*entity.Comment {
	var res []*entity.Comment
	for _, value := range _dp.GetAll(CommentTable) {
		temp := value.(entity.Comment)
		res = append(res, &temp)
	}
	return res
}

func GetComment(id string) *entity.Comment {
	e, _ := _dp.Get(addPrefix(CommentTable, id))
	ce := e.(entity.Comment)
	return &ce
}

func AddVote(e *entity.Vote) {
	_dp.Put(addPrefix(VoteTable, e.Id), *e)
}

func GetVotes() []*entity.Vote {
	var res []*entity.Vote
	for _, value := range _dp.GetAll(VoteTable) {
		temp := value.(entity.Vote)
		res = append(res, &temp)
	}
	return res
}

func GetVote(id string) *entity.Vote {
	e, _ := _dp.Get(addPrefix(VoteTable, id))
	ce := e.(entity.Vote)
	return &ce
}

func addPrefix(table string, key string) string {
	return table + "_" + key
}
