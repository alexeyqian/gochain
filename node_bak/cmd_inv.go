package node

/*
func (nd *Node) handleInv(request []byte) {
	var buff bytes.Buffer
	var payload InvRequest

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		fmt.Println("cannot decode payload")
		return
	}

	// processing ...

}*/

/*
func (nd Node) handleInv(header *protocol.MessageHeader, conn io.ReadWriter) error {
	var inv protocol.MsgInv

	lr := io.LimitReader(conn, int64(header.Length))
	utils.DeserializeWithReader(&inv, lr)

	var getData protocol.MsgGetData
	getData.Inventory = inv.Inventory
	getData.Count = inv.Count

	getDataMsg, err := protocol.NewMessage(nd.Network, "getdata", getData)
	if err != nil {
		return err
	}

	msg := utils.Serialize(getDataMsg)

	if _, err := conn.Write(msg); err != nil {
		return err
	}
}
*/
