package binary

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

const (
	commandLength          = 12
	magicAndChecksumLength = 4
)

type Marshaler interface {
	MarshalBinary() ([]byte, error)
}

func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer

	if reflect.TypeOf(v).Kind() == reflect.Ptr {
		v = reflect.ValueOf(v).Elem().Interface()
	}

	switch val := v.(type) {
	case uint8, int8, uint16, int16, uint32, int32, uint64, int64, bool:
		if err := binary.Write(&buf, binary.BigEndian, val); err != nil {
			return nil, err
		}
	case [magicAndChecksumLength]byte:
		if _, err := buf.Write(val[:]); err != nil {
			return nil, err
		}
	case [commandLength]byte:
		if _, err := buf.Write(val[:]); err != nil {
			return nil, err
		}
	case []byte:
		if _, err := buf.Write(val); err != nil {
			return nil, err
		}
	case string:
		if _, err := buf.Write([]byte(val)); err != nil {
			return nil, err
		}
	case Marshaler:
		b, err := val.MarshalBinary()
		if err != nil {
			return nil, err
		}

		if _, err := buf.Write(b); err != nil {
			return nil, err
		}
	default:
		// is it a struct
		if reflect.ValueOf(v).Kind() == reflect.Struct {
			b, err := marshalStruct(v)
			if err != nil {
				return nil, err
			}

			if _, err := buf.Write(b); err != nil {
				return nil, err
			}
			break
		}

		return nil, fmt.Errorf("unsupported type %s", reflect.TypeOf(v).String())
	}

	return buf.Bytes(), nil
}

// iterate over all fields and serialize each of them seperately
// and concatenates the results
func marshalStruct(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	vv := reflect.ValueOf(v)

	for i := 0; i < vv.NumField(); i++ {
		s, err := Marshal(reflect.Indirect(vv.Field(i)).Interface)
		if err != nil {
			f := reflect.TypeOf(v).Field(i).Name
			return nil, fmt.Errorf("failed to marshal field %s: %v", f, err)
		}

		if _, err := buf.Write(s); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}
