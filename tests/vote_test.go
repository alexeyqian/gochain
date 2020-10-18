package tests

import (
	"testing"

	"github.com/alexeyqian/gochain/chain"
	"github.com/alexeyqian/gochain/core"
	"github.com/alexeyqian/gochain/statusdb"
	"github.com/alexeyqian/gochain/utils"
)

// TestVoteAccount
func TestVoteAccount(t *testing.T) {
	chain.Open(TestDataDir)

	chain.AddPendingTx(CreateTestAccount("alice"))
	chain.AddPendingTx(CreateTestAccount("bob"))
	chain.GenerateBlock()

	var tx core.VoteTransaction
	tx.Id = utils.CreateUuid()
	alice, _ := statusdb.GetAccountByName("alice")
	tx.ParentId = alice.Id
	tx.ParentType = core.VoteParentTypeAccount
	tx.Direction = 1
	tx.Voter = "bob"
	tx.VotePower = 1000

	chain.AddPendingTx(tx)
	chain.GenerateBlock()

	vote, _ := statusdb.GetVote(tx.Id)
	if vote == nil {
		t.Errorf("create vote error")
	}
	// TODO: should not need to re-get account from statusdb
	alice, _ = statusdb.GetAccount(tx.ParentId)
	if alice.UpVotes != 1 {
		t.Errorf("up vote error, expected: %d actual: %d", 1, alice.UpVotes)
	}

	if alice.VotePower != tx.VotePower {
		t.Errorf("vote power error, expected: %d, actual: %d", tx.VotePower, alice.VotePower)
	}

	// validate reduce bob's vote power

	chain.Close()
	chain.Remove()
}

// TestVoteArticle
// TestVoteComment
