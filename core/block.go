package core
import (
	"encoding/json"
)

type GlobalProperties struct{
	bno int
	genesis_time string
	last_witness string
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

type Transaction struct{
	id string
	txtype string
	data string
	creator string
	createdon int
	signature string
}

type Block struct{
	Id string
	Num int
}

type Chain struct{
	id string
	version string
}

func (b Block) String() string{
	return b.Id
}

func SerializeBlock(b *Block) []byte{
	data, _ := json.Marshal(*b)
	return data
}

func UnSerializeBlock(data []byte) Block{
	var b Block
	json.Unmarshal(data, &b)
	return b
}