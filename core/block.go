package core

import (
	"encoding/json"
)

const MaxTransactionsInBlock = 10

type Block struct {
	Id           string
	Num          uint64
	CreatedOn    uint64
	Witness      string
	Transactions []Transactioner
}

func SerializeBlock(b *Block) []byte {
	data, _ := json.Marshal(*b)
	return data
}

func UnSerializeBlock(data []byte) Block {
	var b Block
	json.Unmarshal(data, &b)
	return b
}

func (b *Block) AddTransaction(t Transactioner) {
	b.Transactions = append(b.Transactions, t)
}

/*

type GlobalProperties struct {
	bno            int
	genesis_time   string
	last_witness   string
	current_supply int
	max_block_size int
}

type TransactionType int

const (
	TtCreateAccount TransactionType = iota
	TtTransferCoin
	TtCreateArticle
	TtRegisterWitness
)

type Transaction struct {
	id        string
	txtype    string
	data      string
	creator   string
	createdon int
	signature string
}


type Chain struct {
	id      string
	version string
}

*/
