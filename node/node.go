package node

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"time"

	"github.com/alexeyqian/gochain/binary"
)

type Node struct {
	Address    string
	KnownNodes []string
	// Dependency Injection for testing
	//BlockChain *chain.Chain
	//ApiServer  *network.ApiServer
	Peers  map[string]*Peer
	PingCh chan peerPing
	PongCh chan uint64
}

func NewNode(network, userAgent string)(*Node, error){
	networkMagic, ok := protocol.Networks[network]
	if !ok{
		return nil, fmt.Errorf("unsupported network %s", network)
	}

	return &Node{
		Network: network,
		NetworkMagic: networkMagic,
		Peers: make(map[string]*Peer),
		PingCh: make(chan peerPing),
		pongCh: make(chan uint64),
		UserAgent: userAgent,
	}, nil
}

// before adding a peer, we must first get basic information about it.
// finish a version handshake

// TODO: replace Start
func (nd Node) Run(nodeAddr string) error{
	peerAddr, err := ParseNodeAddr(nodeAddr)
	if err != nil{
		return error
	}

	version, err := protocol.NewVersinMsg(
		nd.Network,
		nd.UserAgent,
		peerAddr.IP,
		peerAddr.Port,
	)

	if err != nil{
		return err
	}

	msgSerialized, err := binary.Marshal(version)
	if err != nil{
		return err
	}

	conn, err := net.Dial("tcp", nodeAddr)
	if err != nil{
		return err
	}
	defer conn.Close()

	_, err = conn.Write(msgSerialized)
	if err != nil{
		return err
	}

	go nd.monitorPeers()

	tmp := make([]byte, protocol.MsgHeaderLength)

	loop:
	for{
		n, err := conn.Read(tmp)
		if err != nil{
			if err != io.EOF {
				return err
			}
			break loop
		}

		var msgHeader protocol.MessageHeader
		if err := binary.NewDecoder(bytes.NewReader(tmp[:n])).Decode(&msgHeader); err != nil{
			fmt.Errorf("invalide header: %+v", err)
			continue
		}

		if err := msgHeader.Validate(); err != nil{
			fmt.Errorf(err)
			continue
		}

		fmt.Printf("received message: %s\n", msgHeader.Command)

		switch msgHeader.CommandString(){
		case "version":
			if err := nd.handleVersion(&msgHeader, conn); err != nil{
				fmt.Errorf("failed to handle 'version': %+v", err)
				continue
			}
		case "verack":
			if err := nd.handleVerack(&msgHeader, conn); err != nil{
				fmt.Errorf("failed to handler 'verack': %+v", err)
				continue
			}
		case "ping":
			if err := no.handlePing(&msgHeader, conn); err != nil {
				fmt.Errorf("failed to handle 'ping': %+v", err)
				continue
			}
		case "pong":
			if err := no.handlePong(&msgHeader, conn); err != nil {
				fmt.Errorf("failed to handle 'pong': %+v", err)
				continue
			}
		case "inv":
			if err := no.handleInv(&msgHeader, conn); err != nil {
				fmt.Errorf("failed to handle 'inv': %+v", err)
				continue
			}
		case "tx":
			if err := no.handleTx(&msgHeader, conn); err != nil {
				fmt.Errorf("failed to handle 'tx': %+v", err)
				continue
			}		
		}

	}

	return nil
}

func (nd *Node) Start() {
	fmt.Printf("start node on: %s\n", nd.Address)

	ln, err := net.Listen(protocol, nd.Address)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	// chain.Open()

	nd.broadcaseMyVersion()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}

		go nd.handleConnection(conn)
	}

}

func (nd *Node) handleConnection(conn net.Conn) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	command := bytesToCommand(request[:commandLength])
	fmt.Printf("%s received command: %s from %s\n", nd.Address, command, conn.RemoteAddr())

	switch command {
	/*case "addr":
		handleAddr(request)
	case "block":
		handleBlock(request, bc)
	case "inv":
		handleInv(request, bc)
	case "getblocks":
		handleGetBlocks(request, bc)
	case "getdata":
		handleGetData(request, bc)
	case "tx":
		handleTx(request, bc)*/
	case commandVersion:
		nd.handleVersionRequest(request)
	default:
		fmt.Println("Unknown command!")
	}

	conn.Close()
}

func (nd *Node) broadcaseMyVersion() {
	for _, addr := range nd.KnownNodes {
		nd.sendMyVersion(addr)
	}
}

func (nd *Node) sendMyVersion(toAddress string) {
	//height := chain.GetBaseHeight()
	height := 100
	payload := gobEncode(nodeVersionRequest{nodeVersion, height, nd.Address})
	request := append(commandToBytes(commandVersion), payload...)
	nd.sendData(toAddress, request)
}

func (nd *Node) handleVersionRequest(request []byte) {
	var buff bytes.Buffer
	var payload nodeVersionRequest
	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		fmt.Println("cannot decode payload")
		return
	}

	myHeight := bc.GetBaseHeight()
	foreignHeight := payload.BaseHight

	if myHeight < foreignHeight {
		sendGetBlocks(payload.AddressFrom)
	} else if myHeight > foreignHeight {
		sendVersion(payload.AddressFrom, bc)
	}

	if !nodeIsKnown(payload.AddressFrom) {
		knownNodes = append(knownNodes, payload.AddressFrom)
	}
}

func (nd *Node) handleInv(request []byte) {
	var buff bytes.Buffer
	var payload InvRequest

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		fmt.Println("cannot decode payload")
		return
	}

	// processing ...

}

func (nd *Node) handleGetBlockHashes(request []byte) {

}

func (nd *Node) sendGetBlockHashesRequest(addr string) {
	payload := gobEncode(GetBlockHashesRequest{nd.Address})
	request := append(commandToBytes(commandGetBlockHashes), payload...)
	nd.sendData(addr, request)
}

func nodeIsKnown(addr string, nodes []string) bool {
	for _, node := range nodes {
		if node == addr {
			return true
		}
	}
	return false
}

func (nd *Node) sendData(addr string, data []byte) {
	conn, err := net.Dial(protocol, addr)
	if err != nil {
		fmt.Printf("%s is not available\n", addr)
		var updatedNodes []string
		for _, n := range nd.KnownNodes {
			if n != addr {
				updatedNodes = append(updatedNodes, n)
			}
		}
		nd.KnownNodes = updatedNodes
		return
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}

// As soon as a peer is added, a peer liveliness monitor should start running. Let’s define how it should work:
// 1. The monitor triggers once in a while and sends a ‘ping’ message to the peer.
// 2. It waits for a ‘pong’ message containing the nonce from the ‘ping’ message.
// 3. If no ‘pong’ message is received in a certain time span, then the peer is considered dead and is removed from the list.

func (nd Node) monitorPeers() {
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

// sends ping messages
// waits for replies and handles 'no replay' case
func (nd *Node) monitorPeer(peer *Peer) {
	for {
		time.Sleep(pingIntervalSec * time.Second)
		ping, nonce, err := protocol.NewPingMsg(nd.Network)
		msg, err := binary.Marshal(ping)
		if err != nil{
			fmt.Fatalf("monitor peer, binary marshal: %v", err)
		}

		if _, err := peer.Connection.Write(msg); err != nil {
			nd.disconnectPeer(peer.ID())
		}

		fmt.Debugf("send 'ping' to %s", peer)

		nd.PingCh <- peerPing{
			nonce: nonce, 
			peerID: peer.ID()
		}

		t := time.NewTimer(pingTimeoutSec * time.Second)

		select{
		case pn := <- peer.PongCh:
			if pn != nonce{
				fmt.Errorf("nonce doesn't match for %s, expected %d, got %d", peer, nonce, pn)
				nd.disconnectPeer(peer.ID)
				return
			}
			ftm.Debugf("got 'pong' from %s", peer)
		case <- t.C:
			nd.disconnectPeer(peer.ID())
			return
		}

		// TODO: timer and return sequence
		t.Stop()
	}
}

func (nd Node) disconnectPeer(peerID string){
	fmt.Printf("disconnectiong peer %s\n", peerID)

	peer := nd.Peers[peerID]
	if peer == nil{
		return
	}

	peer.Connection.Close()
}