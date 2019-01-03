package peersterCrypto

import "crypto/sha256"

//EncryptedSecret is the secret
//to be shared with a peer when
//trying to share a file. All
//of the variables are encrypted
//except for Origin and the
//PublicKey
type EncryptedSecret struct {
	FileName     []byte
	Origin       string
	Publickey    PublicKey
	MetaFileHash []byte
	IV           []byte
	Key          []byte
	Signature    []byte
}

//Hash hashes an Encrypted secret.
func (es *EncryptedSecret) Hash() [32]byte {
	h := sha256.New()
	h.Write(es.FileName)
	h.Write(es.MetaFileHash)
	h.Write(es.IV)
	h.Write(es.Key)
	h.Write(es.Signature)
	marshalled := es.Publickey.Marshall()
	h.Write(marshalled)
	h.Write([]byte(es.Origin))
	var buf [32]byte
	copy(buf[:], h.Sum(nil))
	return buf
}

//Compare bytes returns if two encrypted secrets are the same
func (es *EncryptedSecret) Compare(e *EncryptedSecret) bool {
	h1 := es.Hash()
	h2 := e.Hash()
	return compareBytes(h1[:], h2[:])
}

//NewEncryptedSecret returns a new encrypted secret
func NewEncryptedSecret(IV, fn, key, sign, mfh []byte, pk PublicKey, src string) *EncryptedSecret {
	return &EncryptedSecret{
		MetaFileHash: mfh,
		FileName:     fn,
		IV:           IV,
		Key:          key,
		Signature:    sign,
		Publickey:    pk,
		Origin:       src,
	}
}

//Secret is the decrypted
//secret.
type Secret struct {
	FileName     string
	Origin       string
	PublicKey    PublicKey
	MetaFileHash []byte
	IV           []byte
	Key          []byte
}

func compareBytes(b1, b2 []byte) bool {
	n := len(b1)
	if n != len(b2) {
		return false
	}
	for i := 0; i < n; i++ {
		if b1[i] != b2[i] {
			return false
		}
	}
	return true
}

//CompareSecrets compares two secrets
func CompareSecrets(s1, s2 *Secret) bool {
	if s1.FileName != s2.FileName {
		return false
	}
	if !compareBytes(s1.IV, s2.IV) {
		return false
	}
	if !compareBytes(s1.Key, s2.Key) {
		return false
	}
	if !compareBytes(s1.MetaFileHash, s2.MetaFileHash) {
		return false
	}
	return true
}

//NewSecret returns a pointer to a new secret
func NewSecret(fn, src string, pk PublicKey, mfh, iv, key []byte) *Secret {
	return &Secret{
		FileName:     fn,
		Origin:       src,
		PublicKey:    pk,
		MetaFileHash: mfh,
		IV:           iv,
		Key:          key,
	}
}
