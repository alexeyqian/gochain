package tests

import (
	"testing"
	"time"

	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/utils"
	"github.com/alexeyqian/gochain/wallet"
)

func TestSigning(t *testing.T) {
	flag := false
	privkey, pubkey := wallet.MakeKeyPair()

	var tx core.CreateAccountTransaction
	tx.Id = ""
	tx.Signature = ""
	tx.CreatedBy = "init"
	tx.CreatedOn = uint64(time.Now().Unix())
	tx.PublicKey = string(pubkey)

	tx.AccountId = utils.CreateUuid()
	tx.AccountName = "Alice"

	var itx core.Transactioner = tx
	core.SignTx(privkey, &itx)

	flag = core.VerifyTxSignature(&itx)
	if flag == false {
		t.Errorf("cannot verify tx signature.")
	}
}
