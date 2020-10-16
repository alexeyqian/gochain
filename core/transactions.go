package core

//	TODO: move core as higher level -> core
// move records to lower db level -> statusdb

type Transactioner interface {
	TxType() string
}

const InitWitness = "init"
const InitAmount = 100
const AmountPerBlock = 100
const BlockInterval = 3 // seconds
const BlockZeroId = "00000000-0000-0000-000000000000"
const GenesisTime = 1632830400 //Date and time (GMT): Tuesday, September 28, 2021 12:00:00 PM

// TODO: move consts to TransactionType Enum
const CreateAccountTransactionType = "CreateAccountTransactionType"
const TransferCoinTransactionType = "TransferCoinTransactionType"
const CreateArticleTransactionType = "CreateArticleTransactionType"
const CreateCommentTransactionType = "CreateCommentTransactionType"
const VoteTransactionType = "VoteTransactionType"

const VoteParentTypeArticle = "VoteParentTypeArticle"
const VoteParentTypeComment = "VoteParentTypeComment"
const VoteParentTypeAccount = "VoteParentTypeAccount"

type CreateAccountTransaction struct {
	AccountId   string
	AccountName string
	CreatedBy   string
	CreatedOn   uint64
	ExpiredOn   uint64
	Signature   string // should be [SIGBITS]byte
}

func (tx CreateAccountTransaction) TxType() string {
	return CreateAccountTransactionType
}

type TransferCoinTransaction struct {
	From   string
	To     string
	Amount uint64
}

func (tx TransferCoinTransaction) TxType() string {
	return TransferCoinTransactionType
}

type CreateArticleTransaction struct {
	ArticleId string
	Author    string
	Title     string
	Body      string
	Meta      string
	CreatedOn uint64
}

func (tx CreateArticleTransaction) TxType() string {
	return CreateArticleTransactionType
}

type CreateCommentTransaction struct {
	ParentId  string
	CommentId string
	Commentor string
	Body      string
	CreatedOn uint64
}

func (tx CreateCommentTransaction) TxType() string {
	return CreateCommentTransactionType
}

type VoteTransaction struct {
	Id         string
	ParentId   string
	ParentType string
	Direction  int8
	VotePower  uint64
	Voter      string
}

func (tx VoteTransaction) TxType() string {
	return VoteTransactionType
}
