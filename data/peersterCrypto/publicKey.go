package peersterCrypto

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/binary"
	"math/big"
)

//PrivateKey is a holder for the
//private key that the gossiper is using.
//The underlying crypto system should not matter
//as long as you are using the private key struct
//and its following methods.
type PrivateKey struct {
	privateKey *rsa.PrivateKey
}

//PublicKey is a struct that
//contains the public key
//of a private key
type PublicKey struct {
	E int
	N []byte
}

//GetKey returns the public key of the node
func (pk PublicKey) GetKey() rsa.PublicKey {
	N := big.NewInt(0).SetBytes(pk.N)
	pub := rsa.PublicKey{
		E: pk.E,
		N: N,
	}
	return pub
}

//NewPublicPair returns a new public pair.
func NewPublicPair(key rsa.PublicKey, name string) *PublicPair {
	bs := key.N.Bytes()
	pub := PublicKey{
		E: key.E,
		N: bs,
	}
	pair := &PublicPair{
		PublicKey: pub,
		Origin:    name,
	}
	return pair
}

//NewPrivateKey returns a new instance of
//a randomly generated privatekey.
func NewPrivateKey() *PrivateKey {
	priv := getPrivateKey()
	if priv == nil {
		return nil
	}
	return &PrivateKey{
		privateKey: priv,
	}
}

//GetPrivateKey returns a privatekey
func getPrivateKey() *rsa.PrivateKey {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil
	}
	return priv
}

//GetPublicKey returns the corresponding public key
//for the private key
func (priv *PrivateKey) GetPublicKey() PublicKey {
	pub := priv.privateKey.PublicKey
	return PublicKey{
		E: pub.E,
		N: pub.N.Bytes(),
	}
}

//Encrypt encrypts a message with the public key
func (pk *PublicKey) Encrypt(msg []byte) ([]byte, error) {
	pub := pk.GetKey()
	hash := sha256.New()
	label := []byte("")
	ciphertext, err := rsa.EncryptOAEP(
		hash,
		rand.Reader,
		&pub,
		msg,
		label,
	)
	if err != nil {
		return nil, err
	}
	return ciphertext, err
}

//Marshall marshalls the public key to a byte array
func (pk PublicKey) Marshall() []byte {
	e := pk.E
	E := big.NewInt(int64(e))
	N := pk.N
	return append(buf.Bytes(), N...)
}

//Signature takes in a message and creates a signature. Use the
//concatenated values of the IV and the key to create a signature
//since they are the most unique parts of the message.
func (priv *PrivateKey) Signature(msg []byte) ([]byte, error) {
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto // for simple example
	PSSmessage := msg
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPSS(
		rand.Reader,
		priv.privateKey,
		newhash,
		hashed, &opts)
	return signature, err
}

//VerifySignature takes in a signature and verifies that this message
//was signed by the secret key corresponding to this public key.
func (pk PublicKey) VerifySignature(sign []byte, msg []byte) error {
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(msg)
	hashed := pssh.Sum(nil)
	var opts rsa.PSSOptions
	pub := pk.GetKey()
	err := rsa.VerifyPSS(
		&pub,
		newhash,
		hashed,
		sign, &opts)
	return err
}

//EncryptSecret returns a new encryped secret with a signature.
func (priv *PrivateKey) EncryptSecret(s *Secret, pk PublicKey) (*EncryptedSecret, error) {
	myPk := priv.GetPublicKey()
	encIV, err := pk.Encrypt(s.IV)
	encFN, err := pk.Encrypt([]byte(s.FileName))
	encKey, err := pk.Encrypt(s.Key)
	encMFH, err := pk.Encrypt(s.MetaFileHash)
	if err != nil {
		return nil, err
	}
	msg := append(s.IV, s.Key...)
	sign, err := priv.Signature(msg)
	if err != nil {
		return nil, err
	}
	ecns := NewEncryptedSecret(encIV, encFN, encKey, sign, encMFH, myPk, s.Origin, s.Destination)
	return ecns, nil
}

//DecryptSecret takes an encrypted secret and decrypts it.
func (priv *PrivateKey) DecryptSecret(es *EncryptedSecret) (*Secret, error) {
	pk := es.Publickey
	fn, err := priv.Decrypt(es.FileName)
	IV, err := priv.Decrypt(es.IV)
	MFH, err := priv.Decrypt(es.MetaFileHash)
	key, err := priv.Decrypt(es.Key)
	if err != nil {
		return nil, err
	}
	src := es.Origin
	msg := append(IV, key...)
	verify := pk.VerifySignature(es.Signature, msg)
	if verify != nil {
		return nil, verify
	}
	secret := NewSecret(string(fn), src, es.Destination, pk, MFH, IV, key)
	return secret, nil
}

//Decrypt decrypt a message using the private key.
func (priv *PrivateKey) Decrypt(c []byte) ([]byte, error) {
	p := priv.privateKey
	hash := sha256.New()
	label := []byte("")
	pt, err := rsa.DecryptOAEP(
		hash,
		rand.Reader,
		p,
		c,
		label,
	)
	return pt, err
}
