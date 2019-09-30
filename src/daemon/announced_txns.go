package daemon

import (
	"sync"
	"time"

	"github.com/SkycoinProject/skycoin/src/cipher"
)

type announcedTxnsCache struct {
	sync.Mutex
	cache map[cipher.SHA256]int64
}

func newAnnouncedTxnsCache() *announcedTxnsCache {
	return &announcedTxnsCache{
		cache: make(map[cipher.SHA256]int64),
	}
}

func (c *announcedTxnsCache) add(txns []cipher.SHA256) {
	c.Lock()
	defer c.Unlock()

	t := time.Now().UTC().UnixNano()
	for _, txn := range txns {
		c.cache[txn] = t
	}
}

func (c *announcedTxnsCache) flush() map[cipher.SHA256]int64 {
	c.Lock()
	defer c.Unlock()

	if len(c.cache) == 0 {
		return nil
	}

	cache := c.cache

	c.cache = make(map[cipher.SHA256]int64)

	return cache
}
