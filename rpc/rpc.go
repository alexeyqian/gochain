package rpc

import (
	"fmt"

	"github.com/alexeyqian/gochain/protocol"
)

type Node interface {
	Mempool() map[string]*protocoL.MsgTx
}

type RPC struct {
	node Node
}

type MempoolArgs interface{}
type MempoolReply string

func (r RPC) GetMempool(args *MempoolArgs, replay *MempoolReply) error {
	txs := r.node.Mempool()
	*reply = MempoolReply(formatMempoolReplay(txs))

	return nil
}

func formatMempoolReply(txs map[string]*protocol.MsgTx) string {
	var result string
	for k := range txs {
		result += fmt.Sprintf("%s\n", k)
	}

	result += fmt.Sprintf("Total %d txs.", len(txs))

	return result
}
