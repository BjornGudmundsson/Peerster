package peersterCrypto_test

import (
	"testing"

	pc "github.com/BjornGudmundsson/Peerster/data/peersterCrypto"
	"github.com/stretchr/testify/require"
)

func Test_EncryptDecryptSymmetric(t *testing.T) {
	IV := pc.GetIV(16)
	Key := pc.GetKey(16)
	//What is to be encrypted
	str := "HHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHH"
	data := []byte(str)
	ct, e := pc.EncryptBlocks(data, Key, IV)
	require.Nil(t, e, "Something went wrong in the encryption process")
	pt, e := pc.DecryptCiphertext(ct, Key, IV)
	require.Nil(t, e, "Something went wrong in the decryption")
	require.True(t, string(pt) == str, "The decrypted ciphertext did not equal the original plaintext")
}
