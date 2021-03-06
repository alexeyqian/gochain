package core

import (
	"crypto/sha256"
	"encoding/json"
)

const MaxTransactionsInBlock = 10

// TODO: use array to hold block id/hash, and map to hold id => block
// for fast access block id/hash and data by array index

type Block struct {
	ID           string // 32 bytes (256 bits) hash value of the entire block
	Num          int
	PrevBlockId  string
	MerkleRoot   string
	CreatedOn    int
	Witness      string
	nonce        int
	Transactions []Transactioner
	//Size         int  ??
}

func (b *Block) SerializeBlockWitoutId() []byte {
	// serialize all fields except ID into byte array
	return nil
}

func (b *Block) CalculateHash() {
	var data []byte
	// fill data with block field bytes
	hash := sha256.Sum256(data)
	var temp []byte = hash[:]
	b.ID = string(temp)
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

func (b *Block) HashTransactions() []byte {
	var transactions [][]byte

	for _, tx := range b.Transactions {
		transactions = append(transactions, HashTx(tx))
	}

	mtree := NewMerkleTree(transactions)
	return mtree.RootNode.Data
}
