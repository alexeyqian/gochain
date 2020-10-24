package tests

import (
	"fmt"
	"testing"

	"github.com/alexeyqian/gochain/node"
)

func TestNode(t *testing.T) {
	// dependency inject Chain and ApiServer objects into node
	//var chain1 chain.BlockChain
	//chain1.Open("data1")
	//CreateTestBlocks(chain1, 100)

	var node1 node.Node
	node1.Address = "localhost:2000"
	//node1.Chain = chain
	var node2 node.Node
	node2.Address = "localhost:2001"
	node2.KnownNodes = []string{node1.Address} // in app, this should come from config.ini

	go node1.Start()
	go node2.Start()

	// sleep for a while

	fmt.Println(">>>done")
	for {
	}
}
