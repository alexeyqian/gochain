package node

import (
	"io"

	"github.com/alexeyqian/gochain/protocol"
	"github.com/alexeyqian/gochain/utils"
)

func (nd *Node) handlePong(header *protocol.MessageHeader, conn io.ReadWriter) error {
	var pong protocol.MsgPong

	lr := io.LimitReader(conn, int64(header.Length))
	err := utils.DeserializeWithReader(&pong, lr)
	if err != nil {
		return err
	}

	nd.PongCh <- pong.Nonce
	return nil
}
