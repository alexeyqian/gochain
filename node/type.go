package node

type Node struct {
	Address    string
	KnownNodes []string
	// Dependency Injection for testing
	//BlockChain *chain.Chain
	//ApiServer  *network.ApiServer
	Peers  map[string]*Peer
	PingCh chan peerPing
	PongCh chan uint64
}

type peerPing struct {
	nonce  uint64
	peerID string
}

type nodeVersionRequest struct {
	Version int
	Height  int
	Address string
}

type blockHashesRequest struct {
	From string
}
