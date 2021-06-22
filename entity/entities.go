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
	BlockNum int
	Witness  string
	Time     int
	Version  string
	Supply   int
}

type Wso struct {
	ID                 string
	MajorVersion       int
	MaxBlockSize       int
	AccountCreationFee int
	CurrentWitnesses   []string
}

type SoftForkItem struct {
	ID          string
	BlockNum    int
	BlockData   []byte
	PrevBlockID string
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
	CreatedOn    int
	Coin         int
	Vest         int
	Credit       int
	UpVotes      int
	DownVotes    int
	VotePower    int
	ArticleCount int
}

type Article struct {
	ID        string
	Author    string
	Title     string
	Body      string
	Meta      string
	UpVotes   int
	DownVotes int
	VotePower int
}

type Comment struct {
	ID        string
	ParentId  string
	CommentId string
	Commentor string
	Body      string
	CreatedOn int
	UpVotes   int
	DownVotes int
	VotePower int
}

type Vote struct {
	ID         string
	ParentId   string
	ParentType string
	Direction  int
	VotePower  int
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
