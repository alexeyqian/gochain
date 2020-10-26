package protocol

import (
	"encoding/hex"
	"testing"
	"time"
)

func TestMessageSerialization(t *testing.T) {
	version := MsgVersion{
		Version:   Version,
		Services:  SrvNodeNetwork,
		Timestamp: time.Date(2019, 11, 11, 0, 0, 0, 0, time.UTC).Unix(),
		AddrRecv: NetAddr{
			Services: SrvNodeNetwork,
			IP:       NewIPv4(127, 0, 0, 1),
			Port:     9333,
		},
		AddrFrom: NetAddr{
			Services: SrvNodeNetwork,
			IP:       NewIPv4(127, 0, 0, 1),
			Port:     9334,
		},
		Nonce:       31337,
		UserAgent:   NewUserAgent(),
		StartHeight: -1,
		Relay:       true,
	}
	msg, err := NewMessage("version", "simnet", version)
	if err != nil {
		t.Errorf("unexpected error: %+v", err)
		return
	}

	msgSerialized, err := msg.Serialize()
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
