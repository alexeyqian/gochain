package node

import (
	"fmt"
	"io"
	"net"

	"github.com/alexeyqian/gochain/protocol"
	"github.com/alexeyqian/gochain/utils"
)

// received a version message from peer, handle it here
func (nd *Node) handleVersion(header *protocol.MessageHeader, conn net.Conn) error {
	var version protocol.MsgVersion

	lr := io.LimitReader(conn, int64(header.Length))
	err := utils.DeserializeWithReader(&version, lr)
	if err != nil {
		return err
	}

	peer := Peer{
		Address:    conn.RemoteAddr(),
		Connection: conn,
		PongCh:     make(chan int),
		NodeType:   version.NodeType,
		Version:    version.Version,
	}

	nd.Peers[peer.ID()] = &peer
	// after peer is added, start the monitor indefinity loop
	//go nd.monitorPeer(&peer)

	// TODO: check if it's new peer or existing peer
	fmt.Printf("new peer %s\n", peer)

	// after receiving a version message
	// node send out a version ack message
	verack, err := protocol.NewVerackMsg(nd.Network)
	if err != nil {
		return err
	}

	msg := utils.Serialize(verack)

	if _, err := conn.Write(msg); err != nil {
		return err
	}

	return nil
}
