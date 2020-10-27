package protocol

/*
Transactions and blocks transferring happens this way:

1. When a node gets a new transaction, it sends ‘inv’ message to its peers. ‘inv’ means inventory and it literally says “Hey! I have these…”. But ‘inv’ doesn’t contain full data, only hashes.
2. Any peer that receives the message can decide where it wants to get full data or not.
3. If a peer wants full data, it sends a ‘getdata’ reply specifying a list of hashes it want to get full data for.
4. A node that receives ‘getdata’ checks what objects were requested (transactions or blocks) and sends them in related messages: ‘tx’ for transaction and ‘block’ for block (one transaction/block per message).
*/
type MsgTx struct {
	Version int32
	Witness string
}

// func NewMsgTx
