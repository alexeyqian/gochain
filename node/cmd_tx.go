package node

import "io"

func (nd Node) handleTx(header *protocol.MessageHeader, conn io.ReadWriter) error {
	var tx protocol.MsgTx

	lr := io.LimitReader(conn, int64(header.Length))
	if err := binary.NewDecoder(lr).Decode(&tx); err != nil {
		return err
	}
	// log.Debugf("transaction : %+v", tx)
	return nil
}
