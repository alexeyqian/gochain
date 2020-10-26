package node

import (
	"io"
	"net"
)

// system will handle multiple connections
// one connection per peer
type Peer struct {
	Address    net.Addr
	Connection io.ReadWriteCloser
	PongCh     chan uint64 // pass pong replies
	Services   uint64      // describes feature supported by the peer
	UserAgent  string
	Version    int32
}

func (p Peer) ID() string {
	return p.Address.String()
}
