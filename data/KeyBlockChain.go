package data

import (
	"crypto/rsa"
	"crypto/sha256"
	"encoding/binary"
)

type KeyBlock struct {
	PrevHash [32]byte
	Nonce [32]byte
	Transactions []KeyTransaction
}

type BlockRequest struct {
	Origin string
	Destination string
	HopLimit uint32
	HashValue [32]byte
}

type BlockReply struct {
	Origin string
	Destination string
	HopLimit uint32
	Block KeyBlock
}

type KeyTransaction struct {
	Name string
	Key rsa.PublicKey
}

type KeyPublish struct {
	Transaction *KeyTransaction
	HopLimit uint32
}

type KeyBlockPublish struct {
	Origin string
	Block *KeyBlock
	HopLimit uint32
}



func (b *KeyBlock) Hash() (out [32]byte) {
	h := sha256.New()
	h.Write(b.PrevHash[:])
	h.Write(b.Nonce[:])
	binary.Write(h,binary.LittleEndian, uint32(len(b.Transactions)))
	for _, t := range b.Transactions {
		th := t.Hash()
		h.Write(th[:])
	}
	copy(out[:], h.Sum(nil))
	return
}

func (t *KeyTransaction) Hash() (out [32]byte) {
	h := sha256.New()
	h.Write([]byte(t.Name))
	binary.Write(h,binary.LittleEndian, uint32(t.Key.E))
	h.Write(t.Key.N.Bytes())
	copy(out[:], h.Sum(nil))
	return
}
