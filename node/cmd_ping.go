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

	// after receiving a ping
	// send out a pong response
	pong, err := protocol.NewPongMsg(n.Network, ping.Nonce)
	if err != nil {
		return err
	}

	msg, err := binary.Marshal(pong)
	if err != nil {
		return err
	}

	if _, err := conn.Write(msg); err != nil {
		return err
	}

	return nil
}
