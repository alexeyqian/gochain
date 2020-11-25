package protocol

import "math/rand"

type MsgPing struct {
	Nonce int
}

type MsgPong struct {
	Nonce int
}

func NewPingMsg(network string) (*Message, int, error) {
	nonce := rand.Uint64()
	payload := MsgPing{
		Nonce: nonce,
	}

	msg, err := NewMessage(network, cmdPing, payload)
	if err != nil {
		return nil, 0, err
	}

	return msg, nonce, nil
}

func NewPongMsg(network string, nonce int) (*Message, error) {
	payload := MsgPong{
		Nonce: nonce,
	}

	msg, err := NewMessage(network, cmdPong, payload)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
