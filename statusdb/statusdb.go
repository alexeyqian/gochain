package statusdb

import (
	"github.com/alexeyqian/gochain/entity"
)

var _lastSavedPoint int = 0 // used for fast replay from giving point
var _gpo entity.Gpo
var _accounts []entity.Account
var _articles []entity.Article
var _comments []entity.Comment
var _votes []entity.Vote

func Open() {

}

func Close() {
	// persistent data to disk
}

func Remove() {
	// reset all data
	_gpo = entity.Gpo{}
	_accounts = nil
}

func GetGpo() *entity.Gpo {
	return &_gpo
}

func AddAccount(acc entity.Account) {
	_accounts = append(_accounts, acc)
}

func GetAccounts() []entity.Account {
	return _accounts
}

func GetAccountByName(name string) *entity.Account {
	for index, acc := range _accounts {
		if acc.Name == name {
			return &_accounts[index]
		}
	}
	return nil
}

func GetAccount(id string) *entity.Account {
	for index, acc := range _accounts {
		if acc.Id == id {
			return &_accounts[index]
		}
	}
	return nil
}

func AddArticle(article entity.Article) {
	_articles = append(_articles, article)
}

func GetArticles() []entity.Article {
	return _articles
}

// TODO: use map(hash/id, object for fast access of an object)
func GetArticle(id string) *entity.Article {
	for index, article := range _articles {
		if article.ArticleId == id {
			return &_articles[index]
		}
	}
	return nil
}

func AddComment(comment entity.Comment) {
	_comments = append(_comments, comment)
}

func GetComments() []entity.Comment {
	return _comments
}

// TODO: use general function to remove duplicated code
func GetComment(id string) *entity.Comment {
	for index, comment := range _comments {
		if comment.CommentId == id {
			return &_comments[index]
		}
	}
	return nil
}

func AddVote(r entity.Vote) {
	_votes = append(_votes, r)
}

func GetVotes() []entity.Vote {
	return _votes
}

func GetVote(id string) *entity.Vote {
	for index, vote := range _votes {
		if vote.Id == id {
			return &_votes[index]
		}
	}
	return nil
}
