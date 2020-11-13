package chain

import (
	"github.com/alexeyqian/gochain/config"
	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/entity"
	"github.com/alexeyqian/gochain/statusdb"
	"github.com/alexeyqian/gochain/utils"
)

func (c *Chain) genesis() {

	// update global status
	var gpo entity.Gpo
	gpo.ID = statusdb.GpoKey
	gpo.BlockId = core.BlockZeroId
	gpo.BlockNum = 0
	gpo.Witness = core.InitWitness
	gpo.Time = core.GenesisTime
	gpo.Supply = core.InitAmount
	//fmt.Printf("creating gpo: %+v", gpo)
	c.sdb.CreateGpo(&gpo)

	var wso entity.Wso
	wso.MaxBlockSize = config.DefaultBlockSize
	wso.AccountCreationFee = config.DefaultAccountCreationFee
	c.sdb.CreateWso(&wso)

	// update chain database
	var acc entity.Account
	acc.ID = utils.CreateUuid() // TODO: should be public key string
	acc.Name = core.InitWitness
	acc.CreatedOn = core.GenesisTime
	acc.Coin = core.InitAmount
	err := c.sdb.CreateAccount(&acc)
	if err != nil {
		panic(err)
	}

	// update lgr, create a dummy block 0
	b := core.Block{ID: core.BlockZeroId, Num: 0, CreatedOn: core.GenesisTime, Witness: core.InitWitness}
	c.lgr.Append(utils.Serialize(b))
}
