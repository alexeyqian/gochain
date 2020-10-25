package protocol

import (
	"bytes"
	"encoding/binary"
)

type MsgVersion struct {
	Version     int32
	Services    uint64
	Timestamp   int64
	AddrRecv    NetAddr
	AddrFrom    NetAddr
	Nonce       uint64
	UserAgent   VarStr
	StartHeight int32
	Relay       bool
}

func (v MsgVersion) Serialize() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.BigEndian, v.Version); err != nil {
		return nil, err
	}

	if err := binary.Write(&buf, binary.BigEndian, v.Services); err != nil {
		return nil, err
	}

	if err := binary.Write(&buf, binary.BigEndian, v.Timestamp); err != nil {
		return nil, err
	}

	serializedAddrRecv, err := v.AddrRecv.Serialize()
	if err != nil {
		return nil, err
	}
	if _, err := buf.Write(serializedAddrRecv); err != nil {
		return nil, err
	}

	serializedAddrFrom, err := v.AddrFrom.Serialize()
	if err != nil {
		return nil, err
	}
	if _, err := buf.Write(serializedAddrFrom); err != nil {
		return nil, err
	}

	if err := binary.Write(&buf, binary.BigEndian, v.Nonce); err != nil {
		return nil, err
	}

	serializedUserAgent, err := v.UserAgent.Serialize()
	if err != nil {
		return nil, err
	}
	if _, err := buf.Write(serializedUserAgent); err != nil {
		return nil, err
	}

	if err := binary.Write(&buf, binary.BigEndian, v.StartHeight); err != nil {
		return nil, err
	}

	if err := binary.Write(&buf, binary.BigEndian, v.Relay); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
