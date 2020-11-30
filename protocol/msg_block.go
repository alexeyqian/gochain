package protocol

type MsgBlock struct {
	Version    int
	PrevBlock  [32]byte
	MerkleRoot [32]byte
	Timestamp  int
	Nonce      int
	TxCount    int
	Txs        []MsgTx
}
