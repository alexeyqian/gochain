package protocol

import (
	"log"
)

const (
	cmdPing       = "ping"
	cmdPong       = "pong"
	cmdVersion    = "version"
	commandLength = 12
)

var commands = map[string][commandLength]byte{
	cmdVersion: newCommand(cmdVersion),
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
