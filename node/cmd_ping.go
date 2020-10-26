package node

import (
	"io"

	"github.com/alexeyqian/gochain/protocol"
)

func (nd Node) handlePing(header *protocol.MessageHeader, conn io.ReadWriter) error {
	var ping protocol.MsgPing

	lr := io.LimitReader(comm, int64(header.Length))
	if err := binary.NewDecoder(lr).Decode(&ping); err != nil {
		return err
	}

	pong, err := protocol.NewPongMsg(n.Network, ping.Nonce)
	msg, err := binary.Marshal(pong)

	if _, err := conn.Write(msg); err != nil {
		return err
	}

	return nil
}
