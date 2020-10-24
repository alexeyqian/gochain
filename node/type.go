package node

type Node struct {
	Address    string
	KnownNodes []string
	// Dependency Injection for testing
	//BlockChain *chain.Chain
	//ApiServer  *network.ApiServer
}

type nodeVersionRequest struct {
	Version int
	Height  int
	Address string
}

type blockHashesRequest struct {
	From string
}
