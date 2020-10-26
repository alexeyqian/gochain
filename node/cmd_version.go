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

// As soon as a peer is added, a peer liveliness monitor should start running. Let’s define how it should work:
// 1. The monitor triggers once in a while and sends a ‘ping’ message to the peer.
// 2. It waits for a ‘pong’ message containing the nonce from the ‘ping’ message.
// 3. If no ‘pong’ message is received in a certain time span, then the peer is considered dead and is removed from the list.

func (nd Node) monitorPeers(){
	peerPings := make(map[uint64]string)

	for{
		select{
		case nonce := <- n.PongCh:
			peerID := peerPings[nonce]
			if peerID == ""{ // make sure peer is still in the list
				break
			}

			peer := nd.Peers[peerID]
			if peer == nil {
				break
			}

			peer.PongCh <- nonce
			// after directing the nonce, it should be removed to avoid memory leak
			delete(peerPings, nonce)
		case pp := <- nd.PingCh:
			peerPings[pp.nonce] = pp.peerID
		}
	}
}

// sends ping messages, waits for replies and handles 'no replay' case
func (nd *Node) monitorPeer(peer *Peer){
	for{
		time.Sleep(pingIntervalSec * time.Second)
		ping, nonce, err := protocol.NewPingMsg(nd.Network)
		msg, err := binary.Marshal(ping)

		if _, err := peer.Connection.Write(msg); err != nil{
			nd.disconnectPeer(peer.ID())
		}
	}
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