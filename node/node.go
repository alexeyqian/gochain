package node

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/alexeyqian/gochain/protocol"
)

type Node struct {
	Address    string
	ChainNet   string
	KnownNodes []string
	//Peers      map[string]*Peer
	//PingCh     chan peerPing
	//PongCh     chan uint64
	//mempool    *Mempool
	// pointers to ledger and statusdb
}

func NewNode(addr, chainnet string, knownNodes []string) (*Node, error) {
	ok := protocol.ValidateNetwork(chainnet)
	if !ok {
		return nil, fmt.Errorf("unsupported network %s", chainnet)
	}

	return &Node{
		Address:    addr,
		ChainNet:   chainnet,
		KnownNodes: knownNodes,
		//Peers:      make(map[string]*Peer),
		//PingCh:     make(chan peerPing),
		//pongCh:     make(chan uint64),
	}, nil
}

func (nd *Node) Start() {
	fmt.Printf("start node on: %s\n", nd.Address)

	ln, err := net.Listen("tcp", nd.Address)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	// chain.Open()
	//nd.broadcaseMyVersion()

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
	command := protocol.ExtractCommand(request)
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
	//case commandVersion:
	//nd.handleVersionRequest(request)
	default:
		fmt.Println("Unknown command!")
	}

	conn.Close()
}
