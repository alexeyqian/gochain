package utils

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/gob"
	"encoding/pem"
	"fmt"
	"io"
	"log"
)

func CreateUuid() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

func GenerateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func FindString(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func Serialize(e interface{}) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(e)
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func Deserialize(e interface{}, data []byte) {
	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(e)
	if err != nil {
		panic(err)
	}
}

func DeserializeWithReader(data interface{}, reader io.Reader) error {
	dec := gob.NewDecoder(reader)
	return dec.Decode(data)
}

func EncodePubKeyInPem(publicKey *ecdsa.PublicKey) string {
	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

	return string(pemEncodedPub)
}

/*
func EncodeInPem(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) (string, string) {
	x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

	return string(pemEncoded), string(pemEncodedPub)
}*/

func DecodeInPem(pemEncoded string, pemEncodedPub string) (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	block, _ := pem.Decode([]byte(pemEncoded))
	x509Encoded := block.Bytes
	privateKey, _ := x509.ParseECPrivateKey(x509Encoded)

	blockPub, _ := pem.Decode([]byte(pemEncodedPub))
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
	publicKey := genericPublicKey.(*ecdsa.PublicKey)

	return privateKey, publicKey
}

/*
func test() {
    privateKey, _ := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
    publicKey := &privateKey.PublicKey

    encPriv, encPub := encode(privateKey, publicKey)

    fmt.Println(encPriv)
    fmt.Println(encPub)

    priv2, pub2 := decode(encPriv, encPub)

    if !reflect.DeepEqual(privateKey, priv2) {
        fmt.Println("Private keys do not match.")
    }
    if !reflect.DeepEqual(publicKey, pub2) {
        fmt.Println("Public keys do not match.")
    }
}*/
