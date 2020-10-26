package protocol

import (
	"log"
)

const (
	cmdVersion = "version"
	cmdVerack  = "verack"
	cmdPing    = "ping"
	cmdPong    = "pong"
	cmdInv     = "inv"
	cmdTx      = "tx"
	cmdGetData = "getdata"

	commandLength = 12
)

var commands = map[string][commandLength]byte{
	cmdVersion: newCommand(cmdVersion),
	cmdVerack:  newCommand(cmdVerack),
	cmdPing:    newCommand(cmdPing),
	cmdPong:    newCommand(cmdPong),
	cmdInv:     newCommand(cmdInv),
	cmdTx:      newCommand(cmdTx),
	cmdGetData: newCommand(cmdGetData),
}

func newCommand(command string) [commandLength]byte {
	if len(command) > commandLength {
		log.Panicf("command %s is too long", command)
	}

	var packed [commandLength]byte
	buf := make([]byte, commandLength-1)
	copy(packed[:], append([]byte(command), buf...)[:])

	return packed
}
