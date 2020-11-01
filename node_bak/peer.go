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
	Version    int32
}

func (p *Peer) ID() string {
	return p.Address.String()
}

func (p *Peer) String() string {
	return fmt.Sprintf("%s", p.Address)
}

type peerPing struct {
	nonce  uint64
	peerID string
}

// As soon as a peer is added, a peer liveliness monitor should start running. Let’s define how it should work:
// 1. The monitor triggers once in a while and sends a ‘ping’ message to the peer.
// 2. It waits for a ‘pong’ message containing the nonce from the ‘ping’ message.
// 3. If no ‘pong’ message is received in a certain time span, then the peer is considered dead and is removed from the list.

func (nd *Node) monitorPeers() {
	// TODO: should we use peerID as key?
	// since nonc might be same from different peers
	peerPings := make(map[uint64]string)

	for {
		select {
		case pp := <-nd.PingCh:
			peerPings[pp.nonce] = pp.peerID

		// pass pong messages from the handler
		case nonce := <-nd.PongCh:
			peerID := peerPings[nonce]
			if peerID == "" { // make sure peer is still in the list
				break
			}

			peer := nd.Peers[peerID]
			if peer == nil {
				break
			}

			peer.PongCh <- nonce
			// after directing the nonce, it should be removed to avoid memory leak
			delete(peerPings, nonce)
		}
	}
}

// sends ping messages
// waits for replies and handles 'no replay' case
/*
func (nd *Node) monitorPeer(peer *Peer) {
	for {
		time.Sleep(pingIntervalSec * time.Second)
		ping, nonce, err := protocol.NewPingMsg(nd.Network)
		msg, err := binary.Marshal(ping)
		if err != nil{
			panic(err)
		}

		if _, err := peer.Connection.Write(msg); err != nil {
			nd.disconnectPeer(peer.ID())
		}

		fmt.Printf("send 'ping' to %s", peer)

		nd.PingCh <- peerPing{
			nonce: nonce,
			peerID: peer.ID(),
		}

		t := time.NewTimer(pingTimeoutSec * time.Second)

		select{
		case pn := <- peer.PongCh:
			if pn != nonce{
				fmt.Printf("nonce doesn't match for %s, expected %d, got %d", peer, nonce, pn)
				nd.disconnectPeer(peer.ID)
				return
			}
			ftm.Printf("got 'pong' from %s", peer)
		case <- t.C:
			nd.disconnectPeer(peer.ID())
			return
		}

		// TODO: timer and return sequence
		t.Stop()
	}
}
*/
