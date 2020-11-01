package node

import (
	"encoding/hex"
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
			// TODO: validation
			hash, err := tx.Hash()
			txid := hex.EncodeToString(hash)
			m.txs[txid] = &tx
		case block := <- m.NewBlockCh:
			// TODO: validation
			for _, tx := range block.Txs{
				hash, err := tx.Hash()

				txid := hex.EncodeToString(hash)
				delete(m.txs, txid)
			}
		}

	}
}