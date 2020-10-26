package protocol

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"strings"
)

const (
	checksumLength  = 4
	nodeNetwork     = 1
	magicLength     = 4
	MsgHeaderLength = magicLength + commandLength + checksumLength + 4 // 4 layload length value
)

var (
	magicMainnet = [magicLength]byte{0xf9, 0xbe, 0xb4, 0xd9}
	magicSimnet  = [magicLength]byte{0x16, 0x1c, 0x14, 0x12}
	networks     = map[string][magicLength]byte{
		"mainnet": magicMainnet,
		"simnet":  magicSimnet,
	}
)

type Magic [magicLength]byte

type MessageHeader struct {
	Magic    [magicLength]byte
	Command  [commandLength]byte
	Length   uint32 // length of payload
	Checksum [checksumLength]byte
}

type Message struct {
	MessageHeader
	Payload []byte
}

func NewMessage(cmd, network string, payload interface{}) (*Message, error) {
	serializedPayload, err := binary.Marshal(payload)
	if err != nil {
		return nil, err
	}

	command, ok := commands[cmd]
	if !ok {
		return nil, fmt.Errorf("unsupported comamnd %s", cmd)
	}

	magic, ok := networks[network]
	if !ok {
		return nil, fmt.Errorf("unsupported network %s", network)
	}

	msg := Message{
		MessageHeader: MessageHeader{
			Magic:    magic,
			Command:  command,
			Length:   uint32(len(serializedPayload)),
			Checksum: checksum(serializedPayload),
		},
		Payload: serializedPayload,
	}

	return &msg, nil
}

func (mh MessageHeader) CommandString() string {
	return strings.Trim(string(mh.Command[:]), string(0))
}

func (mh MessageHeader) Validate() error {
	if !mh.HasValidMagic() {
		return fmt.Errorf("invalid magic: %x", mh.Magic)
	}

	if !mh.HasValidCommand() {
		return fmt.Errorf("invalid command %s", mh.CommandString())
	}

	return nil
}

func (mh MessageHeader) HasValidCommand() bool {
	_, ok := commands[mh.CommandString()]
	return ok
}

func (mh MessageHeader) HasValidMagic() bool {
	switch mh.Magic {
	case magicMainnet, magicSimnet:
		return true
	}

	return false
}

func checksum(data []byte) [checksumLength]byte {
	hash := sha256.Sum256(data)
	hash = sha256.Sum256(hash[:])
	var hashArr [checksumLength]byte
	copy(hashArr[:], hash[0:checksumLength])
	return hashArr
}

// There’re different formats of serialization.
// encoding/gob is very Golang way of serialization
// other language don't support it.
// so the customized serialization mechanic:
// take byte representation of every field and concatenate
// them preserving the order.

func (m Message) Serialize() ([]byte, error) {
	var buf bytes.Buffer

	if _, err := buf.Write(m.Magic[:]); err != nil {
		return nil, err
	}

	if _, err := buf.Write(m.Command[:]); err != nil {
		return nil, err
	}

	if err := binary.Write(&buf, binary.BigEndian, m.Length); err != nil {
		return nil, err
	}

	if _, err := buf.Write(m.Checksum[:]); err != nil {
		return nil, err
	}

	if _, err := buf.Write(m.Payload); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
