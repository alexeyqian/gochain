package node

import (
	"io"
)

func (nd Node) handlePong(header *protocol.MessageHeader, conn io.ReadWriter) error {
	var pong protocal.MsgPong

	lr := io.LimitReader(conn, int64(header.Length))
	if err := binary.NewDecoder(lr).Decode(&pong); err != nil {
		return err
	}

	nd.PongCh <- pong.Nonce
	return nil
}
