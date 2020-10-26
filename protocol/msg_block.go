package protocol

type MsgBlock struct {
	Version    int32
	PrevBlock  [32]byte
	MerkleRoot [32]byte
	Timestamp  uint32
	Bits       [4]byte
	Nonce      uint32
	TxCount    uint8
	Txs        []MsgTx
}
