package node

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"

	"github.com/alexeyqian/gochain/protocol"
	"github.com/alexeyqian/gochain/utils"

	"github.com/alexeyqian/gochain/binary"
)

const nodeVersion = 1

// before adding a peer, we must first get basic information about it.
// finish a version handshake
func (nd Node) Run(nodeAddr string) error {
	peerAddr, err := ParseNodeAddr(nodeAddr)
	if err != nil {
		return error
	}

	version, err := protocol.NewVersionMsg(
		nd.Network,
		peerAddr.IP,
		peerAddr.Port,
	)

	if err != nil {
		return err
	}

	msgSerialized, err := binary.Marshal(version)
	if err != nil {
		return err
	}

	conn, err := net.Dial("tcp", nodeAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write(msgSerialized)
	if err != nil {
		return err
	}

	go nd.monitorPeers()
	//go nd.mempool.Run() // TODO: uncomment

	tmp := make([]byte, protocol.MsgHeaderLength)

loop:
	for {
		n, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				return err
			}
			break loop
		}

		var msgHeader protocol.MessageHeader
		if err := utils.DeserializeWithReader(&msgHeader, bytes.NewReader(tmp[:n])); err != nil {
			fmt.Errorf("invalide header: %+v", err)
			continue
		}

		if err := msgHeader.Validate(); err != nil {
			fmt.Errorf(err)
			continue
		}

		fmt.Printf("received message: %s\n", msgHeader.Command)

		switch msgHeader.CommandString() {
		case "version":
			if err := nd.handleVersion(&msgHeader, conn); err != nil {
				fmt.Errorf("failed to handle 'version': %+v", err)
				continue
			}
		case "verack":
			if err := nd.handleVerack(&msgHeader, conn); err != nil {
				fmt.Errorf("failed to handler 'verack': %+v", err)
				continue
			}
		case "ping":
			if err := nd.handlePing(&msgHeader, conn); err != nil {
				fmt.Errorf("failed to handle 'ping': %+v", err)
				continue
			}
		case "pong":
			if err := nd.handlePong(&msgHeader, conn); err != nil {
				fmt.Errorf("failed to handle 'pong': %+v", err)
				continue
			}
		case "inv":
			//if err := no.handleInv(&msgHeader, conn); err != nil {
			//	fmt.Errorf("failed to handle 'inv': %+v", err)
			//	continue
			//}
		case "tx":
			if err := no.handleTx(&msgHeader, conn); err != nil {
				fmt.Printf("failed to handle 'tx': %s", err)
				continue
			}
		}

	}

	return nil
}

func (nd *Node) broadcaseMyVersion() {
	for _, addr := range nd.KnownNodes {
		nd.sendMyVersion(addr)
	}
}

func (nd *Node) sendMyVersion(toAddress string) {
	//height := chain.GetBaseHeight()
	height := 100
	payload := utils.Serialize(nodeVersionRequest{nodeVersion, height, nd.Address})
	request := append(protocol.CommandToBytes(commandVersion), payload...)
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

func (nd *Node) handleGetBlockHashes(request []byte) {

}

func (nd *Node) sendGetBlockHashesRequest(addr string) {
	payload := utils.Serialize(GetBlockHashesRequest{nd.Address})
	request := append(protocol.CommandToBytes(commandGetBlockHashes), payload...)
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
		log.Panic(err) // TODO: panic or return error?
	}
}

func (nd Node) Mempool() map[string]*protocol.MsgTx {
	m := make(map[string]*protocol.MsgTx)

	for k, v := range nd.mempool.txs {
		m[string(k)] = v
	}

	return m
}

func (nd *Node) disconnectPeer(peerID string) {
	fmt.Printf("disconnectiong peer %s\n", peerID)

	peer := nd.Peers[peerID]
	if peer == nil {
		return
	}

	peer.Connection.Close()
}
