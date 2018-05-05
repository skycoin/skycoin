package dbutil

import "github.com/boltdb/bolt"

// Everything in this file is deprecated and will be removed,
// it's here while transitioning away from visor/bucket to visor/dbutil

// Rollback callback function type
type Rollback func()

// TxHandler function type for processing bolt transaction
type TxHandler func(tx *bolt.Tx) (Rollback, error)
