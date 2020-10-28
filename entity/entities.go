package entity

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
	UpVotes      uint32
	DownVotes    uint32
	VotePower    uint64
	ArticleCount uint64
}

type Article struct {
	Id        string
	Author    string
	Title     string
	Body      string
	Meta      string
	UpVotes   uint32
	DownVotes uint32
	VotePower uint64
}

type Comment struct {
	Id        string
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
	Id         string
	ParentId   string
	ParentType string
	Direction  int8
	VotePower  uint64
	Voter      string
}

// TODO: use reflection now, will be replaced with generics later
func HasID(e Entity) bool {
	return true
}

func GetEntityType(e Entity) string {
	// using reflection
	return ""
}
