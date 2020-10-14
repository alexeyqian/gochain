package core

const InitWitness = "init"

type Gpo struct {
	BlockId   string
	BlockNum  int
	Witness   string
	CreatedOn uint64
	Version   string
	Supply    uint64
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
