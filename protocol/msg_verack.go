package protocol

func NewVerackMsg(network string) (*Message, error) {
	msg, err := NewMessage(network, cmdVerack, []byte{})
	if err != nil {
		return nil, err
	}

	return msg, nil
}
