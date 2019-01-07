package peersterCrypto_test

import (
	"testing"

	pc "github.com/BjornGudmundsson/Peerster/data/peersterCrypto"
	"github.com/stretchr/testify/require"
)

func Test_EncryptDecryptAsymmetric(t *testing.T) {
	priv1 := pc.NewPrivateKey()
	priv2 := pc.NewPrivateKey()
	//pub1 := priv1.GetPublicKey()
	pub2 := priv2.GetPublicKey()
	//Get 128 bit key and IV
	IV := pc.GetIV(16)
	key := pc.GetKey(16)
	//Just to get a random byte string
	metafilehash := pc.GetIV(16)
	src := "Bjorn"
	fn := "Secret file"
	secret := pc.NewSecret(fn, src, "jon", pub2, metafilehash, IV, key)
	require.NotNil(t, secret)
	encryptedSecret, e := priv1.EncryptSecret(secret, pub2)
	require.Nil(t, e, "Should be able to encrypt")
	decryptedSecret, e := priv2.DecryptSecret(encryptedSecret)
	require.Nil(t, e, "Decrypting the secret should work")
	require.NotNil(t, decryptedSecret, "Should have gotten back a value")
	require.True(t, pc.CompareSecrets(secret, decryptedSecret), "Secrets should be equal")
	//Decrypting the secret using a wrong public key
	priv3 := pc.NewPrivateKey()
	_, e = priv3.DecryptSecret(encryptedSecret)
	require.NotNil(t, e, "A message should not be decryptable by any public key")
}
