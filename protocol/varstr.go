package protocol

import (
	"bytes"
	"encoding/binary"
)

// to serialize a message, we need to knwo legnths of all fields.
// Since strings aren't fixed, ew also need to store length of each string
type VarStr struct {
	Length uint8
	String string
}

func newVarStr(str string) VarStr {
	return VarStr{
		Length: uint8(len(str)),
		String: str,
	}
}

func (v VarStr) Serialize() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.BigEndian, v.Length); err != nil {
		return nil, err
	}

	if _, err := buf.Write([]byte(v.String)); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
