package tests

import (
	"testing"

	"github.com/alexeyqian/gochain/wallet"
)

func TestMakeKeyPairs(t *testing.T) {
	flag := false

	_, pub := wallet.MakeKeyPair()
	//fmt.Printf("private key: %+v\n", priv)
	//fmt.Printf("public key: %+v\n", pub)
	//fmt.Printf("public key len : %d(bytes), %d(bits)\n", len(pub), len(pub)*8)
	addr := wallet.GenerateAddressFromPubKey(pub)
	//fmt.Printf("address len(bytes): %d\n", len(addr))
	//fmt.Printf("address as bytes: %+v\n", addr)
	//fmt.Printf("address as string: %s\n", addr)

	flag = wallet.ValidateAddress(addr)
	if flag == false {
		t.Errorf("validate address failed")
	}

	flag = wallet.ValidateAddressAgainstPubKey(addr, pub)
	if flag == false {
		t.Errorf("reverse verify failed")
	}
}
