package main

import (
	"fmt"
	"time"

	core "github.com/alexeyqian/gochain/core"
	ledger "github.com/alexeyqian/gochain/ledger"
	utils "github.com/alexeyqian/gochain/utils"
)

/*

func newChain() *chain{
	c := chain{id: "guid-xxx", version: "0.0.1"}
	return &c
}

func createAccountTransaction(name string) *transaction{
	d := "{json data}"
	t := transaction{id: "001", txtype: "CREATE_ACCOUNT", data: d}
	return &t
}

func startup(){

}

func chainGenesis(){
	// init global properties
	// create init witness account

}

func gpo() *global_properties{
	gpo := global_properties{}
	return &gpo
}

const BlockInterval int = 3

func constructBlock() *block{
	g := gpo()
	b := block{last_witness: g.witness, no: g.bno + 1}
	b.createdon = g.time + BlockInterval
	// ...
	return &b
}

func findAccountByName(){

}

func findArticleByPermlink(){

}

func getCommentsByPermlink(){

}

func appendIndex(){

}

func appendBlock(){

}
*/
func main() {
	fmt.Printf("starting ...\n")

	priv := utils.GenerateKey()
	fmt.Printf("private key: %v\n", priv)
	fmt.Printf("public key: %v\n", priv.PublicKey)

	ledger.Open("test_data")
	sec := time.Now().Unix()
	b := core.Block{Id: utils.CreateUuid(), Num: 0, CreatedOn: uint64(sec), Witness: "init_miner"}
	ledger.Append(core.SerializeBlock(&b))
	br := ledger.Read(0)
	fmt.Printf("%+v\n", core.UnSerializeBlock(br))
	ledger.Close()
	ledger.Remove()
	fmt.Println("done")
}
