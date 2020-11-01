package node

import (
	"io"

	"github.com/alexeyqian/gochain/protocol"
)

func (nd *Node) handleVerack(header *protocol.MessageHeader, conn io.ReadWriter) error {
	return nil
}
