package data

import (
	"crypto/rsa"
	"crypto/sha256"
	"encoding/binary"

	"github.com/BjornGudmundsson/Peerster/data/peersterCrypto"
)

type KeyBlock struct {
	PrevHash     [32]byte
	Nonce        [32]byte
	Transactions []KeyTransaction
}

type BlockRequest struct {
	Origin      string
	Destination string
	HopLimit    uint32
	HashValue   [32]byte
}

type BlockReply struct {
	Destination     string
	KeyBlockPublish *KeyBlockPublish
}

type KeyTransaction struct {
	KeyPublish *peersterCrypto.PublicPair
	Secret     *peersterCrypto.EncryptedSecret
}

type KeyPublish struct {
	Transaction *KeyTransaction
	HopLimit    uint32
}

type KeyBlockPublish struct {
	Origin   string
	Block    *KeyBlock
	HopLimit uint32
}

//Hash hashes a key blocks
func (b *KeyBlock) Hash() (out [32]byte) {
	h := sha256.New()
	h.Write(b.PrevHash[:])
	h.Write(b.Nonce[:])
	binary.Write(h, binary.LittleEndian, uint32(len(b.Transactions)))
	for _, t := range b.Transactions {
		th := t.Hash()
		h.Write(th[:])
	}
	copy(out[:], h.Sum(nil))
	return
}

//NewEncryptionKeyTransaction returns a new transactions that has a public key
//publication.
func NewEncryptionKeyTransaction(key rsa.PublicKey, name string) *KeyTransaction {
	publish := peersterCrypto.NewPublicPair(key, name)
	return &KeyTransaction{
		KeyPublish: publish,
	}
}

//NewSecretKeyTransaction creates a new key transactions that has
//a secret sharing transactions
func NewSecretKeyTransaction(secret *peersterCrypto.EncryptedSecret) *KeyTransaction {
	return &KeyTransaction{
		Secret: secret,
	}
}

//Hash returns the hash of a key transaction.
//Acts accordingly to which transactions is present.
func (t *KeyTransaction) Hash() (out [32]byte) {
	if t.KeyPublish != nil {
		out = t.KeyPublish.Hash()
	}
	if t.Secret != nil {
		out = t.Secret.Hash()
	}
	return
}

//IsKeyPublish checks if a transaction is for publishing a public key name pair
//or a secret sharing mechanism.
func (t *KeyTransaction) IsKeyPublish() bool {
	return !(t.KeyPublish == nil)
}

//GetName returns the name of the node that published a public key
//returns an empty string if there is no key publish in the transaction.
func (t *KeyTransaction) GetName() string {
	if !t.IsKeyPublish() {
		return ""
	}
	return t.KeyPublish.Origin
}

//Compare compares two key transactions. It checks which kind of transaction
//it is and uses the corresponding compare method.
func (t *KeyTransaction) Compare(tx *KeyTransaction) bool {
	if t.KeyPublish != nil && tx.KeyPublish != nil {
		return t.KeyPublish.Compare(tx.KeyPublish)
	}
	if t.Secret != nil && tx.Secret != nil {
		return t.Secret.Compare(tx.Secret)
	}
	return false
}

//GetPublicKey returns a pointer to a public pair. If the transaction
//has a secret that has been shared it will return nil.
func (t *KeyTransaction) GetPublicKey() *peersterCrypto.PublicPair {
	return t.KeyPublish
}
