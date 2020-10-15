package core

const InitWitness = "init"
const InitAmount = 100
const AmountPerBlock = 100
const BlockInterval = 3 // seconds
const BlockZeroId = "00000000-0000-0000-000000000000"
const GenesisTime = 1632830400 //Date and time (GMT): Tuesday, September 28, 2021 12:00:00 PM

type Gpo struct {
	BlockId  string
	BlockNum uint64
	Witness  string
	Time     uint64
	Version  string
	Supply   uint64
}

type Witness struct {
	Id   string
	Name string
}

type Account struct {
	Id           string
	Name         string
	CreatedOn    uint64
	Coin         uint64
	Vest         uint64
	Credit       uint64
	ArticleCount uint64
}

type Article struct {
	ArticleId string
	Author    string
	Title     string
	Body      string
	Meta      string
}
