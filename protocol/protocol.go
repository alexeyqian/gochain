package protocol

const (
	Version   = 70015
	UserAgent = "/Satoshi:5.64/tinybit:0.0.1/"

	SrvNodeNetwork        = 1
	SrvNodeGetUTXO        = 2
	SrvNodeBloom          = 4
	SrvNodeWitness        = 8
	SrvNodeNetworkLimited = 1024
)

func NewUserAgent() VarStr {
	return newVarStr(UserAgent)
}
