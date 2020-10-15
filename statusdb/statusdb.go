package statusdb

import (
	"github.com/alexeyqian/gochain/core"
)

var _lastSavedPoint int = 0 // used for fast replay from giving point
var _gpo core.Gpo
var _accounts []core.Account
var _pendingTransactions []core.Transactioner
var _articles []core.Article

func Open() {

}

func Close() {
	// persistent data to disk
}

func Remove() {
	// reset all data
	_gpo = core.Gpo{}
	_accounts = nil
	_pendingTransactions = nil
}

func GetGpo() *core.Gpo {
	return &_gpo
}

func GetPendingTransactions() []core.Transactioner {
	return _pendingTransactions
}

// TODO: should only be alowwed to invoke from chain package
func AddPendingTransaction(tx core.Transactioner) {
	_pendingTransactions = append(_pendingTransactions, tx)
}

func MovePendingTransactionsToBlock(b *core.Block) {
	i := 0
	for _, tx := range _pendingTransactions {
		if i >= core.MaxTransactionsInBlock {
			break
		}
		b.AddTransaction(tx)
		i++
	}

	if len(_pendingTransactions) > core.MaxTransactionsInBlock {
		_pendingTransactions = _pendingTransactions[core.MaxTransactionsInBlock:]
	}
}

func AddAccount(acc core.Account) {
	_accounts = append(_accounts, acc)
}

func GetAccounts() []core.Account {
	return _accounts
}

func GetAccount(name string) *core.Account {
	for index, acc := range _accounts {
		if acc.Name == name {
			return &_accounts[index]
		}
	}
	return nil
}

func AddArticle(article core.Article) {
	_articles = append(_articles, article)
}

func GetArticles() []core.Article {
	return _articles
}

// TODO: use map(hash/id, object for fast access of an object)
func GetArticle(id string) *core.Article {
	for index, article := range _articles {
		if article.ArticleId == id {
			return &_articles[index]
		}
	}
	return nil
}
