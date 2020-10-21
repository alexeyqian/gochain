package core

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"math/big"
	"reflect"
)

type Transactioner interface {
	Apply() error
	Validate() error
	// QuickValidate() error
}

// GetRawTransaction
// DecodeRawTransaction

const InitWitness = "init"
const InitAmount = 100
const AmountPerBlock = 100
const BlockInterval = 3 // seconds
const BlockZeroId = "00000000-0000-0000-000000000000"
const GenesisTime = 1632830400 //Date and time (GMT): Tuesday, September 28, 2021 12:00:00 PM

const VoteParentTypeArticle = "VoteParentTypeArticle"
const VoteParentTypeComment = "VoteParentTypeComment"
const VoteParentTypeAccount = "VoteParentTypeAccount"

type CreateAccountTransaction struct {
	Id        string // [256]byte hash256
	Signature string // should be [SIGBITS]byte
	CreatedBy string
	CreatedOn uint64
	PublicKey string

	AccountId   string
	AccountName string
}

type TransferCoinTransaction struct {
	Id        string
	Signature string
	CreatedBy string
	CreatedOn uint64
	PublicKey string

	From   string
	To     string
	Amount uint64
}

type CreateArticleTransaction struct {
	Id        string
	Signature string
	CreatedBy string
	CreatedOn uint64
	PublicKey string

	ArticleId string
	Author    string
	Title     string
	Body      string
	Meta      string
}

type CreateCommentTransaction struct {
	Id        string
	Signature string
	CreatedBy string
	CreatedOn uint64
	PublicKey string

	CommentId string
	ParentId  string
	Commentor string
	Body      string
}

type VoteTransaction struct {
	Id        string
	Signature string
	CreatedBy string
	CreatedOn uint64
	PublicKey string

	ParentId   string
	ParentType string
	Direction  int8
	VotePower  uint64
	Voter      string
}

func SerializeTx(tx Transactioner) []byte {
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

func HashTx(tx Transactioner) []byte {

	id := GetTxId(tx)
	signature := GetTxSignature(tx)
	if id != "" || signature != "" {
		log.Panic("should not call hash after setting id and signature")
	}
	pubkey := GetTxPubKey(tx)
	if pubkey == "" {
		log.Panic("should not call hash if pubkey is empty")
	}

	hash := sha256.Sum256(SerializeTx(tx))
	return hash[:]
}

// Sign sign a transaction id/hash
func SignTx(privkey ecdsa.PrivateKey, tx Transactioner) {
	id := HashTx(tx)
	//fmt.Printf("tx hash len %d and value: %v\n", len(id), id)
	SetTxId(tx, string(id))
	r, s, _ := ecdsa.Sign(rand.Reader, &privkey, id)
	temp := append(r.Bytes(), s.Bytes()...)
	SetTxSignature(tx, string(temp))
}

func VerifyTxSignature(tx Transactioner) bool {
	// unpack values of signature which is a pair of numbers
	signature := []byte(GetTxSignature(tx))
	siglen := len(signature)
	r := big.Int{}
	s := big.Int{}
	r.SetBytes(signature[:(siglen / 2)])
	s.SetBytes(signature[(siglen / 2):])

	// unpack values of public key which is a pair of coordinates
	pubkey := []byte(GetTxPubKey(tx))
	keylen := len(pubkey)
	x := big.Int{}
	y := big.Int{}
	x.SetBytes(pubkey[:(keylen / 2)])
	y.SetBytes(pubkey[(keylen / 2):])

	curve := elliptic.P256()
	rawPubkey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}

	return ecdsa.Verify(&rawPubkey, []byte(GetTxId(tx)), &r, &s)
}

func GetTxId(tx Transactioner) string {
	return reflect.ValueOf(tx).Elem().FieldByName("Id").String()
}

func SetTxId(tx Transactioner, id string) {
	reflect.ValueOf(tx).Elem().FieldByName("Id").SetString(id)
}

func GetTxSignature(tx Transactioner) string {
	return reflect.ValueOf(tx).Elem().FieldByName("Signature").String()
}

func SetTxSignature(tx Transactioner, sig string) {
	reflect.ValueOf(tx).Elem().FieldByName("Signature").SetString(sig)
}

func GetTxPubKey(tx Transactioner) string {
	return reflect.ValueOf(tx).Elem().FieldByName("PublicKey").String()
}
