package core

type Transactioner interface {
	TxType() string
}

const CreateAccountTransactionType = "CreateAccountTransactionType"
const TransferCoinTransactionType = "TransferCoinTransactionType"

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
