package core

type Transactioner interface {
	Apply() error
	Validate() error
}

const InitWitness = "init"
const InitAmount = 100
const AmountPerBlock = 100
const BlockInterval = 3 // seconds
const BlockZeroId = "00000000-0000-0000-000000000000"
const GenesisTime = 1632830400 //Date and time (GMT): Tuesday, September 28, 2021 12:00:00 PM

const VoteParentTypeArticle = "VoteParentTypeArticle"
const VoteParentTypeComment = "VoteParentTypeComment"
const VoteParentTypeAccount = "VoteParentTypeAccount"

type CreateAccountTransaction struct {
	Id          string
	AccountId   string
	AccountName string
	CreatedBy   string
	Signature   string // should be [SIGBITS]byte
	CreatedOn   uint64
}

type TransferCoinTransaction struct {
	Id        string
	From      string
	To        string
	Amount    uint64
	CreatedOn uint64
}

type CreateArticleTransaction struct {
	Id        string
	ArticleId string
	Author    string
	Title     string
	Body      string
	Meta      string
	CreatedOn uint64
}

type CreateCommentTransaction struct {
	Id        string
	ParentId  string
	CommentId string
	Commentor string
	Body      string
	CreatedOn uint64
}

type VoteTransaction struct {
	Id         string
	ParentId   string
	ParentType string
	Direction  int8
	VotePower  uint64
	Voter      string
	CreatedOn  uint64
}
