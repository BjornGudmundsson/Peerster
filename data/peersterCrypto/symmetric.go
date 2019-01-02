package peersterCrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

//GetIV returns a new randomly
//chosen IV of a given size.
func GetIV(size int) []byte {
	buf := make([]byte, size)
	rand.Read(buf)
	return buf
}

//GetKey returns a key of a given lenght
func GetKey(size int) []byte {
	return GetIV(size)
}

//EncryptBlocks takes in a byte array that must be a
//multiple of 16, 24 or 32 bytes. If the message to be encrypted is not a multiple of
//16, 24, 32 pad it with zero bytes.
func EncryptBlocks(blocks, key, IV []byte) ([]byte, error) {
	ciph, e := aes.NewCipher(key)
	if e != nil {
		return nil, e
	}
	cbcenc := cipher.NewCBCEncrypter(ciph, IV)
	n := len(blocks)
	dst := make([]byte, n)
	copy(dst, blocks)
	cbcenc.CryptBlocks(dst, blocks)
	return dst, nil
}

//DecryptCiphtertext takes in a ciphertext, key and the IV
//and returns the corresponding plaintext
func DecryptCiphertext(ct, key, IV []byte) ([]byte, error) {
	ciph, e := aes.NewCipher(key)
	if e != nil {
		return nil, e
	}
	n := len(ct)
	cbcdec := cipher.NewCBCDecrypter(ciph, IV)
	dst := make([]byte, n)
	copy(dst, ct)
	cbcdec.CryptBlocks(dst, ct)
	return dst, nil
}
