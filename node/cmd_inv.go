package node

import (
	"io"
)

func (nd Node) handleInv(header *protocol.MessageHeader, conn io.ReadWriter) error {
	var inv protocol.MsgInv

	lr := io.LimitReader(conn, int64(header.Length))
	if err := binary.NewDecoder(lr).Decode(&inv); err != nil {
		return err
	}

	var getData protocol.MsgGetData
	getData.Inventory = inv.Inventory
	getData.Count = inv.Count

	getDataMsg, err := protocol.NewMessage("getdata", nd.Network, getData)
	if err != nil {
		return err
	}

	msg, err := binary.Marshal(getDataMsg)
	if err != nil {
		return err
	}

	if _, err := conn.Write(msg); err != nil {
		return err
	}
}