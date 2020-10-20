package core

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"reflect"
)

type Transactioner interface {
	Apply() error
	Validate() error
	//FastValidate() error // used to validate received tx from network, and called by Validate
	Hash()
	// TrimmedCopy: get trimmed copy for signing
	// GetId()
	// GetPubKey()
	// GetSignature()
	// // Sign()
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

func SerializeTx(tx *Transactioner) []byte {
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(*tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

func HashTx(tx *Transactioner) []byte {
	//fv := reflect.ValueOf(tx).Elem().FieldByName("Id")
	//id := fv.String()
	//fv = reflect.ValueOf(tx).Elem().FieldByName("Signature")
	//signature := fv.String()
	//if id != nil || signature != nil {
	//	log.Panic("should not call hash after setting id and signature")
	//}
	// assert public key is not mepty

	hash := sha256.Sum256(SerializeTx(tx))
	return hash[:]
}

// Sign sign a transaction id/hash
func Sign(privkey ecdsa.PrivateKey, tx *Transactioner) {
	id := HashTx(tx)
	SetTxId(tx, string(id))
	r, s, _ := ecdsa.Sign(rand.Reader, &privkey, id)
	temp := append(r.Bytes(), s.Bytes()...)
	SetTxSignature(tx, string(temp))
}

func GetTxId(tx *Transactioner) {
	reflect.ValueOf(tx).Elem().FieldByName("Id")
}

func SetTxId(tx *Transactioner, id string) {
	reflect.ValueOf(&tx).Elem().FieldByName("Id").SetString(id)
}

func SetTxSignature(tx *Transactioner, sig string) {
	reflect.ValueOf(&tx).Elem().FieldByName("Signature").SetString(sig)
}
