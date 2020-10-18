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

}

func Close() {
	_dp.Close()
}

func Remove() {
	_dp.Remove()
}

func GetGpo() *entity.Gpo {
	return _dp.Get(GpoKey).(entity.Gpo)
}

func SaveGpo(e *entity.Gpo) {
	_dp.Put(GpoKey, gpo)
}

func AddAccount(e entity.Account) {
	_dp.Put(addPrefix(AccountTable, e.Id), e)
}

func GetAccounts() []entity.Account {
	var res []entity.Account
	for key, value := range _dp.GetAll(AccountTable) {
		res = append(res, value.(entity.Account))
	}
	return res
}

func GetAccount(id string) *entity.Account {
	return _dp.Get(addPrefix(AccountTable, id)).(entity.Account)
}

func GetAccountByName(name string) (*entity.Account, error) {
	for index, acc := range GetAccounts() {
		if acc.Name == name {
			return &_accounts[index], nil
		}
	}
	return nil, nil
}

func AddArticle(e entity.Article) {
	_dp.Put(addPrefix(ArticleTable, e.Id), e)
}

func GetArticles() []entity.Article {
	var res []entity.Article
	for key, value := range _dp.GetAll(ArticleTable) {
		res = append(res, value.(entity.Article))
	}
	return res
}

func GetArticle(id string) *entity.Article {
	return _dp.Get(addPrefix(ArticleTable, id)).(entity.Article)
}

func AddComment(comment entity.Comment) {
	_dp.Put(addPrefix(CommentTable, e.Id), e)
}

func GetComments() []entity.Comment {
	var res []entity.Comment
	for key, value := range _dp.GetAll(CommentTable) {
		res = append(res, value.(entity.Comment))
	}
	return res
}

func GetComment(id string) *entity.Comment {
	return _dp.Get(addPrefix(CommentTable, id)).(entity.Comment)
}

func AddVote(r entity.Vote) {
	_dp.Put(addPrefix(VoteTable, e.Id), e)
}

func GetVotes() []entity.Vote {
	var res []entity.Vote
	for key, value := range _dp.GetAll(VoteTable) {
		res = append(res, value.(entity.Vote))
	}
	return res
}

func GetVote(id string) *entity.Vote {
	return _dp.Get(addPrefix(VoteTable, id)).(entity.Vote)
}

func addPrefix(table string, key string) {
	return table + "_" + key
}
