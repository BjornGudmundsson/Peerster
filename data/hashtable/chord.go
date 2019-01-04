package hashtable

import (
	"crypto/sha1"
	"math/big"
)

//HashStringInt takes a string and
//returns the corresponding big int
//representation
func HashStringInt(s string) *big.Int {
	h := sha1.New()
	h.Write([]byte(s))
	b := h.Sum(nil)
	bigint := big.NewInt(0)
	return bigint.SetBytes(b)
}

//FindPlaceInChord find the place where a big int belongs in a chord
func FindPlaceInChord(a []*big.Int, v *big.Int) (*big.Int, *big.Int) {
	n := len(a)
	if n == 0 {
		return nil, nil
	}
	if n < 2 {
		return a[0], nil
	}
	for i := 0; i < n-1; i++ {
		val := a[i]
		val2 := a[i+1]
		if val.Cmp(v) <= 0 && val2.Cmp(v) > 0 {
			return val, val2
		}
	}
	return a[n-1], a[0]
}
