package tests

import (
	"time"

	core "github.com/alexeyqian/gochain/core"
	utils "github.com/alexeyqian/gochain/utils"
)

func TstCreateAccount(name string) core.Transactioner {
	var tx core.CreateAccountTransaction
	tx.AccountId = utils.CreateUuid()
	tx.AccountName = name
	tx.CreatedOn = uint64(time.Now().Unix())
	tx.ExpiredOn = tx.CreatedOn + uint64(1000000)
	return tx
}
