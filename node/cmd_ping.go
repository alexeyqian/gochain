package node

import (
	"io"

	"github.com/alexeyqian/gochain/protocol"
	"github.com/alexeyqian/gochain/utils"
)

func (nd *Node) handlePing(header *protocol.MessageHeader, conn io.ReadWriter) error {
	var ping protocol.MsgPing

	lr := io.LimitReader(conn, int64(header.Length))
	utils.DeserializeWithReader(&ping, lr)

	// after receiving a ping
	// send out a pong response
	pong, err := protocol.NewPongMsg(nd.Network, ping.Nonce)
	if err != nil {
		return err
	}

	msg := utils.Serialize(pong)

	if _, err := conn.Write(msg); err != nil {
		return err
	}

	return nil
}
