package entity

import (
	"bytes"
	"encoding/gob"
	"reflect"
)

type Entity interface {
}

type Gpo struct {
	BlockId  string
	BlockNum uint64
	Witness  string
	Time     uint64
	Version  string
	Supply   uint64
}

type Witness struct {
	ID   string
	Name string
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
	return reflect.ValueOf(e).Elem().FieldByName("ID").String() != ""
}

func GetID(e Entity) string {
	return reflect.ValueOf(e).Elem().FieldByName("ID").String()
}

func Serialize(e Entity) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(e)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func Deserialize(e Entity, data []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(data))
	return dec.Decode(e)
}
