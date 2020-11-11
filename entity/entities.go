package entity

import (
	"reflect"

	"github.com/alexeyqian/gochain/utils"
)

type Entity interface {
}

type Gpo struct {
	ID       string
	BlockId  string
	BlockNum uint64
	Witness  string
	Time     uint64
	Version  string
	Supply   uint64
}

type Wso struct {
	ID                 string
	MajorVersion       int
	MaxBlockSize       int
	AccountCreationFee int
	CurrentWitnesses   []string
}

type Witness struct {
	ID      string
	Name    string
	Votes   int
	Version int
}

type Account struct {
	ID           string
	Name         string
	CreatedOn    uint64
	Coin         uint64
	Vest         uint64
	Credit       uint64
	UpVotes      uint32
	DownVotes    uint32
	VotePower    uint64
	ArticleCount uint64
}

type Article struct {
	ID        string
	Author    string
	Title     string
	Body      string
	Meta      string
	UpVotes   uint32
	DownVotes uint32
	VotePower uint64
}

type Comment struct {
	ID        string
	ParentId  string
	CommentId string
	Commentor string
	Body      string
	CreatedOn uint64
	UpVotes   uint32
	DownVotes uint32
	VotePower uint64
}

type Vote struct {
	ID         string
	ParentId   string
	ParentType string
	Direction  int8
	VotePower  uint64
	Voter      string
}

func HasID(e Entity) bool {
	//fmt.Printf("reflect %+v", reflect.ValueOf(e).Elem())
	return reflect.ValueOf(e).FieldByName("ID").String() != ""
}

func GetID(e Entity) string {
	return reflect.ValueOf(e).FieldByName("ID").String()
}

func Serialize(e Entity) []byte {
	return utils.Serialize(e)
}

func Deserialize(e Entity, data []byte) {
	utils.Deserialize(e, data)
}
