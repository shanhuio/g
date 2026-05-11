package rand

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	mrand "math/rand"
	"sync"
	"time"
)

var randMutex sync.Mutex
var fallbackRand = mrand.New(
	mrand.NewSource(time.Now().UnixNano()),
)

// Bytes returns a byte slice of random bytes.
func Bytes(n int) []byte {
	ret := make([]byte, n)
	if _, err := rand.Read(ret); err == nil {
		return ret
	}

	randMutex.Lock()
	defer randMutex.Unlock()
	if _, err := fallbackRand.Read(ret); err != nil {
		panic(err)
	}

	return ret
}

// HexBytes returns the hex encoding of a random hex bytes
func HexBytes(n int) string {
	return hex.EncodeToString(Bytes(n))
}

// LowerLetters returns a random ID of n random letters, lower-case only.
func LowerLetters(n int) string {
	r := New()
	var ret bytes.Buffer

	for i := 0; i < n; i++ {
		x := r.Int31n(26)
		ret.WriteRune('a' + x)
	}
	return ret.String()
}

// Letters returns a random ID of n random case-sensitive letters.
// They might have lower case or upper case.
func Letters(n int) string {
	r := New()
	var ret bytes.Buffer

	for i := 0; i < n; i++ {
		x := r.Int31n(52)
		if x < 26 {
			ret.WriteRune('a' + x)
		} else {
			ret.WriteRune('A' + x - 26)
		}
	}
	return ret.String()
}

// Digits returns a string of n random digits.
func Digits(n int) string {
	r := New()
	var ret bytes.Buffer
	for i := 0; i < 10; i++ {
		x := r.Int31n(10)
		ret.WriteRune('0' + x)
	}
	return ret.String()
}

// New returns a new math/rand.Rand that is seeded with crypto rand.
func New() *mrand.Rand {
	seed := int64(binary.LittleEndian.Uint64(Bytes(8)))
	src := mrand.NewSource(seed)
	return mrand.New(src)
}
