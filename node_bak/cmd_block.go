package node

import (
	"io"

	"github.com/alexeyqian/gochain/protocol"
	"github.com/alexeyqian/gochain/utils"
)

func (nd *Node) handleBlock(header *protocol.MessageHeader, conn io.ReadWriter) error {
	var block protocol.MsgBlock

	lr := io.LimitReader(conn, int64(header.Length))
	err := utils.DeserializeWithReader(&block, lr)
	if err != nil {
		return err
	}

	nd.mempool.NewBlockCh <- block
	return nil
}
