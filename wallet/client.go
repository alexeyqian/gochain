package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"
	"math/big"

	"github.com/alexeyqian/gochain/core"

	"golang.org/x/crypto/ripemd160"
)

const pkVersion = byte(0x00)
const checksumLength = 4

func GetRawTransaction(id string) {

}

func DecodeRawTransaction(data []byte) {

}

// GetBlockByNum, GetBlockById
func GetBlock(num uint64) {

}

func GetBlockHash(num uint64) string {
	return ""
}

// MakeKeyPair create private and public key
func MakeKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pubKey
}

// GetAddress from public key to address
func GetAddress(pubkey []byte) []byte {
	pkHash := HashPubKey(pubkey)
	versionedPayload := append([]byte{pkVersion}, pkHash...)
	checksum := checksum(versionedPayload, checksumLength)
	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)
	return address
}

// ValidateAddress to check if address is good
func ValidateAddress(address string) bool {
	pkHash := Base58Decode([]byte(address))
	actualChecksum := pkHash[len(pkHash)-checksumLength:]
	//version := pkHash[0]
	pkHash = pkHash[1 : len(pkHash)-checksumLength]
	targetChecksum := checksum(append([]byte{pkVersion}, pkHash...), checksumLength)
	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

// AddressToPubKeyHash convert addr to public key hash
func AddressToPubKeyHash(address []byte) []byte {
	versionedAndChecksumedHash := Base58Decode(address)
	pkHash := versionedAndChecksumedHash[1 : len(versionedAndChecksumedHash)-checksumLength]
	return pkHash
}

// HashPubKey use sha256 and ripemd to hash pk twice
func HashPubKey(pubkey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubkey)
	ripemdHasher := ripemd160.New()
	_, err := ripemdHasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	ripemd := ripemdHasher.Sum(nil)
	return ripemd
}

// Sign sign a transaction
func Sign(privkey ecdsa.PrivateKey) {
	// sign a trimmed version of tx, not all fields
	// only sign important fields
	// tx.id and tx.signature is empty at this point
	// tx.pubkey is not empty at this point
	// tx.id is a hash of tx body except id and signature
	// tx.id = tx.Hash()
	// r, s, err := ecdsa.Sign(rand.Reader, &privkey, tx.id)
	// tx.signature := (r.Bytes(), s.Bytes()...)

}

func VerifySignature(tx core.Transactioner) bool {
	// unpack values of signature which is a pair of numbers
	signature := tx.GetSignature()
	siglen := len(signature)
	r := big.Int{}
	s := big.Int{}
	r.SetBytes(signature[:(siglen / 2)])
	s.SetBytes(signature[(siglen / 2):])

	// unpack values of public key which is a pair of coordinates
	pubkey := tx.GetPubKey()
	keylen := len(pubkey)
	x := big.Int{}
	y := big.Int{}
	x.SetBytes(pubkey[:(keylen / 2)])
	y.SetBytes(pubkey[(keylen / 2):])

	rawPubkey := ecdsa.PublicKey{curve, &x, &y}

	return ecdsa.Verify(&rawPubkey, tx.GetId(), &r, &s)
}

func checksum(payload []byte, length int) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])
	return secondSHA[:length]
}
