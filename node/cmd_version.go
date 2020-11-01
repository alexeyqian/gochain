package node

import (
	"fmt"
	"io"
	"net"

	"github.com/alexeyqian/gochain/binary"
	"github.com/alexeyqian/gochain/utils"

	"github.com/alexeyqian/gochain/protocol"
)

// received a version message from peer, handle it here
func (nd *Node) handleVersion(header *protocol.MessageHeader, conn net.Conn) error {
	var version protocol.MsgVersion

	lr := io.LimitReader(conn, int64(header.Length))
	utils.DeserializeWithReader(&version, lr)

	peer := Peer{
		Address:    conn.RemoteAddr(),
		Connection: conn,
		PongCh:     make(chan uint64),
		NodeType:   version.NodeType,
		UserAgent:  version.UserAgent.String,
		Version:    version.Version,
	}

	n.Peers[peer.ID()] = &peer
	// after peer is added, start the monitor indefinity loop
	//go nd.monitorPeer(&peer)

	// TODO: check if it's new peer or existing peer
	fmt.Printf("new peer %s\n", peer)

	// after receiving a version message
	// node send out a version ack message
	varack, err := protocol.NewVeractMsg(n.Network)
	if err != nil {
		return err
	}

	msg, err := binary.Marshal(verack)
	if err != nil {
		return err
	}

	if _, err := conn.Write(msg); err != nil {
		return err
	}

	return nil
}
