package core

import (
	"errors"

	"github.com/alexeyqian/gochain/entity"
	"github.com/alexeyqian/gochain/statusdb"
)

func (tx *CreateAccountTransaction) Validate(sdb *statusdb.StatusDB) error {
	if tx.AccountName == "" {
		return errors.New("name cannot be empty.")
	}
	return nil
}

func (tx *CreateAccountTransaction) Apply(sdb *statusdb.StatusDB) error {
	err := tx.Validate(sdb)
	if err != nil {
		return err
	}

	var acc entity.Account
	acc.ID = tx.AccountId
	acc.Name = tx.AccountName
	sdb.AddAccount(&acc)

	return nil
}
