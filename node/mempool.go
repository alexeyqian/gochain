package node

import (
	"github.com/alexeyqian/gochain/protocol"
)

type Mempool struct{
	NewBlockCh chan protocol.MsgBlock
	NewTxCh chain protocol.MsgTx

	txs map[string]*protocol.MsgTx
}

func (m Mempool) Run(){
	for{
		select{
		case tx := <- m.NewTxCh:
			hash, err := tx.Hash()
			txid := hex.EncodeToString(hash)
			m.txs[txid] = &tx
		}
	}
}