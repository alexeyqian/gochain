package protocol

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/alexeyqian/gochain/utils"
)

func TestMessageSerialization(t *testing.T) {
	version := MsgVersion{
		Version:   Version,
		NodeType:  NodeTypeFull,
		Timestamp: time.Date(2019, 11, 11, 0, 0, 0, 0, time.UTC).Unix(),
		AddrRecv: VersionNetAddr{
			NodeType: NodeTypeFull,
			IP:       NewIPv4(127, 0, 0, 1),
			Port:     9333,
		},
		AddrFrom: VersionNetAddr{
			NodeType: NodeTypeFull,
			IP:       NewIPv4(127, 0, 0, 1),
			Port:     9334,
		},
		Nonce:       31337,
		UserAgent:   NewVarStr(UserAgent),
		StartHeight: -1,
	}
	msg, err := NewMessage("simnet", "version", version)
	if err != nil {
		t.Errorf("unexpected error: %+v", err)
		return
	}

	msgSerialized, err := utils.SerializeStruct(&msg)
	if err != nil {
		t.Errorf("unexpected error: %+v", err)
		return
	}

	actual := hex.EncodeToString(msgSerialized)
	expected := "161c141276657273696f6e000000000000000072fd02fade0001117f0000000000000001000000005dc8a480000000000000000100000000000000000000ffff7f0000012475000000000000000100000000000000000000ffff7f00000124760000000000007a691c2f5361746f7368693a352e36342f74696e796269743a302e302e312fffffffff01"
	if actual != expected {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}

}
