package secp256k1

import (
	crand "crypto/rand" // secure, system random number generator
	"crypto/sha256"
	"hash"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	// Memory pool for SHA256 hashes
	sha256HashPoolSize = 30
	sha256HashPool     chan hash.Hash
)

// SumSHA256 sum sha256
func SumSHA256(b []byte) []byte {
	sha256Hash := <-sha256HashPool
	sha256Hash.Reset()
	// sha256.Write never returns an error
	sha256Hash.Write(b) //nolint:errcheck
	sum := sha256Hash.Sum(nil)
	sha256HashPool <- sha256Hash
	return sum[:]
}

/*
Entropy pool needs
- state (an array of bytes)
- a compression function (two 256 bit blocks to single block)
- a mixing function across the pool

- Xor is safe, as it cannot make value less random
-- apply compression function, then xor with current value
--

*/

// EntropyPool entropy pool
type EntropyPool struct {
	Ent  [32]byte // 256 bit accumulator
	lock sync.Mutex
}

// Mix256 mixes in 256 bits, outputs 256 bits
func (ep *EntropyPool) Mix256(in []byte) (out []byte) {
	// hash input
	val1 := SumSHA256(in)
	// return value
	ep.lock.Lock()
	val2 := SumSHA256(append(val1, ep.Ent[:]...))
	// next ent value
	val3 := SumSHA256(append(val1, val2...))

	for i := 0; i < 32; i++ {
		ep.Ent[i] = val3[i]
		val3[i] = 0x00
	}
	ep.lock.Unlock()

	return val2
}

// Mix take in N bytes, salts, return N
func (ep *EntropyPool) Mix(in []byte) []byte {
	length := len(in) - len(in)%32 + 32
	buff := make([]byte, length)
	for i := 0; i < len(in); i++ {
		buff[i] = in[i]
	}
	iterations := (len(in) / 32) + 1
	for i := 0; i < iterations; i++ {
		tmp := ep.Mix256(buff[32*i : 32+32*i]) //32 byte slice
		for j := 0; j < 32; j++ {
			buff[i*32+j] = tmp[j]
		}
	}
	return buff[:len(in)]
}

/*
Note:

- On windows cryto/rand uses CrytoGenRandom which uses RC4 which is insecure
- Android random number generator is known to be insecure.
- Linux uses /dev/urandom , which is thought to be secure and uses entropy pool

Therefore the output is salted.
*/

/*
Note:

Should allow pseudo-random mode for repeatability for certain types of tests
*/

var _ent EntropyPool

// seed pseudo random number generator with
// - hash of system time in nano seconds
// - hash of system environmental variables
// - hash of process id
func init() {
	// init the hash reuse pool
	sha256HashPool = make(chan hash.Hash, sha256HashPoolSize)
	for i := 0; i < sha256HashPoolSize; i++ {
		sha256HashPool <- sha256.New()
	}

	seed1 := []byte(strconv.FormatUint(uint64(time.Now().UnixNano()), 16))
	seed2 := []byte(strings.Join(os.Environ(), ""))
	seed3 := []byte(strconv.FormatUint(uint64(os.Getpid()), 16))

	seed4 := make([]byte, 256)
	_, err := io.ReadFull(crand.Reader, seed4) // system secure random number generator
	if err != nil {
		log.Panic(err)
	}

	// seed entropy pool
	_ent.Mix256(seed1)
	_ent.Mix256(seed2)
	_ent.Mix256(seed3)
	_ent.Mix256(seed4)
}

// RandByte Secure Random number generator for forwards security
// On Unix-like systems, Reader reads from /dev/urandom.
// On Windows systems, Reader uses the CryptGenRandom API.
// Pseudo-random sequence, seeded from program start time, environmental variables,
// and process id is mixed in for forward security. Future version should use entropy pool
func RandByte(n int) []byte {
	buff := make([]byte, n)
	_, err := io.ReadFull(crand.Reader, buff) // system secure random number generator
	if err != nil {
		log.Panic(err)
	}

	// XORing in sequence, cannot reduce security (even if sequence is bad/known/non-random)
	buff2 := _ent.Mix(buff)
	for i := 0; i < n; i++ {
		buff[i] ^= buff2[i]
	}
	return buff
}
