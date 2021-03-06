package network

import (
	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/protocol"
	"github.com/alexeyqian/gochain/utils"
)

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
	payload := utils.Serialize(data)
	request := append(protocol.CommandToBytes("block"), payload...)
	sendData(addr, request)
}

func handleGetBlockHashes(request []byte, bc *BlockChain) {
	// ...
	blocks := bc.GetBlockHashes()
	sendBlockHashes(payload.AddressFrom, "block", blocks)
}
