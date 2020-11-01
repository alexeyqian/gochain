package node

import (
	"fmt"
	"io"

	"github.com/alexeyqian/gochain/protocol"
	"github.com/alexeyqian/gochain/utils"
)

func (nd *Node) handleTx(header *protocol.MessageHeader, conn io.ReadWriter) error {
	var tx protocol.MsgTx

	lr := io.LimitReader(conn, int64(header.Length))
	utils.DeserializeWithReader(&tx, lr)

	fmt.Printf("transaction : %+v\n", tx)

	nd.mempool.NewTxCh <- tx
	return nil
}
