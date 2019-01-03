package peersterCrypto

import "crypto/sha256"

//PublicPair is the public key and
//the corresponding origin
type PublicPair struct {
	PublicKey PublicKey
	Origin    string
}

//Compare compares if two PublicPairs are equivalent
func (pp *PublicPair) Compare(p *PublicPair) bool {
	h1 := pp.Hash()
	h2 := p.Hash()
	return compareBytes(h1[:], h2[:])
}

//Hash hashes a public pair and returns the
//32 byte representation of it.
func (pp *PublicPair) Hash() [32]byte {
	marshalledKey := pp.PublicKey.Marshall()
	originMarshalled := []byte(pp.Origin)
	h := sha256.New()
	h.Write(marshalledKey)
	h.Write(originMarshalled)
	var buf [32]byte
	copy(buf[:], h.Sum(nil))
	return buf
}
