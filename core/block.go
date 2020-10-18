package core

import (
	"crypto/sha256"
	"encoding/json"
)

const MaxTransactionsInBlock = 10

// TODO: use array to hold block id/hash, and map to hold id => block
// for fast access block id/hash and data by array index

type Block struct {
	Id           string // 32 bytes (256 bits) hash value of the entire block
	PrevBlockId  string
	Num          uint64
	MerkleRoot   string
	CreatedOn    uint64
	Witness      string
	nonce        uint64
	Transactions []Transactioner
	//Size         uint64  ??

}

func (b *Block) SerializeBlockWitoutId() []byte {
	// serialize all fields except Id into byte array
	return nil
}

func (b *Block) CalculateHash() {
	var data []byte
	// fill data with block field bytes
	hash := sha256.Sum256(data)
	var temp []byte = hash[:]
	b.Id = string(temp)
}

func SerializeBlock(b *Block) ([]byte, error) {
	return json.Marshal(*b)
}

func UnSerializeBlock(data []byte) (*Block, error) {
	var b Block
	json.Unmarshal(data, &b)
	return &b, nil
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
