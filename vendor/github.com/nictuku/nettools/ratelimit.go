package nettools

import (
	"time"

	"github.com/youtube/vitess/go/cache"
)

// NewThrottle creates a new client throttler that blocks spammy clients.
// UPDATED in 2015-01-17: clients now have to specify the limits. Use 10 and
// 1000 if you want to use the old default values.
func NewThrottler(maxPerMinute int, maxHosts int64) *ClientThrottle {
	r := ClientThrottle{
		maxPerMinute: maxPerMinute,
		c:            cache.NewLRUCache(maxHosts),
		blocked:      cache.NewLRUCache(maxHosts),
		stop:         make(chan bool),
	}
	go r.cleanup()
	return &r
}

// ClientThrottle identifies and blocks hosts that are too spammy. It only
// cares about the number of operations per minute.
type ClientThrottle struct {
	maxPerMinute int

	// Rate limiter.
	c *cache.LRUCache

	// Hosts that get blocked once go to a separate cache, and stay forever
	// until they stop hitting us enough to fall off the blocked cache.
	blocked *cache.LRUCache

	// This channel will be closed when the ClientThrottle should be stopped.
	stop chan bool
}

// Stops the ClientThrottle and all internal goroutines.
func (r *ClientThrottle) Stop() {
	close(r.stop)

}

func (r *ClientThrottle) CheckBlock(host string) bool {
	_, blocked := r.blocked.Get(host)
	if blocked {
		// Bad guy stays there.
		return false
	}

	v, ok := r.c.Get(host)
	var h hits
	if !ok {
		h = hits(59)
	} else {
		h = v.(hits) - 1
	}
	if int(h) < 60-r.maxPerMinute {
		// fmt.Printf("blocking because int(h)=%v < 60-r.maxPerMinute = (60-%v) => %v\n", int(h), r.maxPerMinute, 60-r.maxPerMinute)
		r.c.Set(host, h-300)
		// New bad guy.
		r.blocked.Set(host, h) // The value here is not relevant.
		return false
	}
	r.c.Set(host, h)
	return true
}

// refill the buckets.
// this is the first way I thought of how to implement client rate limiting.
// Need to think and research more.
func (r *ClientThrottle) cleanup() {
	// Check the bucket faster than the rate period, to reduce the pressure in the cache.
	t := time.Tick(5 * time.Second)

	for {
		select {
		case <-t:
			var h hits
			// This is ridiculously inefficient but it'll have to do for now.
			for _, item := range r.c.Items() {
				h = item.Value.(hits) + 5
				if h > 60 {
					// Reduce pressure in the LRU.
					r.c.Delete(item.Key)
				} else {
					r.c.Set(item.Key, h)
				}
			}
		case <-r.stop:
			return
		}
	}
}

type hits int

func (h hits) Size() int {
	return 1
}
