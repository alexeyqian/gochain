package core

import (
	"math/rand"
	"sort"
	"time"

	"github.com/alexeyqian/gochain/config"
	"github.com/alexeyqian/gochain/entity"
	"github.com/alexeyqian/gochain/statusdb"
)

// Selector witness selector type
type Selector struct {
	sdb *statusdb.StatusDB
}

// NewSelector standard new func
func NewSelector(db *statusdb.StatusDB) *Selector {
	return &Selector{
		sdb: db,
	}
}

func (s *Selector) getNextWitness() string {
	headBlockNumber := s.sdb.HeadBlockNumber()
	// only happends at beginning of each round
	if isNewRound(headBlockNumber) {
		updateWitnessSchedule(s.sdb)
	}

	wso, err := s.sdb.GetWso()
	if err != nil {
		panic("cannot get next witness")
	}

	return wso.CurrentWitnesses[headBlockNumber%len(wso.CurrentWitnesses)]
}

func isNewRound(blockNumber int) bool {
	return blockNumber%config.MaxWitnesses == 0
}

func updateWitnessSchedule(sdb *statusdb.StatusDB) {
	// update current witnesses
	topVoted := getTopVotedWitnesses(sdb, config.MaxWitnesses)
	shuffleWitnesses(topVoted)
	wso, _ := sdb.GetWso()
	for _, item := range topVoted {
		wso.CurrentWitnesses = append(wso.CurrentWitnesses, item.Name)
	}

	// update other median data, such as block size, account creation fee etc.
	// TODO: ...
	// sort them by account creation fee to get medium creation fee
	// position = list.size() / 2
	// sort them by max block size to get mideum block size
	// get mid interest rate

	// update wso
	// wso.mid_account_creation_fee =
	// wso.mid_block_size =
	// wso other features update

	sdb.UpdateWso(wso)
}

func getTopVotedWitnesses(sdb *statusdb.StatusDB, max int) []*entity.Witness {
	witnesses := sdb.GetWitnesses()

	// sort by votes desc
	sort.Slice(witnesses, func(i, j int) bool {
		return witnesses[i].Votes > witnesses[j].Votes
	})

	return witnesses[:max]
}

func shuffleWitnesses(witnesses []*entity.Witness) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(witnesses), func(i, j int) { witnesses[i], witnesses[j] = witnesses[j], witnesses[i] })
}
