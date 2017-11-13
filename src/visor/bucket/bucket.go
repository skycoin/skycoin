package bucket

import (
	"encoding/binary"

	"github.com/boltdb/bolt"
)

// Itob converts uint64 to bytes
func Itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// Btoi converts bytes to uint64
func Btoi(v []byte) uint64 {
	return binary.BigEndian.Uint64(v)
}

// Rollback callback function type
type Rollback func()

// TxHandler function type for processing bolt transaction
type TxHandler func(tx *bolt.Tx) (Rollback, error)
