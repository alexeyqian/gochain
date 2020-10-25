package protocol

const (
	Version        = 70015
	UserAgent      = "/weku:0.0.1"
	SrvNoteNetwork = 1
	SrvNodeWitness = 2
	SrvNodeFull    = 3
)

func NewUserAgent() VarStr {
	return newVarStr(UserAgent)
}
