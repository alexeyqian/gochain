package protocol

import (
	"fmt"
	"io"
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

func (v *VarStr) UnmarshalBinary(r io.Reader) error {
	lengthBuf := make([]byte, 1)
	if _, err := r.Read(lengthBuf); err != nil {
		return fmt.Errorf("varstr unmarshalBinary: %+v", err)
	}

	v.Length = uint8(lengthBuf[0])

	stringBuf := make([]byte, v.Length)
	if _, err := r.Read(stringBuf); err != nil {
		return fmt.Errorf("VarStr.UnmarshalBinary: %+v", err)
	}
	v.String = string(stringBuf)

	return nil
}
