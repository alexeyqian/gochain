package node

type peerPing struct {
	nonce  int
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
