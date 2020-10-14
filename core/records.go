package core

const InitWitness = "init"
const InitAmount = 100
const AmountPerBlock = 100

type Gpo struct {
	PrevBlockId string
	BlockId     string
	BlockNum    uint64
	Witness     string
	Time        uint64
	Version     string
	Supply      uint64
}

type Witness struct {
	Id   string
	Name string
}

type Account struct {
	Id        string
	Name      string
	CreatedOn uint64
	Coin      uint64
	Vest      uint64
	Credit    uint64
}

type Article struct {
	Author   string
	Title    string
	Content  string
	Premlink string
}
