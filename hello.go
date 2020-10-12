package main
import "fmt"
import "time"
import "./ledger"
/*
func createUuid() string{
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	//fmt.Println(uuid)
	return uuid
}

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
func main(){
	fmt.Println("hello world!")
	now := time.Now() 
	sec := now.Unix()
	fmt.Println("sec: %v", sec)
	var lg ledger.Ledger
	//ledger.OpenLedger("data")
	fmt.Println("success")
}

