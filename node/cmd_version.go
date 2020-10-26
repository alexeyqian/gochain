package node

import (
	"github.com/alexeyqian/gochain/protocol"
	"net"
)

// before adding a peer, we must first get basic information about it.
// finish a version handshake

func (nd Node) handleVersion(header *protocol.MessageHeader, conn net.Conn) error {
	var version protocol.MsgVersion

	lr := io.LimitReader(conn, int64(header.Length))
	if err := binary.NewDecoder(lr).Decode(&version); err != nil{
		return err
	}

	peer := Peer{
		Address: conn.RemoteAddr()
		Connection: conn,
		PonCh: make(chan uint64),
		Services: version.Services,
		UserAgent: version.UserAgent.String,
		Version: version.Version
	}

	n.Peers[peer.ID()] = &peer
	// after peer is added, start the monitor cycle
	go nd.monitorPeer(&peer)
	 // ...
}

func (nd Node) handlePong(header *protocol.MessageHeader, conn io.ReadWriter) error{
	var pong protocol.MsgPing // ?? MsgPong

	lr := io.LimitReader(conn, int64(header.Length))
	if err := binary.NewDecoder(lr).Decode(&pong); err != nil{
		return err
	}

	nd.PongCh <- pong.Nonce
	return nil
}