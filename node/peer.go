package node

import (
	"fmt"
	"io"
	"net"
)

const (
	pingIntervalSec = 120
	pingTimeoutSec  = 30
)

// system will handle multiple connections
// one connection per peer
type Peer struct {
	Address    net.Addr
	Connection io.ReadWriteCloser
	PongCh     chan uint64 // pass pong replies
	NodeType   uint64      // describes feature supported by the peer
	UserAgent  string
	Version    int32
}

func (p Peer) ID() string {
	return p.Address.String()
}

func (p Peer) String() string {
	return fmt.Sprintf("%s (%s)", p.UserAgent, p.Address)
}

type peerPing struct {
	nonce  uint64
	peerID string
}
