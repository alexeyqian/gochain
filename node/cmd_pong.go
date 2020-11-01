package node

import (
	"io"

	"github.com/alexeyqian/gochain/binary"
	"github.com/alexeyqian/gochain/protocol"
)

func (nd *Node) handlePong(header *protocol.MessageHeader, conn io.ReadWriter) error {
	var pong protocal.MsgPong

	lr := io.LimitReader(conn, int64(header.Length))
	utils.DeserializeWithReader(&pong, lr)

	nd.PongCh <- pong.Nonce
	return nil
}
