package core

import (
	"errors"

	"github.com/alexeyqian/gochain/entity"
	"github.com/alexeyqian/gochain/statusdb"
)

func (tx *TransferCoinTransaction) Validate(sdb *statusdb.StatusDB) error {
	var err error
	var fromAcc *entity.Account

	fromAcc, err = sdb.GetAccountByName(tx.From)
	if err != nil {
		return errors.New("transfer coin: from account is not exist")
	} else if fromAcc.Coin < tx.Amount {
		return errors.New("transfer coin: no enough coin")
	}

	_, err = sdb.GetAccountByName(tx.To)
	if err != nil {
		return errors.New("transfer coin: to account is not exist")
	}

	return nil
}

func (tx *TransferCoinTransaction) Apply(sdb *statusdb.StatusDB) error {
	err := tx.Validate(sdb)
	if err != nil {
		return err
	}

	fromAcc, _ := sdb.GetAccountByName(tx.From)
	toAcc, _ := sdb.GetAccountByName(tx.To)
	fromAcc.Coin -= tx.Amount
	toAcc.Coin += tx.Amount

	// TODO: check errors
	sdb.UpdateAccount(fromAcc)
	sdb.UpdateAccount(toAcc)

	return nil
}
