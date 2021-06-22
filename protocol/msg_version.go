package protocol

import (
	"math/rand"
	"time"
)

type VersionNetAddr struct {
	Time     int
	NodeType int
	IP       IPv4
	Port     uint16
}

type MsgVersion struct {
	Version     int32
	NodeType    int
	Timestamp   int64
	AddrRecv    VersionNetAddr
	AddrFrom    VersionNetAddr
	Nonce       int
	StartHeight int32
}

func NewVersionMsg(network, userAgent string, peerIP IPv4, peerPort uint16) (*Message, error) {
	payload := MsgVersion{
		Version:   Version,
		NodeType:  NodeTypeFull,
		Timestamp: time.Now().UTC().Unix(),
		AddrRecv: VersionNetAddr{
			NodeType: NodeTypeFull,
			IP:       peerIP,
			Port:     peerPort,
		},
		AddrFrom: VersionNetAddr{
			NodeType: NodeTypeFull,
			IP:       NewIPv4(127, 0, 0, 1),
			Port:     9334,
		},
		Nonce:       rand.Uint64(),
		StartHeight: -1,
	}

	msg, err := NewMessage(network, cmdVersion, payload)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
