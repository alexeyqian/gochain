package node

import (
	"io"
	"log"
)

func (nd Node) handleTx(header *protocol.MessageHeader, conn io.ReadWriter) error {
	var tx protocol.MsgTx

	lr := io.LimitReader(conn, int64(header.Length))
	if err := binary.NewDecoder(lr).Decode(&tx); err != nil {
		return err
	}
	log.Printf("transaction : %+v\n", tx)

	nd.mempool.NewTxCh <- tx
	return nil
}
