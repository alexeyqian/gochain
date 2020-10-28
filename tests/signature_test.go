package tests

import (
	"testing"
	"time"

	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/utils"
	"github.com/alexeyqian/gochain/wallet"
)

func TestGetAndSetTx(t *testing.T) {
	_, pubkey := wallet.MakeKeyPair()

	var tx core.CreateAccountTransaction
	tx.ID = "hello"
	tx.Signature = "world"
	tx.PublicKey = string(pubkey)

	id := core.GetTxId(&tx)
	sig := core.GetTxSignature(&tx)
	pub := core.GetTxPubKey(&tx)

	if id != "hello" || sig != "world" || pub != string(pubkey) {
		t.Errorf("cannot get tx data by using reflection.")
	}

	core.SetTxId(&tx, "updated_id")
	core.SetTxSignature(&tx, "updated_sig")
	if tx.ID != "updated_id" || tx.Signature != "updated_sig" {
		t.Errorf("cannot set tx data by using reflection.")
	}
	//fmt.Println("set tx id and sig success!")
}

func TestSigning(t *testing.T) {
	flag := false
	privkey, pubkey := wallet.MakeKeyPair()

	var tx core.CreateAccountTransaction
	tx.ID = ""
	tx.Signature = ""
	tx.CreatedBy = "init"
	tx.CreatedOn = uint64(time.Now().Unix())
	tx.PublicKey = string(pubkey)

	tx.AccountId = utils.CreateUuid()
	tx.AccountName = "Alice"

	core.SignTx(privkey, &tx)
	//fmt.Printf("signature: %v\n", tx.Signature)
	flag = core.VerifyTxSignature(&tx)
	if flag == false {
		t.Errorf("cannot verify tx signature.")
	}
}
