package node

import (
	"fmt"
	"testing"
)

func TestNode(t *testing.T) {
	// dependency inject Chain and ApiServer objects into node
	//var chain1 chain.BlockChain
	//chain1.Open("data1")
	//CreateTestBlocks(chain1, 100)

	node1, _ := NewNode("localhost:2000", "simnet", nil)
	go node1.Start()

	node2, _ := NewNode("localhost:2001", "simnet", []string{node1.Address})
	go node2.Start()

	fmt.Println(">>>running ...")
	for {
	}
}
