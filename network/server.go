package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/alexeyqian/gochain/core"
)

type NodeVersion struct {
	Version    int
	BaseHeight int
	Address    string
}

const protocol = "tcp"
const nodeVersion = 1
const commandLength = 20

var nodeAddress string
var knowNodes = []string{"localhost:3000"}
var blocksInTransit = [][]byte{}

func StartServer(port int, address string) {
	nodeAddress := fmt.Sprintf("localhost:%s", port)

	ln, err := net.Listen(protocol, nodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	bc := NewBlockChain(port)

	knownNodes := GetKnownNodes()
	for _, node := range knownNodes {
		sendVersion(node, bc)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConnection(conn, bc)
	}
}

func GetKnownNodes() []string {
	var nodes []string
	nodes = append(nodes, "localhost:2000")
	//nodes = append(nodes, "localhost:2001")
	return nodes
}

func sendVersion(address string) {
	height := chain.GetBaseHeight()
	payload := gobEncode(version{version, height, myAddress})

	request := append(commandToBytes("version"), payload...)
	sendData(address, request)
}

// convert string to bytes data, leave the rest bytes empty
func commandToBytes(command string) []byte {
	var data [commandLength]byte

	for i, c := range command {
		data[i] = byte(c)
	}

	return data[:]
}

// extract command string from bytes data
func bytesToCommand(data []byte) string {
	var command []byte

	for _, b := range data {
		if b != 0x0 {
			command = append(command, b)
		}
	}
}

func extractCommand(request []byte) []byte {
	return request[:commandLength]
}

func requestBlocks() {
	for _, node := range knowNodes {
		sendGetBlocks(node)
	}
}

/*
func sendAddresses(address string){
	nodes :=
}*/

type SendingBlock struct {
	Address   string
	BlockData []byte
}

func sendBlock(addr string, b *core.Block) {
	data := SendingBlock{Address: nodeAddress, BlockData: b.Serialize()}
	payload := gobEncode(data)
	request := append(commandToBytes("block"), payload...)
	sendData(addr, request)
}

func handleConnection(conn net.Conn, bc *BlockChain) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	command := bytesToCommand(request[:commandLength])
	fmt.Printf("received command: %s\n", command)

	switch command {
	case "addr":
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
		handleTx(request, bc)
	case "version":
		handleVersion(request, bc)
	default:
		fmt.Println("Unknown command!")
	}

	conn.Close()
}

func handleVersion(request []byte, bc *BlockChain) {
	var buff bytesToCommand
	var payload Version
	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)

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

func handleGetBlockHashes(request []byte, bc *BlockChain) {
	// ...
	blocks := bc.GetBlockHashes()
	sendBlockHashes(payload.AddressFrom, "block", blocks)
}

func gobEncode(data interface{}) []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func nodeIsKnown(addr string, nodes []string) bool {
	for _, node := range nodes {
		if node == addr {
			return true
		}
	}
	return false
}
