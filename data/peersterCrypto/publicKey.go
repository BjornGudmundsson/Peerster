package peersterCrypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
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
	publicKey rsa.PublicKey
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
		publicKey: pub,
	}
}

//Encrypt encrypts a message with the public key
func (pk *PublicKey) Encrypt(msg []byte) ([]byte, error) {
	pub := &pk.publicKey
	hash := sha256.New()
	label := []byte("")
	ciphertext, err := rsa.EncryptOAEP(
		hash,
		rand.Reader,
		pub,
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
	e := pk.publicKey.E
	E := big.NewInt(int64(e))
	N := pk.publicKey.N
	return append(E.Bytes(), N.Bytes()...)
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
	err := rsa.VerifyPSS(
		&pk.publicKey,
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
	ecns := NewEncryptedSecret(encIV, encFN, encKey, sign, encMFH, myPk, s.Origin)
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
	secret := NewSecret(string(fn), src, pk, MFH, IV, key)
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