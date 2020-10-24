package node

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
)

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
