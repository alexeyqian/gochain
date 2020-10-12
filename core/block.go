package core

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
	id string
	no int
}

type Chain struct{
	id string
	version string
}

func (b Block) String() string{
	return b.id
}