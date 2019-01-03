package transactions

import pc "github.com/BjornGudmundsson/Peerster/data/peersterCrypto"

//Transaction is the transction
//of either a public key being
//logged or a secret being shared
//Either the EncryptedSecret is
//not nil or the Public pair is
//not nil, not both.
type Transaction struct {
	EncryptedSecret *pc.EncryptedSecret
	PublicPair      *pc.PublicPair
}

//Hash returns the hash of a transaction
func (t *Transaction) Hash() [32]byte {
	if t.EncryptedSecret != nil {
		return t.EncryptedSecret.Hash()
	}
	if t.PublicPair != nil {
		return t.PublicPair.Hash()
	}
	var buf [32]byte
	return buf
}
