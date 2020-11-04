package chain

import (
	"testing"

	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/utils"
)

// TestVoteAccount
func TestVoteAccount(t *testing.T) {
	c := SetupTestChain()

	c.AddPendingTx(CreateTestAccount("alice"))
	c.AddPendingTx(CreateTestAccount("bob"))
	c.GenerateBlock()

	var tx core.VoteTransaction
	tx.ID = utils.CreateUuid()
	alice, _ := c.sdb.GetAccountByName("alice")
	tx.ParentId = alice.ID
	tx.ParentType = core.VoteParentTypeAccount
	tx.Direction = 1
	tx.Voter = "bob"
	tx.VotePower = 1000

	c.AddPendingTx(tx)
	c.GenerateBlock()

	vote, _ := c.sdb.GetVote(tx.ID)
	if vote == nil {
		t.Errorf("create vote error")
	}
	// TODO: should not need to re-get account from statusdb
	alice, _ = c.sdb.GetAccount(tx.ParentId)
	if alice.UpVotes != 1 {
		t.Errorf("up vote error, expected: %d actual: %d", 1, alice.UpVotes)
	}

	if alice.VotePower != tx.VotePower {
		t.Errorf("vote power error, expected: %d, actual: %d", tx.VotePower, alice.VotePower)
	}

	// validate reduce bob's vote power

	TearDownTestChain(c)
}

// TestVoteArticle
// TestVoteComment
