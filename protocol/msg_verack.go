package protocol

func NewVerackMsg(network string) (*Message, error) {
	msg, err := NewMessage(cmdVerack, network, []byte{})
	if err != nil {
		return nil, err
	}

	return msg, nil
}
