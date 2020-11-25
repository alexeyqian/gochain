package chain

import (
	"github.com/alexeyqian/gochain/config"
	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/entity"
	"github.com/alexeyqian/gochain/statusdb"
	"github.com/alexeyqian/gochain/utils"
)

func (c *Chain) genesis() {
	// validation
	if config.GenesisTime%config.BlockInterval != 0 {
		panic("genesis time is incorrect.")
	}

	// update global status
	var gpo entity.Gpo
	gpo.ID = statusdb.GpoKey
	gpo.BlockId = config.BlockZeroId
	gpo.BlockNum = 0
	gpo.Witness = config.InitWitness
	gpo.Time = config.GenesisTime
	gpo.Supply = config.InitAmount
	//fmt.Printf("creating gpo: %+v", gpo)
	c.sdb.CreateGpo(&gpo)

	var wso entity.Wso
	wso.MaxBlockSize = config.DefaultBlockSize
	wso.AccountCreationFee = config.DefaultAccountCreationFee
	c.sdb.CreateWso(&wso)

	// update chain database
	var acc entity.Account
	acc.ID = utils.CreateUuid() // TODO: should be public key string
	acc.Name = config.InitWitness
	acc.CreatedOn = config.GenesisTime
	acc.Coin = config.InitAmount
	err := c.sdb.CreateAccount(&acc)
	if err != nil {
		panic(err)
	}

	// update lgr, create a dummy block 0
	b := core.Block{ID: config.BlockZeroId, Num: 0, CreatedOn: config.GenesisTime, Witness: config.InitWitness}
	c.lgr.Append(utils.Serialize(b))
}
