package core

type Transactioner interface {
	Evaluate()
}

type CreateAccountTransaction struct {
	AccountId   string
	AccountName string
	CreatedBy   string
	CreatedOn   uint64
	ExpiredOn   uint64
	Signature   string // should be [SIGBITS]byte
}

func (t CreateAccountTransaction) Evaluate() {
	return
}

func SignTransaction(tx *Transactioner) {

}
