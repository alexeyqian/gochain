package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"golang.org/x/crypto/ripemd160"
)

const pkVersion = byte(0x00)
const checksumLength = 4

// return a pirvate key, and a 64 bytes (512 bits) long public key raw []byte
func MakeKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pubKey
}

// HashPubKey use sha256 and ripemd to hash pk twice
// from 64 bytes public key to 20 bytes ripemd hash
func hashPubKeyToRipemd(pubkey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubkey)
	ripemdHasher := ripemd160.New()
	_, err := ripemdHasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	ripemd := ripemdHasher.Sum(nil)
	return ripemd
}

func getPubKeyHashFromAddress(address []byte) []byte {
	pkHash := Base58Decode(address)
	return pkHash[1 : len(pkHash)-4]
}

func checksum(payload []byte, length int) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])
	return secondSHA[:length]
}

// create a 34 bytes address from 64 bytes public key
// by hash a public key from 512 bits to 256 bits sha256,
// then hash it again to 160 bits ripemd hash
// then add version (1 bytes), generate versioned hash with checksum (4 bytes)
// so the total len of pre encode address is: 1(version) + 20(ripemd) + 4(checksum) = 25 bytes
// then the address is encoded by base58, which increases the address to 34 bytes
// examples: displayed in base58 encoding
//1Lhqun1E9zZZhodiTqxfPQBcwr1CVDV2sy
//1DNp9T85JZ3fxJvnKgfKB2wovVkCDc4pDH
func GenerateAddressFromPubKey(pubkey []byte) []byte {
	pkHash := hashPubKeyToRipemd(pubkey)
	versionedPayload := append([]byte{pkVersion}, pkHash...)
	checksum := checksum(versionedPayload, checksumLength)
	fullPayload := append(versionedPayload, checksum...)
	//fmt.Printf("full payload len: %d\n", len(fullPayload)) // should be 25 bytes
	address := Base58Encode(fullPayload)
	//fmt.Printf("address len: %d\n", len(address)) // TODO: 34 bytes??
	return address
}

// ValidateAddress to check if address is good
func ValidateAddress(address []byte) bool {
	pkHash := Base58Decode(address)
	actualChecksum := pkHash[len(pkHash)-checksumLength:]
	//version := pkHash[0]
	pkHash = pkHash[1 : len(pkHash)-checksumLength]
	targetChecksum := checksum(append([]byte{pkVersion}, pkHash...), checksumLength)
	return bytes.Equal(actualChecksum, targetChecksum)
}

func ValidateAddressAgainstPubKey(address []byte, pubkey []byte) bool {
	// reverse verifying
	pkHash := getPubKeyHashFromAddress(address)
	origPkHash := hashPubKeyToRipemd(pubkey)
	return bytes.Equal(pkHash, origPkHash)
}

/*

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

*/
