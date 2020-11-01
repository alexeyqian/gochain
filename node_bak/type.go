package node

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
