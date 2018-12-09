package transactions

import (
	"math/rand"
	"time"
)

//GetNonce returns a new nonce
func GetNonce() [32]byte {
	temp := make([]byte, 32)
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	r1.Read(temp)
	var val [32]byte
	copy(val[:], temp[:])
	return val
}
